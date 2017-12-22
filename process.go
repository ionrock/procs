package procs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

func lineReader(wg *sync.WaitGroup, r io.Reader, out chan string) {
	defer wg.Done()

	reader := bufio.NewReader(r)
	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)

		if n, err := reader.Read(buf); err != nil {
			return
		} else {
			buf = buf[:n]
		}

		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				break
			}
			buffer.Write(buf[0:i])
			out <- buffer.String()
			buffer.Reset()
			buf = buf[i+1:]
		}
		buffer.Write(buf)
	}

	out <- buffer.String()
}

type Process struct {
	// CmdString takes a string and parses it into the relevant cmds
	CmdString string

	// Cmds is the list of command delmited by pipes.
	Cmds []*exec.Cmd

	// Env provides a map[string]string that can mutated before
	// running a command.
	Env map[string]string

	// Dir defines the directory the command should run in. The
	// Default is the current dir.
	Dir string

	// OutputHandler can be defined to perform any sort of processing
	// on the output. The simple interface is to accept a string (a
	// line of output) and return a string that will be included in the
	// buffered output and/or output written to stdout.
	OutputHandler func(string) string

	// TODO: This can really be private
	Pipes []*io.PipeWriter

	// When no output is given, we'll buffer output in these vars.
	errBuffer bytes.Buffer
	outBuffer bytes.Buffer

	// When a output handler is provided, we ensure we're handling a
	// single line at at time.
	outputWait *sync.WaitGroup
	stdoutChan chan string
	stderrChan chan string
}

// NewProcess creates a new *Process from a command string.
func NewProcess(command string) *Process {
	return &Process{CmdString: command}
}

// internal expand method to use the os env or proc env.
func (p *Process) expand(s string) string {
	if p.Env != nil {
		return os.ExpandEnv(s)
	}

	return os.Expand(s, func(key string) string {
		v, _ := p.Env[key]
		return v
	})
}

// addCmd adds a new command to the list of commands, ensuring the Dir
// and Env have been added to the underlying *exec.Cmd instances.
func (p *Process) addCmd(cmdparts []string) {
	var cmd *exec.Cmd
	if len(cmdparts) == 1 {
		cmd = exec.Command(cmdparts[0])
	} else {
		cmd = exec.Command(cmdparts[0], cmdparts[1:]...)
	}

	if p.Dir != "" {
		cmd.Dir = p.Dir
	}

	if p.Env != nil {
		env := []string{}
		for k, v := range p.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, p.expand(v)))
		}

		cmd.Env = env
	}

	p.Cmds = append(p.Cmds, cmd)
}

// findCmds parses the CmdString to find the commands that should be
// run by spliting the lexically parsed command by pipes ("|").
func (p *Process) findCmds() {
	// Skip if the cmd set is already set. This allows manual creation
	// of piped commands.
	if len(p.Cmds) > 0 {
		return
	}

	if p.CmdString == "" {
		return
	}

	parts := SplitCommand(p.CmdString)
	for i := range parts {
		parts[i] = p.expand(parts[i])
	}

	cmd := []string{}
	for _, part := range parts {
		if part == "|" {
			p.addCmd(cmd)
			cmd = []string{}
		} else {
			cmd = append(cmd, part)
		}
	}

	p.addCmd(cmd)
}

func (p *Process) setupOutputHandler() error {
	last := len(p.Cmds) - 1
	stdout, err := p.Cmds[last].StdoutPipe()
	if err != nil {
		fmt.Println("error creating stdout pipe")
		return err
	}

	stderr, err := p.Cmds[last].StderrPipe()
	if err != nil {
		fmt.Println("error creating stderr pipe")
		return err
	}

	p.stdoutChan = make(chan string)
	p.stderrChan = make(chan string)

	p.outputWait = new(sync.WaitGroup)
	p.outputWait.Add(2)

	go lineReader(p.outputWait, stdout, p.stdoutChan)
	go lineReader(p.outputWait, stderr, p.stderrChan)

	go func() {
		for line := range p.stdoutChan {
			_, err := p.outBuffer.WriteString(p.OutputHandler(line))
			if err != nil {
				fmt.Printf("failed to write to buffer: %s", err)
			}
		}
	}()

	go func() {
		for line := range p.stderrChan {
			_, err := p.outBuffer.WriteString(p.OutputHandler(line))
			if err != nil {
				fmt.Printf("failed to write to buffer: %s", err)
			}
		}
	}()

	return nil
}

func (p *Process) setupPipes() error {
	if p.Pipes == nil {
		p.Pipes = make([]*io.PipeWriter, len(p.Cmds)-1)
	}

	i := 0
	for ; i < len(p.Cmds)-1; i++ {
		stdinPipe, stdoutPipe := io.Pipe()
		p.Cmds[i].Stdout = stdoutPipe
		p.Cmds[i].Stderr = &p.errBuffer

		// set the input to the outoput
		p.Cmds[i+1].Stdin = stdinPipe
		p.Pipes[i] = stdoutPipe
	}

	if p.OutputHandler != nil {
		err := p.setupOutputHandler()
		return err
	}

	p.Cmds[i].Stdout = &p.outBuffer
	p.Cmds[i].Stderr = &p.errBuffer

	return nil
}

func (p *Process) call(index int, wait bool) error {
	// This hasn't already been started so start it
	if p.Cmds[index].Process == nil {
		if err := p.Cmds[index].Start(); err != nil {
			return err
		}
	}

	// See if we have more cmds to run and start them by recursively calling Call
	if len(p.Cmds[index:]) > 1 {
		err := p.Cmds[index+1].Start()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				p.Pipes[index].Close()
				err = p.call(index+1, wait)
			}
		}()
	}

	last := len(p.Cmds) - 1
	if index == last && wait == true {
		return nil
	}

	return p.Cmds[index].Wait()
}

func (p *Process) Run() (string, error) {
	p.findCmds()
	p.setupPipes()

	if err := p.call(0, false); err != nil {
		fmt.Printf("error calling command: %q\n", err)
		fmt.Println(string(p.errBuffer.Bytes()))
		return "", err
	}

	return p.outBuffer.String(), nil
}

func (p *Process) Start() error {
	p.findCmds()
	p.setupPipes()

	if err := p.call(0, true); err != nil {
		fmt.Printf("error calling command: %q\n", err)
		fmt.Println(string(p.errBuffer.Bytes()))
		return err
	}

	return nil
}

func (p *Process) Wait() error {
	if p.outputWait != nil {
		p.outputWait.Wait()
	}

	last := len(p.Cmds) - 1
	return p.Cmds[last].Wait()
}

func (p *Process) Output() string {
	return p.outBuffer.String()
}

package procs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

type PipedProc struct {
	Cmds          []*exec.Cmd
	Pipes         []*io.PipeWriter
	OutputHandler func(string) string

	errBuffer bytes.Buffer
	outBuffer bytes.Buffer

	outputWait *sync.WaitGroup
	stdoutChan chan string
	stderrChan chan string
}

func (pp *PipedProc) setupOutputHandler() error {
	last := len(pp.Cmds) - 1
	stdout, err := pp.Cmds[last].StdoutPipe()
	if err != nil {
		fmt.Println("error creating stdout pipe")
		return err
	}

	stderr, err := pp.Cmds[last].StderrPipe()
	if err != nil {
		fmt.Println("error creating stderr pipe")
		return err
	}

	pp.stdoutChan = make(chan string)
	pp.stderrChan = make(chan string)

	pp.outputWait = new(sync.WaitGroup)
	pp.outputWait.Add(2)

	go lineReader(pp.outputWait, stdout, pp.stdoutChan)
	go lineReader(pp.outputWait, stderr, pp.stderrChan)

	go func() {
		for line := range pp.stdoutChan {
			// TODO: Make a decision to buffer or not. For now buffer
			_, err := pp.outBuffer.WriteString(pp.OutputHandler(line))
			if err != nil {
				fmt.Printf("failed to write to buffer: %s", err)
			}
		}
	}()

	return nil
}

func (pp *PipedProc) setupPipes() error {
	if pp.Pipes == nil {
		pp.Pipes = make([]*io.PipeWriter, len(pp.Cmds)-1)
	}

	i := 0
	for ; i < len(pp.Cmds)-1; i++ {
		stdinPipe, stdoutPipe := io.Pipe()
		pp.Cmds[i].Stdout = stdoutPipe
		pp.Cmds[i].Stderr = &pp.errBuffer

		// set the input to the outoput
		pp.Cmds[i+1].Stdin = stdinPipe
		pp.Pipes[i] = stdoutPipe
	}

	if pp.OutputHandler != nil {
		err := pp.setupOutputHandler()
		return err
	}

	pp.Cmds[i].Stdout = &pp.outBuffer
	pp.Cmds[i].Stderr = &pp.errBuffer

	return nil
}

func (pp *PipedProc) call(index int) error {
	// This hasn't already been started so start it
	if pp.Cmds[index].Process == nil {
		if err := pp.Cmds[index].Start(); err != nil {
			return err
		}
	}

	// See if we have more cmds to run and start them by recursively calling Call
	if len(pp.Cmds[index:]) > 1 {
		err := pp.Cmds[index+1].Start()
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				pp.Pipes[index].Close()
				err = pp.call(index + 1)
			}
		}()
	}

	return pp.Cmds[index].Wait()
}

func (pp *PipedProc) Run() (string, error) {
	pp.setupPipes()

	if pp.OutputHandler != nil {
	}

	if err := pp.call(0); err != nil {
		fmt.Printf("error calling command: %q\n", err)
		fmt.Println(string(pp.errBuffer.Bytes()))
		return "", err
	}

	return pp.outBuffer.String(), nil
}

package procs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Proc struct {
	Cmd       *exec.Cmd
	Output    *Output
	pipesWait *sync.WaitGroup
}

func New(command string) *Proc {
	var cmd *exec.Cmd
	parts := SplitCommand(command)
	if len(parts) == 0 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	return &Proc{
		Cmd:       cmd,
		pipesWait: new(sync.WaitGroup),
	}
}

// Run the exec.Cmd handling stdout / stderr according to the
// configured Output struct.
func (p *Proc) Run() error {
	if p.Output == nil {
		return p.Cmd.Run()
	}

	err := p.Start()
	if err != nil {
		p.Output.SystemOutput(fmt.Sprint("Failed to start ", p.Cmd.Args, ": ", err))
		return err
	}

	return p.Wait()
}

// Start the exec.Cmd start.
func (p *Proc) Start() error {
	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error creating stdout pipe")
		return err
	}
	stderr, err := p.Cmd.StderrPipe()
	if err != nil {
		fmt.Println("error creating stderr pipe")
		return err
	}

	if p.pipesWait == nil {
		p.pipesWait = new(sync.WaitGroup)
	}
	p.pipesWait.Add(2)

	go p.Output.LineReader(p.pipesWait, p.Output.Name, stdout, false)
	go p.Output.LineReader(p.pipesWait, p.Output.Name, stderr, true)

	return p.Cmd.Start()
}

func (p *Proc) Wait() error {
	if p.pipesWait != nil {
		p.pipesWait.Wait()
	}
	return p.Cmd.Wait()
}

func ParseEnv(environ []string) map[string]string {
	env := make(map[string]string)
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}
	return env
}

func Env(env map[string]string, useEnv bool) []string {
	envlist := []string{}

	// update our env by loading our env and overriding any values in
	// the provided env.
	if useEnv {
		environ := ParseEnv(os.Environ())
		for k, v := range env {
			environ[k] = v
		}
		env = environ
	}

	for key, val := range env {
		if key == "" {
			continue
		}
		envlist = append(envlist, fmt.Sprintf("%s=%s", key, val))
	}

	return envlist
}

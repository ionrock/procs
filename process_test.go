package procs_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/ionrock/procs"
)

func newProcess() *procs.Process {
	return &procs.Process{
		Cmds: []*exec.Cmd{
			exec.Command("echo", "foo"),
			exec.Command("grep", "foo"),
		},
	}
}

func TestProcess(t *testing.T) {
	p := newProcess()

	out, err := p.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	if strings.TrimSpace(out) != "foo" {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestProcessWithOutput(t *testing.T) {
	p := newProcess()

	p.OutputHandler = func(line string) string {
		return fmt.Sprintf("x | %s", line)
	}

	out, err := p.Run()

	if err != nil {
		t.Fatalf("error running program: %s", err)
	}
	expected := "x | foo"
	if strings.TrimSpace(out) != expected {
		t.Errorf("wrong output: expected %q but got %q", expected, out)
	}
}

func TestProcessStartAndWait(t *testing.T) {
	p := newProcess()

	p.Start()
	p.Wait()

	out := p.Output()
	if strings.TrimSpace(out) != "foo" {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestProcessStartAndWaitWithOutput(t *testing.T) {
	p := newProcess()
	p.OutputHandler = func(line string) string {
		return fmt.Sprintf("x | %s", line)
	}

	p.Start()
	p.Wait()

	out := p.Output()
	expected := "x | foo"
	if strings.TrimSpace(out) != expected {
		t.Errorf("wrong output: expected %q but got %q", expected, out)
	}
}

func TestProcessFromString(t *testing.T) {
	p := procs.NewProcess("echo 'foo'")
	out, err := p.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	if strings.TrimSpace(out) != "foo" {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestProcessFromStringWithPipe(t *testing.T) {
	p := procs.NewProcess("echo 'foo' | grep foo")
	out, err := p.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	if strings.TrimSpace(out) != "foo" {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

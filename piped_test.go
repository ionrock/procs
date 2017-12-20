package procs_test

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/ionrock/procs"
)

func newPipedProc() *procs.PipedProc {
	return &procs.PipedProc{
		Cmds: []*exec.Cmd{
			exec.Command("echo", "foo"),
			exec.Command("grep", "foo"),
		},
		Pipes: make([]*io.PipeWriter, 2),
	}
}

func TestPipedProc(t *testing.T) {
	pp := newPipedProc()

	out, err := pp.Run()
	if err != nil {
		t.Fatalf("error running program: %s", err)
	}

	if strings.TrimSpace(out) != "foo" {
		t.Errorf("wrong output: expected foo but got %s", out)
	}
}

func TestPipedProcWithOutput(t *testing.T) {
	pp := newPipedProc()

	pp.OutputHandler = func(line string) string {
		return fmt.Sprintf("x | %s", line)
	}

	out, err := pp.Run()

	if err != nil {
		t.Fatalf("error running program: %s", err)
	}
	expected := "x | foo"
	if strings.TrimSpace(out) != expected {
		t.Errorf("wrong output: expected %q but got %q", expected, out)
	}
}

package procs

import (
	"log"
	"os"
	"strings"

	shlex "github.com/flynn/go-shlex"
)

// SplitCommand parses a command and splits it into lexical arguments
// like a shell, returning a []string that can be used as arguments to
// exec.Command.
func SplitCommand(cmd string) []string {
	return SplitCommandEnv(cmd, os.Getenv)
}

// SplitCommandEnv parses a command and splits it into lexical
// arguments like a shell, returning a []string that can be used as
// arguments to exec.Command. It also allows providing an expansion
// function that will be used when expanding values within the parsed
// arguments.
func SplitCommandEnv(cmd string, getenv func(key string) string) []string {
	parts, err := shlex.Split(strings.TrimSpace(strings.Trim(os.Expand(cmd, getenv), "`")))
	if err != nil {
		log.Fatal(err)
	}
	return parts
}

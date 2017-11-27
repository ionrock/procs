package procs

import (
	"log"
	"os"
	"strings"

	shlex "github.com/flynn/go-shlex"
)

func SplitCommand(cmd string) []string {
	return SplitCommandEnv(cmd, os.Getenv)
}

func SplitCommandEnv(cmd string, getenv func(key string) string) []string {
	parts, err := shlex.Split(strings.Trim(os.Expand(cmd, getenv), "`"))
	if err != nil {
		log.Fatal(err)
	}
	return parts
}

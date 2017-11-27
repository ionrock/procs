package procs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func New(command string) *exec.Cmd {
	var cmd *exec.Cmd
	parts := SplitCommand(command)
	if len(parts) == 0 {
		cmd = exec.Command(parts[0])
	} else {
		cmd = exec.Command(parts[0], parts[1:]...)
	}

	return cmd
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

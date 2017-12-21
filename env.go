package procs

import (
	"fmt"
	"os"
	"strings"
)

func ParseEnv(environ []string) map[string]string {
	env := make(map[string]string)
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)

		// There is a chance we can get an env with empty values
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
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

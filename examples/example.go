package foo

import (
	"os/exec"
	"strings"

	"github.com/ionrock/procs"
)

type Knife struct {
	Builder *procs.Builder
}

func (k *Knife) Command(model, action string, args []string) *exec.Cmd {
	command := k.Builder.Command(map[string]string{
		"model":  model,
		"action": action,
		"args":   strings.Join(args, " "),
	})

	// This will parse the command correctly into is distince parts.
	cmd := procs.New(command)

	return cmd
}

func main() {
	knife := &Knife{
		Builder: &procs.Builder{
			// Provide some defaults
			Context: map[string]string{
				"server": "https://chef.exmaple.org",
				"key":    "/etc/chef/knife.pem",
			},

			// Create a list of templates that use normal shell expansion
			Templates: []string{
				"knife", "${model}", "${action}", "${args}",
				"-s ${server}", "-k ${key}", "-Fj",
			},
		},
	}

	cmd := knife.Command("cookbook", "show", "foo")
	cmd.Run()
}

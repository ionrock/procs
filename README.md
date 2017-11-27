# Procs

Procs is a library to make working with command line applications more
like working with an API.

The primary use case is when you have to use a command line client in
place of an API. Often times you want to do things like output stdout
within your own logs or ensure that every time the command is called,
there are a standard set of flags that are used.

## Example

Let's use the `knife` command for working wih a
[Chef](https://chef.io) server. The `knife` command can use a common
set of arguments that are required to be in certain order. It also can
output JSON, which is typical in scripting interactions with the
server.

Lets create a new Procs factory to work with the `knife` command.

```
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
```

The `procs.Builder` allows creating simple templates and data for
building commands. The templates allow simple variable expansion as
you might find in a bash. The `procs.New` function creates a new
`*exec.Cmd` after safely parsing the command and expanding any env
vars.

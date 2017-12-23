/*
The procs package provides some higher level abstractions on top of
the exec package.

The majority of this functionality is included in the procs.Process struct.

Defining a Command

A command can be defined by a string rather than a []string. Normally,
this also implies that the library will run the command in a shell,
exposing a potential man in the middle attack. Rather than using a
shell, procs lexically parses the command for the different
arguments. It also allows for pipes in order to string commands
together.

	p := procs.NewProcess("kubectl get events | grep dev")


Output Handling

One use case that is cumbersome is using the piped output from a
command. For example, lets say we wanted to start a couple commands
and have each command have its own prefix in stdout, while still
capturing the output of the command as-is.

	p := procs.NewProcess("cmd1")
	p.OutputHandler = func(line string) string {
		fmt.Printf("cmd1 | %s\n")
		return line
	}

	out, _ := p.Run()
	fmt.Println(out)

Whatever is returned from the OutputHandler will be in the buffered
output. In this way you can choose to filter or skip output buffering
completely.

Environment Variables

Rather than use the exec.Cmd []string environment variables, a
procs.Process uses a map[string]string for environment variables.

	p := procs.NewProcess("echo $FOO")
	p.Env = map[string]string{"FOO": "foo"}

Also, environment variables will be expanded automatically using the
os.Expand semantics and the provided environment. If no environment is
provided, the parent process environment is used.

There is also a Env function that can help to merge an existing
environment with the parent environment to allow overriding parent
variables.

	myenv := map[string]string{"USER": "foo"}
	env := ParseEnv(Env(myenv, true))

The Env function will overlay the passed in environment with the
parent environment and ParseEnv takes a environment []string and
converts it to a map[string]string for use with a Proc.
*/

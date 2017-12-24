# Procs

[![](https://travis-ci.org/ionrock/procs.svg?branch=master)](https://travis-ci.org/ionrock/procs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ionrock/procs)](https://goreportcard.com/report/github.com/ionrock/procs)
[![GoDoc](https://godoc.org/github.com/ionrock/procs?status.svg)](https://godoc.org/github.com/ionrock/procs)

Procs is a library to make working with command line applications a
little nicer.

The primary use case is when you have to use a command line client in
place of an API. Often times you want to do things like output stdout
within your own logs or ensure that every time the command is called,
there are a standard set of flags that are used.

For example, I'm automating [Chef](https://chef.io) usage within CI
using the [knife]() tool rather than reimplement and maintain an API
client in client in Go. This also ensures that when things go wrong,
the commands can be debugged directly.

## Basic Usage

The majority of this functionality is intended to be included the
procs.Process.

### Defining a Command

A command can be defined by a string rather than a []string. Normally,
this also implies that the library will run the command in a shell,
exposing a potential man in the middle attack. Rather than using a
shell, procs lexically parses the command for the different
arguments. It also allows for pipes in order to string commands
together.

```go
p := procs.NewProcess("kubectl get events | grep dev")
```

### Output Handling

One use case that is cumbersome is using the piped output from a
command. For example, lets say we wanted to start a couple commands
and have each command have its own prefix in stdout, while still
capturing the output of the command as-is.

```go
p := procs.NewProcess("cmd1")
p.OutputHandler = func(line string) string {
	fmt.Printf("cmd1 | %s\n")
	return line
}
out, _ := p.Run()
fmt.Println(out)
```

Whatever is returned from the `OutputHandler` will be in the buffered
output. In this way you can choose to filter or skip output buffering
completely.

### Environment Variables

Rather than use the `exec.Cmd` `[]string` environment variables, a
`procs.Process` uses a `map[string]string` for environment variables.

```go
p := procs.NewProcess("echo $FOO")
p.Env = map[string]string{"FOO": "foo"}
```

Also, environment variables will be expanded automatically using the
`os.Expand` semantics and the provided environment. If no environment
is provided, the parent process environment is used.

There is also a `Env` function that can help to merge an existing
environment with the parent environment to allow overriding parent
variables.

```go
myenv := map[string]string{"USER": "foo"}
env := ParseEnv(Env(myenv, true))
```

The `Env` function will overlay the passed in environment with the
parent environment and `ParseEnv` takes a environment `[]string` and
converts it to a `map[string]string` for use with a Proc.

## Example Applications

Take a look in the [`cmd`](./cmd/) dir for some simple applications
that use the library.

### Prelog

The `prelog` command allows running a command and prefixing the output
with a value.

```bash
$ prelog -prefix foo -- echo 'hello world!'
Running the command
foo ===> hello world!
Accessing the output without a prefix.
hello world!
Running the command with Start / Wait
foo ===> hello world!
```

### Cmdtmpl

The `cmdtmpl` command uses the `procs.Builder` to create a command
based on some paramters. It will take a `data.yml` file and
`template.yml` file to create a command.

```bash
$ cat example/data.json
{
  "source": "https://my.example.org",
  "user": "foo",
  "model": "widget",
  "action": "create",
  "args": "-f new -i improved"
}
$ cat example/template.json
[
  "mysvc ${model} ${action} ${args}",
  "--endpoint ${source}",
  "--username ${user}"
]
$ cmdtmpl -data example/data.json -template example/template.json
Command: mysvc foo widget create -f new -i imporoved --endpoint https://my.example.org --username foo
$ cmdtmpl -data example/data.json -template example/template.json -field user=bar
Command: mysvc foo widget create -f new -i imporoved --endpoint https://my.example.org --username bar
```

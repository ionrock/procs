# Procs

![](https://travis-ci.org/ionrock/procs.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ionrock/procs)](https://goreportcard.com/report/github.com/ionrock/procs)
![GoDoc](https://godoc.org/github.com/ionrock/procs?status.svg)

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

The essential usage is to use the `procs.Process` in place of `exec.Cmd`.

```go
proc := procs.NewProcess("knife search -i node \"role:foo\" | grep test")
out, err := proc.Run()
if err != nil {
	panic(err)
}

fmt.Print(out)
```

The pipes are connected within the `procs.Process` type and no shell is used.

You can also work with the environment using a `map[string]string`
rather than a `[]string`.

```go
proc.Env = map[string]string{"FOO": "bar"}
```

There are helpers as well to allow mixing the initial environment with extra values.

```go
// Parse the []string to a map[string]string
env := procs.ParseEnv(os.Environ())
env["FOO"] = "bar"
proc.Env = env
```

## Examples

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

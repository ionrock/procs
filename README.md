# Procs

Procs is a library to make working with command line applications more
like working with an API.

The primary use case is when you have to use a command line client in
place of an API. Often times you want to do things like output stdout
within your own logs or ensure that every time the command is called,
there are a standard set of flags that are used.

## Examples

Take a look in the [`cmd`](./cmd/) dir for some simple applications
that use the library.

### Prelog

The `prelog` command allows running a command and prefixing the output
with a value.

```
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

```
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
$ cmdtmpl -data example/data.yml -template example/template.yml
Command: mysvc foo widget create -f new -i imporoved --endpoint https://my.example.org --username foo
$ cmdtmpl -data example/data.yml -template example/template.yml -field user=bar
Command: mysvc foo widget create -f new -i imporoved --endpoint https://my.example.org --username bar
```

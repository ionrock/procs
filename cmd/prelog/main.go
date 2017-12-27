package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ionrock/procs"
)

type prefixer struct {
	Prefix string
}

func (p *prefixer) handler(line string) string {
	// We can output to stdout with our prefix
	fmt.Printf("%s | %s\n", p.Prefix, line)

	// By returning the line as-is we keep the original output as-is.
	// This also allows for avoiding the buffering of the output by
	// returning an empty string.
	return line
}

func prefixHandler(prefix string) procs.OutHandler {
	return func(line string) string {
		// We can output to stdout with our prefix
		fmt.Printf("%s | %s\n", prefix, line)

		// By returning the line as-is we keep the original output as-is.
		// This also allows for avoiding the buffering of the output by
		// returning an empty string.
		return line
	}
}

func main() {
	prefix := flag.String("prefix", "prelog", "log prefix")

	flag.Parse()

	// We'll concat our command line args to create a single string
	// for our command.
	command := strings.Join(flag.Args(), " ")

	// The procs.New returns a *procs.Proc. The command is safely
	// parsed and will expand any variables found in the environment
	// using shell syntax ($ or ${} using os.ExpandEnv)
	cmd := procs.NewProcess(command)

	// Add an OutputHandler for adding our prefix.
	cmd.OutputHandler = prefixHandler(*prefix)

	// Run our command. This will pipe the output to stdout prefixed
	// with the name and delimiter defined in the Output.
	fmt.Println("Running the command")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error: %q\n", err)
		os.Exit(1)
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("error: %q\n", err)
		os.Exit(1)
	}

	// We can access the output as a normal string as well.
	fmt.Println("Accessing the output without a prefix.")
	fmt.Println(string(out))

	// You can also use Start / Wait and reuse the OutputHandler function
	cmd = procs.NewProcess(command)
	cmd.OutputHandler = prefixHandler(*prefix)

	fmt.Println("Running the command with Start / Wait")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	cmd.Wait()
}

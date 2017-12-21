package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ionrock/procs"
)

func main() {
	prefix := flag.String("prefix", "prelog", "log prefix")

	flag.Parse()

	// We'll concat our command line args to create a single string
	// for our command.
	command := strings.Join(flag.Args(), " ")

	// The procs.New returns a *procs.Proc. The command is safely
	// parsed and will expand any variables found in the environment
	// using shell syntax ($ or ${} using os.ExpandEnv)
	cmd := procs.New(command)

	// Add an OutputHandler for adding our prefix.
	cmd.OutputHandler = func(line string) string {
		return fmt.Sprintf("%s | %s", *prefix, line)
	}

	// Run our command. This will pipe the output to stdout prefixed
	// with the name and delimiter defined in the Output.
	fmt.Println("Running the command")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	// We can access the output as a normal string as well.
	fmt.Println("Accessing the output without a prefix.")
	fmt.Println(out.Output())

	// You can also use Start / Wait and reuse the Output instance.
	cmd = procs.New(command)
	cmd.Output = out

	fmt.Println("Running the command with Start / Wait")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	cmd.Wait()
}

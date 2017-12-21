package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ionrock/procs"
)

type Prefixer struct {
	Prefix string
}

func (p *Prefixer) Handler(line string) string {
	return fmt.Sprintf("%s | %s", p.Prefix, line)
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
	prefixer := &Prefixer{*prefix}
	cmd.OutputHandler = prefixer.Handler

	// Run our command. This will pipe the output to stdout prefixed
	// with the name and delimiter defined in the Output.
	fmt.Println("Running the command")
	out, err := cmd.Run()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	// We can access the output as a normal string as well.
	fmt.Println("Accessing the output without a prefix.")
	fmt.Println(out)

	// You can also use Start / Wait and reuse the OutputHandler function
	cmd = procs.NewProcess(command)
	cmd.OutputHandler = prefixer.Handler

	fmt.Println("Running the command with Start / Wait")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	cmd.Wait()
}

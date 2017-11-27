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

	out := &procs.Output{
		Name:      *prefix,
		Capture:   true,
		Delimiter: "===>", // default is "|"
	}

	command := strings.Join(flag.Args(), " ")
	cmd := procs.New(command)
	cmd.Output = out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	fmt.Println("The output was")
	fmt.Println(out.Output())

	cmd = procs.New(command)
	cmd.Output = out

	err = cmd.Start()
	if err != nil {
		fmt.Printf("error: %q\n", err)
	}

	cmd.Wait()
}

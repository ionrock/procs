package procs_test

import (
	"fmt"

	"github.com/ionrock/procs"
)

func Example() {

	b := procs.Builder{
		Context: map[string]string{
			"NAME": "eric",
		},
		Templates: []string{
			"echo $NAME |",
			"grep $NAME",
		},
	}

	p := procs.NewProcess(b.Command())

	p.Run()
	out, _ := p.Output()
	fmt.Println(string(out))
	// Output:
	// eric
}

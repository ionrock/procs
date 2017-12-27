package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ionrock/procs"
)

func loadData(path string) (map[string]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := make(map[string]string)

	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func prefixHandler(name string) procs.OutHandler {
	return func(line string) string {
		fmt.Printf("%s | %s\n", name, line)
		return ""
	}
}

func main() {
	procfilePath := flag.String("procfile", "p", "a file with key/value JSON to act like a Procfile")

	flag.Parse()

	procfile, err := loadData(*procfilePath)
	if err != nil {
		fmt.Printf("error loading procfile data: %q\n", err)
		os.Exit(1)
	}

	m := procs.NewManager()
	for k, v := range procfile {
		outHandler := prefixHandler(k)
		outHandler(fmt.Sprintf("Starting %s with %s", k, v))

		p := procs.NewProcess(v)
		p.OutputHandler = outHandler
		p.ErrHandler = outHandler

		m.StartProcess(k, p)
	}
	m.Wait()

	fmt.Println("All processes exited. Exiting...")
}

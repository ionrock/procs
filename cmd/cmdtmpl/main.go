package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func loadTmpls(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmpls := []string{}

	err = json.Unmarshal(content, &tmpls)
	if err != nil {
		return nil, err
	}

	return tmpls, nil
}

type fieldArray []string

func (f *fieldArray) String() string {
	return strings.Join(*f, " ")
}

func (f *fieldArray) Set(val string) error {
	*f = append(*f, val)
	return nil
}

func main() {
	data := flag.String("data", "", "a set of key/value JSON for passing to a template")
	tmpl := flag.String("template", "", "a list of strings that have shell style templates")

	var fields fieldArray
	flag.Var(&fields, "field", "field in the form of KEY=VALUE to overide a")

	flag.Parse()

	ctx, err := loadData(*data)
	if err != nil {
		fmt.Printf("error loading data: %q\n", err)
		os.Exit(1)
	}

	tmpls, err := loadTmpls(*tmpl)
	if err != nil {
		fmt.Printf("error loading templates: %q\n", err)
		os.Exit(1)
	}

	b := procs.Builder{
		Context:   ctx,
		Templates: tmpls,
	}

	if len(fields) > 0 {
		ctx := procs.ParseEnv(fields)
		fmt.Printf("Command: %s\n", b.CommandContext(ctx))
	} else {
		fmt.Printf("Command: %s\n", b.Command())
	}

}

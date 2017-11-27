package procs_test

import (
	"testing"

	"github.com/ionrock/procs"
)

func TestCommand(t *testing.T) {
	b := &procs.Builder{
		Context: map[string]string{
			"options": "-Fj -s https://example.com/chef -k knife.pem",
		},

		Templates: []string{
			"knife",
			"${model} ${action}",
			"${args}",
			"${options}",
		},
	}

	cmd := b.Command(map[string]string{
		"model":  "data bag",
		"action": "from file",
		"args":   "foo data_bags/foo/bar.json",
	})

	expected := "knife data bag from file foo data_bags/foo/bar.json -Fj -s https://example.com/chef -k knife.pem"
	if cmd != expected {
		t.Fatalf("failed building command: %q != %q", cmd, expected)
	}
}

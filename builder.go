package procs

import (
	"fmt"
	"os"
	"strings"
)

// Builder helps construct commands using templates.
type Builder struct {
	Context   map[string]string
	Templates []string
}

func (b *Builder) getConfig(ctx map[string]string) func(string) string {
	return func(key string) string {
		if v, ok := ctx[key]; ok {
			return v
		}

		return ""
	}

}

func (b *Builder) expand(v string, ctx map[string]string) string {
	return os.Expand(v, b.getConfig(ctx))
}

func (b *Builder) Command(ctx map[string]string) string {
	for k, v := range b.Context {
		ctx[k] = b.expand(v, ctx)
	}

	parts := []string{}
	for _, t := range b.Templates {
		fmt.Println(t)
		parts = append(parts, b.expand(t, ctx))
	}

	fmt.Println(parts)
	return strings.Join(parts, " ")
}

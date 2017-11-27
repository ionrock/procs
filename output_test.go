package procs_test

import (
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/ionrock/procs"
)

func TestOutputFactory(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("capture", func() {
		o.BeforeEach(func(t *testing.T) (*testing.T, *procs.Output) {
			return t, &procs.Output{Padding: 5, Capture: true}
		})

		o.Spec("writing a line stores the output", func(t *testing.T, of *procs.Output) {
			of.WriteLine("foo", "hello world", false)
			Expect(t, of.Output()).To(Equal("hello world\n"))
		})
	})
}

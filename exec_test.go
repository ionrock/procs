package procs_test

import (
	"os"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/ionrock/procs"
)

func TestNewCmd(t *testing.T) {
	newCmdTests := []struct {
		cmd   *procs.Proc
		parts []string
	}{
		// Make sure we split more than one part
		{procs.New("/bin/echo 'hello world'"), []string{"/bin/echo", "hello world"}},

		// Make sure we only call with a path argument
		{procs.New("/bin/echo"), []string{"/bin/echo"}},
	}

	for _, cmdTest := range newCmdTests {
		for i, part := range cmdTest.parts {
			Expect(t, part).To(Equal(cmdTest.cmd.Cmd.Args[i]))
		}
	}
}

func TestEnv(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) (*testing.T, string) {
		envvar := "PROC_ENVVAR_FALLBACK_TEST"
		os.Setenv("PROC_ENVVAR_FALLBACK_TEST", "bar")
		return t, envvar
	})

	o.Spec("include env", func(t *testing.T, envvar string) {
		env := map[string]string{
			"foo": "hello world",
		}
		environ := procs.ParseEnv(procs.Env(env, true))

		Expect(t, environ).To(HaveKey(envvar))
		Expect(t, environ).To(HaveKey("foo"))
	})
}

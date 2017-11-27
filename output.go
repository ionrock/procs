// Most of this code is taken from forego, with the colorization removed.
// https://github.com/ddollar/forego
package procs

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Output provdes a helper to watch an io.Reader from a process
// for writing to stdout with a prefix.
type Output struct {
	Name      string
	Delimiter string
	Padding   int
	Capture   bool

	pipesWait    *sync.WaitGroup
	outputBuffer *bytes.Buffer
	sync.Mutex
}

// LineReader is intended to be used for the stdout/err pipe from a
// process to prefix the output.
func (of *Output) LineReader(wg *sync.WaitGroup, name string, r io.Reader, isError bool) {
	defer wg.Done()

	reader := bufio.NewReader(r)

	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)

		if n, err := reader.Read(buf); err != nil {
			return
		} else {
			buf = buf[:n]
		}

		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				break
			}
			buffer.Write(buf[0:i])
			of.WriteLine(name, buffer.String(), isError)
			buffer.Reset()
			buf = buf[i+1:]
		}

		buffer.Write(buf)
	}
}

func (of *Output) SystemOutput(str string) {
	of.WriteLine("forego", str, false)
}

func (of *Output) ErrorOutput(str string) {
	fmt.Printf("ERROR: %s\n", str)
	os.Exit(1)
}

// WriteLine writes out a single with a prefix.
func (of *Output) WriteLine(left, right string, isError bool) {
	of.Lock()
	defer of.Unlock()

	delimiter := "|"
	if of.Delimiter != "" {
		delimiter = of.Delimiter
	}

	formatter := fmt.Sprintf("%%-%ds %s ", of.Padding, delimiter)
	fmt.Printf(formatter, left)
	fmt.Println(right)
	if of.Capture {
		if of.outputBuffer == nil {
			of.outputBuffer = bytes.NewBufferString(right + "\n")
		} else {
			of.outputBuffer.WriteString(right + "\n")
		}
	}
}

// Output returns the contexts of the output as a list of strings
// where each item is a line.
func (of *Output) Output() string {
	return strings.TrimSpace(of.outputBuffer.String())
}

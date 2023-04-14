package execmd

import (
	"bytes"
	"io"
	"log"
	"strings"
)

// prefixedStream is a custom writer that wraps a logger and allows
// adding a prefix to each line of output, as well as optionally saving
// the output data to a buffer.
type prefixedStream struct {
	Logger   *log.Logger
	buffer   *bytes.Buffer
	data     *bytes.Buffer
	prefix   string
	saveData bool
}

// newPrefixedStream creates a new PrefixedStream with the provided logger,
// prefix, and a flag indicating whether to save the output data.
func newPrefixedStream(logger *log.Logger, prefix string, saveData bool) *prefixedStream {
	return &prefixedStream{
		Logger:   logger,
		buffer:   bytes.NewBuffer(nil),
		prefix:   prefix,
		saveData: saveData,
		data:     bytes.NewBuffer(nil),
	}
}

// Write writes data to the buffer and processes any complete lines,
// adding the prefix and logging them.
func (p *prefixedStream) Write(data []byte) (int, error) {
	n, err := p.buffer.Write(data)
	if err != nil {
		return n, err
	}

	return n, p.outputLines()
}

// Get retrieves the data buffer.
func (p *prefixedStream) Get() *bytes.Buffer {
	return p.data
}

// Close closes the buffer and flushes any remaining data.
func (p *prefixedStream) Close() error {
	if err := p.flush(); err != nil {
		return err
	}

	p.buffer = bytes.NewBuffer(nil)
	return nil
}

// flush reads the remaining data from the buffer, processes it,
// and logs it with the prefix.
func (p *prefixedStream) flush() error {
	var data []byte

	if _, err := p.buffer.Read(data); err != nil {
		return err
	}

	p.output(string(data))
	return nil
}

// outputLines reads and processes lines from the buffer, logging them
// with the prefix. Any incomplete lines are left in the buffer.
func (p *prefixedStream) outputLines() error {
	for {
		line, err := p.buffer.ReadString('\n')

		if len(line) > 0 {
			if strings.HasSuffix(line, "\n") {
				p.output(line)
			} else {
				if _, err := p.buffer.WriteString(line); err != nil {
					return err
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// ClearData resets the data buffer.
func (p *prefixedStream) ClearData() {
	p.data.Reset()
}

// output logs the given text with the prefix and, if saveData is true,
// appends the text to the data buffer.
func (p *prefixedStream) output(text string) {
	if len(text) < 1 {
		return
	}

	if p.saveData {
		p.data.WriteString(text)
	}

	text = p.prefix + text

	p.Logger.Print(text)
}

package execmd

// main idea taken from logstreamer
// see: https://github.com/kvz/logstreamer

import (
	"bytes"
	"io"
	"log"
	"strings"
)

type pStream struct {
	Logger   *log.Logger
	buf      *bytes.Buffer
	prefix   string
	saveData bool
	data     *bytes.Buffer
}

func newPStream(logger *log.Logger, prefix string, saveData bool) *pStream {
	return &pStream{
		Logger:   logger,
		buf:      bytes.NewBuffer([]byte("")),
		prefix:   prefix,
		saveData: saveData,
		data:     bytes.NewBuffer([]byte("")),
	}
}

func (p *pStream) Write(b []byte) (n int, err error) {
	if n, err = p.buf.Write(b); err != nil {
		return
	}

	err = p.OutputLines()
	return
}

func (p *pStream) Close() error {
	if err := p.Flush(); err != nil {
		return err
	}

	p.buf = bytes.NewBuffer([]byte(""))
	return nil
}

func (p *pStream) Flush() error {
	var b []byte

	if _, err := p.buf.Read(b); err != nil {
		return err
	}

	p.out(string(b))
	return nil
}

func (p *pStream) OutputLines() error {
	for {
		line, err := p.buf.ReadString('\n')

		if len(line) > 0 {
			if strings.HasSuffix(line, "\n") {
				p.out(line)
			} else {
				// put back into buffer, it's not a complete line yet
				//  Close() or Flush() have to be used to flush out
				//  the last remaining line if it does not end with a newline
				if _, err := p.buf.WriteString(line); err != nil {
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

func (p *pStream) FlushData() {
	p.data.Reset()
}

func (p *pStream) Get() *bytes.Buffer {
	return p.data
}

func (p *pStream) out(str string) {
	if len(str) < 1 {
		return
	}

	if p.saveData {
		p.data.WriteString(str)
	}

	str = p.prefix + str

	p.Logger.Print(str)
}

package process

import (
	"io/ioutil"
	"os"
	"github.com/glycerine/rbuf"
)

type LogWriter interface {
	Write(p []byte) (n int, err error)
	String() string
	Len() int64
	Close()
}

type FileLogWriter struct {
	filename string
	file     *os.File
}

func NewFileLogWriter(file string) (*FileLogWriter, error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}

	flw := &FileLogWriter{
		filename: file,
		file:     f,
	}
	return flw, nil
}

func (flw FileLogWriter) Close() {
	flw.file.Close()
}

func (flw FileLogWriter) Write(p []byte) (n int, err error) {
	return flw.file.Write(p)
}

func (flw FileLogWriter) String() string {
	b, err := ioutil.ReadFile(flw.filename)
	if err == nil {
		return string(b)
	}
	return ""
}

func (flw FileLogWriter) Len() int64 {
	s, err := os.Stat(flw.filename)
	if err == nil {
		return s.Size()
	}
	return 0
}

type InMemoryLogWriter struct {
	buffer *rbuf.FixedSizeRingBuf
}

func NewInMemoryLogWriter() InMemoryLogWriter {
	imlw := InMemoryLogWriter{}
	imlw.buffer = rbuf.NewFixedSizeRingBuf(1024 * 1024 * 2) // 2M size
	return imlw
}

func (imlw InMemoryLogWriter) Write(p []byte) (n int, err error) {
	return imlw.buffer.Write(p)
}

func (imlw InMemoryLogWriter) String() string {
	return string(imlw.buffer.Bytes())
}

func (imlw InMemoryLogWriter) Len() int64 {
	return int64(imlw.buffer.ContigLen())
}

func (imlw InMemoryLogWriter) Close() {
}

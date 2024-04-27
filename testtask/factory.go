package testtask

import (
	"io"
	"os"
)

func NewStringEntryFactory(parser Parser[io.Reader, *StringEntry]) StringEntryFactory {
	return StringEntryFactory{
		parser: parser,
	}
}

type StringEntryFactory struct {
	parser Parser[io.Reader, *StringEntry]
}

func (s StringEntryFactory) Create(r io.Reader) Entry {
	e := &StringEntry{}
	e, err := s.parser.Parse(r, e)
	if err != nil {
		return &StringEntry{err: err}
	}
	return e
}

type FileEntryFactory struct {
}

func (f FileEntryFactory) Create(filename string) OSEntry[io.Reader] {
	file, err := os.Open(filename)
	if err != nil {
		return &FileEntry{err: err}
	}
	return &FileEntry{readCloser: file, err: nil}
}

package testtask

import (
	"io"
	"time"
)

type StringEntry struct {
	message   string
	timestamp time.Time
	err       error
}

func (s *StringEntry) String() string {
	return "some"
}

func (s *StringEntry) Err() error {
	return s.err
}

type FileEntry struct {
	readCloser io.ReadCloser
	err        error
}

func (f *FileEntry) Val() io.Reader {
	return f.readCloser
}

func (f *FileEntry) Err() error {
	return f.err
}

func (f *FileEntry) Close() {
	f.readCloser.Close()
}

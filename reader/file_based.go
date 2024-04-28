package reader

import (
	"bufio"
	"os"
)

type FileBased struct {
	reader *bufio.Reader
	file   *os.File
	delim  byte
}

func (f *FileBased) Read() (string, error) {
	return f.reader.ReadString(f.delim)
}

func (f *FileBased) Close() {
	f.file.Close()
}

func NewFileBasedFactory() FileBasedFactory {
	return FileBasedFactory{}
}

type FileBasedFactory struct {
}

func (f FileBasedFactory) Create(delim byte) func(name string) (Reader[string], error) {
	return func(name string) (Reader[string], error) {
		file, err := os.Open(name)
		if err != nil {
			return nil, err
		}

		return &FileBased{
			file:   file,
			reader: bufio.NewReader(file),
			delim:  delim,
		}, nil
	}
}

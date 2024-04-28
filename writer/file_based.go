package writer

import (
	"bufio"
	"fmt"
	"os"
)

type FileBased struct {
	writer *bufio.Writer
	delim  byte
	file   *os.File
}

func (f *FileBased) Write(content string) error {
	fmt.Println(f.file.Name(), " : ", content)
	_, err := f.writer.WriteString(content + string(f.delim))
	return err
}

func (f *FileBased) Close() {
	fmt.Println(f.file.Name(), " : write")
	f.writer.Flush()
}

func NewFileBasedFactory() FileBasedFactory {
	return FileBasedFactory{}
}

type FileBasedFactory struct {
}

func (f FileBasedFactory) Create(delim byte) func(name string) (Writer[string], error) {
	return func(name string) (Writer[string], error) {
		file, err := os.Create(name)
		if err != nil {
			return nil, err
		}

		return &FileBased{
			writer: bufio.NewWriter(file),
			delim:  delim,
			file:   file,
		}, nil
	}
}

package source

import (
	"github.com/Alieksieiev0/test-task/reader"
	"github.com/Alieksieiev0/test-task/testtask"
)

func NewFileContentFactory() FileContentFactory {
	return FileContentFactory{}
}

type FileContentFactory struct {
}

func (f FileContentFactory) Create(reader reader.Reader[string]) testtask.Source[string] {
	return &FileContent{
		reader: reader,
	}
}

func NewAsyncFileContentFactory() AsyncFileContentFactory {
	return AsyncFileContentFactory{}
}

type AsyncFileContentFactory struct {
}

func (f AsyncFileContentFactory) Create(message <-chan string) testtask.Source[string] {
	return &AsyncFileContent{
		message: message,
	}
}

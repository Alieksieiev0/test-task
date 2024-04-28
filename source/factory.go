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

func NewSortedFilePairsFactory() SortedFilePairsFactory {
	return SortedFilePairsFactory{}
}

type SortedFilePairsFactory struct {
}

func (s SortedFilePairsFactory) CreateFromMap(files map[string]string) *SortedFilePairs {
	var keys []string
	var values []string
	for k, v := range files {
		keys = append(keys, k)
		values = append(values, v)
	}

	return &SortedFilePairs{
		keys:   NewFileCollection(keys),
		values: NewFileCollection(values),
	}
}

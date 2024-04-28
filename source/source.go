package source

import (
	"fmt"
	"io"

	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/reader"
	"go.uber.org/atomic"
	"golang.org/x/exp/maps"
)

func NewFileCollection(files map[string]string) *FileCollection {
	return &FileCollection{
		files:  maps.Keys(files),
		offset: atomic.NewInt64(0),
	}
}

type FileCollection struct {
	files  []string
	offset *atomic.Int64
}

func (f *FileCollection) Data() iterator.Iterator[string] {
	return iterator.NewCallbackBased(func() *iterator.Step[string] {
		curr := f.offset.Load()
		if curr == int64(len(f.files)) {
			return iterator.NewStepErr[string](io.EOF)
		}
		f.offset.Inc()
		return iterator.NewStepVal(f.files[curr])
	})
}

func NewFileContent() *FileContent {
	return &FileContent{}
}

type FileContent struct {
	reader reader.Reader[string]
}

func (f *FileContent) Data() iterator.Iterator[string] {
	return iterator.NewCallbackBased(func() *iterator.Step[string] {
		if f.reader == nil {
			return iterator.NewStepErr[string](fmt.Errorf("file contents were not provided"))
		}
		content, err := f.reader.Read()
		if err != nil {
			return iterator.NewStepErr[string](err)
		}
		return iterator.NewStepVal(content)
	})
}

type AsyncFileContent struct {
	message <-chan string
}

func (f *AsyncFileContent) Data() iterator.Iterator[string] {
	return iterator.NewCallbackBased(func() *iterator.Step[string] {
		if f.message == nil {
			return iterator.NewStepErr[string](fmt.Errorf("file contents were not provided"))
		}
		content, ok := <-f.message
		if !ok {
			return iterator.NewStepErr[string](io.EOF)
		}
		return iterator.NewStepVal(content)
	})
}

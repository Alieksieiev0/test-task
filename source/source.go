package source

import (
	"bufio"
	"io"

	"github.com/Alieksieiev0/test-task/iterator"
	"go.uber.org/atomic"
)

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

type FileContent struct {
	r     *bufio.Reader
	delim byte
}

func (f *FileContent) Load(file io.Reader) {
	f.r = bufio.NewReader(file)
}

func (f *FileContent) Data() iterator.Iterator[string] {
	return iterator.NewCallbackBased(func() *iterator.Step[string] {
		content, err := f.r.ReadString(f.delim)
		if err != nil {
			return iterator.NewStepErr[string](err)
		}
		return iterator.NewStepVal(content)
	})
}

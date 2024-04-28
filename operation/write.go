package operation

import (
	"io"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

func NewAsyncWriteOperation[Input Entry[error], Output AsyncEntry[error]]() *AsyncWriteOperation[Input, Output] {
	ctx, close := context.WithCancel(context.Background())
	return &AsyncWriteOperation[Input, Output]{ctx: ctx, close: close}
}

type AsyncWriteOperation[Input Entry[error], Output AsyncEntry[error]] struct {
	ctx   context.Context
	close func()
}

func (a *AsyncWriteOperation[Input, Output]) Close() error {
	a.close()
	return nil
}

func (r *AsyncWriteOperation[Input, Output]) Run(
	iterator iterator.Iterator[Input],
	results Output,
) {
	for {
		select {
		case <-r.ctx.Done():
			results.PassVal(io.EOF)
			results.Close()
			return
		default:
			entry := iterator.Next()
			if entry.Err() != nil {
				results.PassVal(entry.Err())
				results.Close()
				return
			}
		}
	}
}

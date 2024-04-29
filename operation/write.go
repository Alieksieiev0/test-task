package operation

import (
	"io"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

type AsyncWriteOperation[Input Entry[error], Output AsyncEntry[error]] struct {
	ctx context.Context
}

func (a *AsyncWriteOperation[Input, Output]) Run(
	iterator iterator.Iterator[Input],
	results Output,
) {
	for {
		select {
		case <-a.ctx.Done():
			results.Close()
			return
		default:
			entry := iterator.Next()
			err := entry.Err()
			if err == nil {
				err = entry.Val().Val()
			}

			if err == io.EOF {
				results.Close()
				return
			} else if err != nil {
				results.PassVal(entry.Err())
				results.Close()
				return
			}

		}
	}
}

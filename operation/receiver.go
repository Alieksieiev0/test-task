package operation

import (
	"io"
	"log"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

type ReadOperation[T Entry] struct {
	ctx context.Context
}

func (r *ReadOperation[T]) Run(iterator iterator.Iterator[T]) error {
	for {
		select {
		case <-r.ctx.Done():
			return io.EOF
		default:
			entry := iterator.Next()
			if entry.Err() != nil {
				return entry.Err()
			}
			log.Print("sink: ", entry.Val().String())
		}
	}
}

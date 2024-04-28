package operation

import (
	"fmt"
	"io"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

func NewAsyncReadOperation[Input ErrorProneEntry[string], Output AsyncErrorProneEntry[string]]() *AsyncReadOperation[Input, Output] {
	ctx, close := context.WithCancel(context.Background())
	return &AsyncReadOperation[Input, Output]{ctx: ctx, close: close}
}

type AsyncReadOperation[Input ErrorProneEntry[string], Output AsyncErrorProneEntry[string]] struct {
	ctx   context.Context
	close func()
}

func (a *AsyncReadOperation[Input, Output]) Close() error {
	a.close()
	return nil
}

func (a *AsyncReadOperation[Input, Output]) Run(iterator iterator.Iterator[Input], results Output) {
	fmt.Println("---read-start---")
	for {
		select {
		case <-a.ctx.Done():
			results.Close()
			return
		default:
			entry := iterator.Next()

			if entry.Err() == io.EOF {
				results.Close()
				return
			} else if entry.Err() != nil {
				results.PassErr(entry.Err())
				results.Close()
				return
			}
			results.PassVal(entry.Val().Val())
		}
	}
}

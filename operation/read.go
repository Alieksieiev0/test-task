package operation

import (
	"fmt"
	"io"
	"time"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

func NewAsyncReadOperation[Input ErrorProneEntry[string], Output AsyncErrorProneEntry[string]](
	delay time.Duration,
) *AsyncReadOperation[Input, Output] {
	ctx, close := context.WithCancel(context.Background())
	return &AsyncReadOperation[Input, Output]{ctx: ctx, delay: delay, close: close}
}

type AsyncReadOperation[Input ErrorProneEntry[string], Output AsyncErrorProneEntry[string]] struct {
	ctx   context.Context
	delay time.Duration
	close func()
}

func (a *AsyncReadOperation[Input, Output]) Close() error {
	a.close()
	return nil
}

func (a *AsyncReadOperation[Input, Output]) Run(iterator iterator.Iterator[Input], results Output) {
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

			time.Sleep(a.delay)
			line := entry.Val().Val()
			fmt.Println("LINE: ", line)
			results.PassVal(line)
		}
	}
}

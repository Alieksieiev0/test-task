package operation

import (
	"fmt"
	"sync"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

func NewAsyncFileDecoder[T AsyncErrorProneEntry[string], V AsyncEntry[error]](
	wg *sync.WaitGroup,
) *AsyncFileDecoder[T, V] {
	ctx, cancel := context.WithCancel(context.Background())
	return &AsyncFileDecoder[T, V]{ctx: ctx, cancel: cancel, wg: wg}
}

type AsyncFileDecoder[T AsyncErrorProneEntry[string], V AsyncEntry[error]] struct {
	ctx    context.Context
	cancel func()
	wg     *sync.WaitGroup
}

func (a *AsyncFileDecoder[T, V]) Close() {
	a.cancel()
}

func (a *AsyncFileDecoder[T, V]) Run(
	readCallbackIterator iterator.Iterator[func(context.Context) T],
	writeCallbackIterator iterator.Iterator[func(context.Context, T) V],
) {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		readCallbackEntry := readCallbackIterator.Next()
		if readCallbackEntry.Err() != nil {
			break
		}

		writeCallbackEntry := writeCallbackIterator.Next()
		if writeCallbackEntry.Err() != nil {
			break
		}

		readResults := readCallbackEntry.Val()(ctx)
		writeResults := writeCallbackEntry.Val()(ctx, readResults)

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			for {
				select {
				case <-a.ctx.Done():
					cancel()
					return
				case err, ok := <-readResults.Err():
					if !ok {
						continue
					}
					fmt.Println("ERROR READING FILE ", err)
					cancel()
					return
				case err, more := <-writeResults.Val():
					if !more {
						return
					}
					cancel()
					fmt.Println("ERROR SAVING FILE:", err)
					return
				}

			}
		}()
	}
	a.wg.Wait()
}

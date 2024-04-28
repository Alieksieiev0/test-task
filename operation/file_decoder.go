package operation

import (
	"fmt"
	"sync"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

func NewAsyncFileDecoder[T AsyncErrorProneEntry[string], U func(AsyncEntry[string]) AsyncEntry[error]]() *AsyncFileDecoder[T, U] {
	ctx, close := context.WithCancel(context.Background())
	return &AsyncFileDecoder[T, U]{ctx: ctx, close: close}
}

type AsyncFileDecoder[T AsyncErrorProneEntry[string], U func(AsyncEntry[string]) AsyncEntry[error]] struct {
	ctx   context.Context
	close func()
}

func (a *AsyncFileDecoder[T, U]) Close() error {
	a.close()
	return nil
}

func (a *AsyncFileDecoder[T, U]) Run(
	readIterator iterator.Iterator[T],
	writeCallbackIterator iterator.Iterator[U],
) {
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		readEntry := readIterator.Next()
		readResults := readEntry.Val()
		if readEntry.Err() != nil {
			// nil check?
			readResults.Close()
			wg.Done()
			break
		}

		writeCallbackEntry := writeCallbackIterator.Next()
		if writeCallbackEntry.Err() != nil {
			wg.Done()
			break
		}
		writeResults := writeCallbackEntry.Val()(readResults)

		go func() {
			for {
				select {
				case err, ok := <-readResults.Err():
					if !ok {
						continue
					}
					fmt.Println("ERROR READING FILE ", err)
					wg.Done()
					return
				case err, more := <-writeResults.Val():
					if !more {
						fmt.Println("NO MORE")
						wg.Done()
						return
					}
					fmt.Println("ERROR SAVING FILE:", err)
					wg.Done()
					return
					//cancelFunc()
				}

			}
		}()
	}
	wg.Wait()
}

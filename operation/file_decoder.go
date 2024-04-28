package operation

import (
	"fmt"
	"sync"

	"github.com/Alieksieiev0/test-task/iterator"
	"golang.org/x/net/context"
)

func NewAsyncFileDecoder[T AsyncErrorProneEntry[string], V AsyncEntry[error]]() *AsyncFileDecoder[T, V] {
	ctx, close := context.WithCancel(context.Background())
	return &AsyncFileDecoder[T, V]{ctx: ctx, close: close}
}

type AsyncFileDecoder[T AsyncErrorProneEntry[string], V AsyncEntry[error]] struct {
	ctx   context.Context
	close func()
}

func (a *AsyncFileDecoder[T, V]) Close() error {
	a.close()
	return nil
}

func (a *AsyncFileDecoder[T, V]) Run(
	readIterator iterator.Iterator[T],
	writeCallbackIterator iterator.Iterator[func(T) V],
) {
	var wg sync.WaitGroup
	for {
		wg.Add(1)
		readEntry := readIterator.Next()
		readResults := readEntry.Val()
		if readEntry.Err() != nil {
			fmt.Println(readEntry.Err())
			wg.Done()
			break
		}

		writeCallbackEntry := writeCallbackIterator.Next()
		if writeCallbackEntry.Err() != nil {
			wg.Done()
			break
		}
		writerCallback := writeCallbackEntry.Val()
		writeResults := writerCallback(readResults)

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

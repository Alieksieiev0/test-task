package operation

import (
	"context"
	"time"

	"github.com/Alieksieiev0/test-task/testtask"
)

func NewAsyncReadOperationFactory() AsyncReadOperationFactory {
	return AsyncReadOperationFactory{}
}

type AsyncReadOperationFactory struct {
}

func (a AsyncReadOperationFactory) Create(
	delayFunc func() time.Duration,
) func(context.Context) testtask.Operation[testtask.ErrorProneEntry[string], testtask.AsyncErrorProneEntry[string]] {
	return func(
		ctx context.Context,
	) testtask.Operation[testtask.ErrorProneEntry[string], testtask.AsyncErrorProneEntry[string]] {
		return &AsyncReadOperation[testtask.ErrorProneEntry[string], testtask.AsyncErrorProneEntry[string]]{
			ctx:       ctx,
			delayFunc: delayFunc,
		}
	}
}

func NewAsyncWriteOperationFactory() AsyncWriteOperationFactory {
	return AsyncWriteOperationFactory{}
}

type AsyncWriteOperationFactory struct {
}

func (a AsyncWriteOperationFactory) Create() func(context.Context) testtask.Operation[testtask.Entry[error], testtask.AsyncEntry[error]] {
	return func(ctx context.Context) testtask.Operation[testtask.Entry[error], testtask.AsyncEntry[error]] {
		return &AsyncWriteOperation[testtask.Entry[error], testtask.AsyncEntry[error]]{
			ctx: ctx,
		}
	}
}

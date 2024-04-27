package iterator

import (
	"io"

	"go.uber.org/atomic"
)

type Cb[T any] func() *Step[T]

func NewCallbackBased[T any](cb Cb[T]) Iterator[T] {
	return &CallbackBased[T]{cb: cb, closed: atomic.NewBool(false)}
}

type CallbackBased[T any] struct {
	cb     Cb[T]
	closed *atomic.Bool
}

func (c *CallbackBased[T]) Close() {
	c.closed.Store(true)
}

func (c *CallbackBased[T]) Next() Iteration[T] {
	if c.closed.Load() {
		return NewStepErr[T](io.EOF)
	}
	res := c.cb()
	if res.Err() == io.EOF {
		c.Close()
	}
	return res
}

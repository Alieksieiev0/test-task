package iterator

type IteratorWrapper[Input, Output any] interface {
	Wrap(Iterator[Input]) Iterator[Output]
}

type Iterator[T any] interface {
	Next() Iteration[T]
	Close()
}

type Iteration[T any] interface {
	Val() T
	Err() error
}

func NewStepVal[T any](val T) *Step[T] {
	return &Step[T]{
		val: val,
	}
}

func NewStepErr[T any](err error) *Step[T] {
	return &Step[T]{
		err: err,
	}
}

type Step[T any] struct {
	val T
	err error
}

func (s *Step[T]) Val() T {
	return s.val
}

func (s *Step[T]) Err() error {
	return s.err
}

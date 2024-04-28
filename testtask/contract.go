package testtask

import (
	"github.com/Alieksieiev0/test-task/iterator"
)

type Source[T any] interface {
	Data() iterator.Iterator[T]
}

type KeyValueSource[K, V any] interface {
	Keys() iterator.Iterator[K]
	Values() iterator.Iterator[V]
}

type Operation[Input any, Output any] interface {
	Run(iterator iterator.Iterator[Input], results Output)
}

type MultiOperation[T, U any] interface {
	Run(firstIter iterator.Iterator[T], secondIter iterator.Iterator[U])
}

type Processor[Input, Output any] interface {
	Process(iterator iterator.Iterator[Input], process func(Input) Output) iterator.Iterator[Output]
}

type Parser[Input, Output any] interface {
	Parse(data Input, out Output) (Output, error)
}

type Entry[T any] interface {
	Val() T
}

type ErrorProneEntry[T any] interface {
	Entry[T]
	Err() error
}

type Closable interface {
	Close()
}

type AsyncEntry[T any] interface {
	Entry[<-chan T]
	Closable
	PassVal(T)
}

type AsyncErrorProneEntry[T any] interface {
	AsyncEntry[T]
	Err() <-chan error
	PassErr(error)
}

type PlainFactory[T any] interface {
	Create() T
}

type Factory[Input, Output any] interface {
	Create(data Input) Output
}

type FactoryFunc[Input, Output any] func(i Input) Output
type ErrorFactoryFunc[Input, Output any] func(i Input) (Output, error)

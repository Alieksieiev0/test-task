package testtask

import (
	"github.com/Alieksieiev0/test-task/iterator"
)

type Source[T any] interface {
	Data() iterator.Iterator[T]
}

type Loader[T any] interface {
	Load(source T)
}

type SourceLoader[Input, Output any] interface {
	Loader[Input]
	Source[Output]
}

type Operation[Input any, Output any] interface {
	Run(iterator iterator.Iterator[Input], results Output)
}

type Processor[Input, Output any] interface {
	Process(iterator iterator.Iterator[Input], process func(Input) Output) iterator.Iterator[Output]
}

type Parser[Input, Output any] interface {
	Parse(data Input, out Output) (Output, error)
}

type Result[T any] interface {
	Val() T
	Err() error
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

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

type Operation[T any] interface {
	Run(iterator iterator.Iterator[T]) error
}

type SourceLoader[Input, Output any] interface {
	Loader[Input]
	Source[Output]
}

type Processor[Input, Output any] interface {
	Process(iterator iterator.Iterator[Input], process func(Input) Output) iterator.Iterator[Output]
}

type Parser[Input, Output any] interface {
	Parse(data Input, out Output) (Output, error)
}

type Entry interface {
	String() string
	Err() error
}

type IOEntry[T any] interface {
	Val() T
	Err() error
}

type OSEntry[T any] interface {
	IOEntry[T]
	Close()
}

type Factory[Input, Output any] interface {
	Create(data Input) Output
}

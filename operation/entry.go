package operation

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

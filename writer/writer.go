package writer

type Writer[T any] interface {
	Write(T) error
	Close()
}

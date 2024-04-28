package reader

type Reader[T any] interface {
	Read() (T, error)
	Close()
}

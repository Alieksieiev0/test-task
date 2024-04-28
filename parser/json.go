package parser

import (
	"encoding/json"
	"io"
)

func NewJsonParser[T any]() *JsonParser[T] {
	return &JsonParser[T]{}
}

type JsonParser[T any] struct {
}

func (j *JsonParser[T]) Parse(r io.Reader, v T) (T, error) {
	err := json.NewDecoder(r).Decode(v)
	return v, err
}

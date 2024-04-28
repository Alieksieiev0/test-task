package parser

import (
	"encoding/json"
	"io"
)

func NewJsonParser[T json.Unmarshaler]() *JsonParser[T] {
	return &JsonParser[T]{}
}

type JsonParser[T json.Unmarshaler] struct {
}

func (j *JsonParser[T]) Parse(r io.Reader, v T) (T, error) {
	err := json.NewDecoder(r).Decode(v)
	return v, err
}

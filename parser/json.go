package parser

import (
	"encoding/json"
	"io"
)

type JsonParser[T any] struct {
}

func (j *JsonParser[T]) Parse(r io.Reader, v T) (T, error) {
	err := json.NewDecoder(r).Decode(v)
	return v, err
}

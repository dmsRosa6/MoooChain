package core

import (
	"errors"
	"io"
)

type Encoder[T any] interface{
	Encode(w io.Writer) error
}

type BlockEncoder struct{}

//TODO implement this
func (be BlockEncoder) Encode(w io.Writer) error{
	return errors.New("Not implemented")
}

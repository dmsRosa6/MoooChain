package core

import (
	"errors"
	"io"
)

type Decoder[T any] interface{
	Decode(r io.Reader) error
}

type BlockDecoder struct{}

//TODO implement this
func (bc BlockDecoder) Decode(r io.Reader) error{
	return errors.New("Not implemented")
}
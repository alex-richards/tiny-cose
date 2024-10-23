package cose

import (
	"errors"
	"io"
)

func (s *Sign1) ReadCBOR(in io.Reader) error {
	return errors.New("TODO")
}

func (s *Sign1) WriteCBOR(out io.Writer) error {
	return errors.New("TODO")
}

func (s *Sign1) toBeSigned(payload []byte, aad []byte) ([]byte, error) {
	return nil, errors.New("TODO")
}

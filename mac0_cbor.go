package cose

import (
	"bytes"
	"io"

	cbor "github.com/alex-richards/tiny-cbor"
)

func (m *MAC0) ReadCBOR(in io.Reader) error {
	err := cbor.ReadArray(
		in,
		func(indefinite bool, length uint64) error {
			if length != 4 {
				return ErrFormat
			}
			return nil
		},
		func(i uint64, in io.Reader) error {
			switch i {
			default:
				return ErrFormat
			}
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *MAC0) ReadProtected() (*Headers, error) {
	var protected Headers
	err := protected.ReadCBOR(bytes.NewReader(m.Protected))
	if err != nil {
		return nil, err
	}
	return &protected, nil
}

func toBeMaced(protected, payload, aad []byte) ([]byte, error) {
	out := bytes.NewBuffer(make([]byte, 0, 1+3+len(protected)+3+len(payload)+3+len(aad)))

	_, err := cbor.WriteArrayHeader(out, 4)
	if err != nil {
		return nil, err
	}

	_, err = cbor.WriteString(out, "MAC0")
	if err != nil {
		return nil, err
	}

	if len(protected) == 0 {
		return nil, ErrFormat
	}

	_, err = cbor.WriteBytes(out, protected)
	if err != nil {
		return nil, err
	}

	_, err = cbor.WriteBytes(out, aad)
	if err != nil {
		return nil, err
	}

	if payload == nil {
		return nil, ErrFormat
	}

	_, err = cbor.WriteBytes(out, payload)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

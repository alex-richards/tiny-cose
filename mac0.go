package cose

import (
	"bytes"
)

type MAC0 struct {
	Protected   []byte
	Unprotected Headers
	Payload     []byte
	Tag         []byte
}

func NewMAC0(payload, aad []byte, detach bool, tagger Tagger) (*MAC0, error) {
	protected := Headers{}
	protected.Alg = tagger.Alg()

	protectedOut := [3]byte{}
	protectedBuf := bytes.NewBuffer(protectedOut[:0])
	err := protected.WriteCBOR(protectedBuf)
	if err != nil {
		return nil, err
	}

	protectedEncoded := protectedOut[:]

	tbm, err := toBeMaced(protectedEncoded, payload, aad)
	if err != nil {
		return nil, err
	}

	tag, err := tagger.CreateTag(tbm)

	if detach {
		payload = nil
	}

	return &MAC0{
		Protected: protectedEncoded,
		Payload:   payload,
		Tag:       tag,
	}, nil
}

type Tagger interface {
	Alg() Alg
	CreateTag(toBeMaced []byte) ([]byte, error)
}

package cose

import (
	"bytes"
	"io"

	"github.com/alex-richards/tiny-cbor"
)

const (
	keyFieldKty int32 = 1
	keyFieldCrv int32 = -1
	keyFieldX   int32 = -2
	keyFieldY   int32 = -3
)

func (k *Key) ReadCBOR(in io.Reader) error {
	err := cbor.ReadMap(
		in,
		func(indefinite bool, length uint64) error {
			return nil
		},
		func(in io.Reader) error {
			f, err := cbor.ReadSigned[int32](in)
			if err != nil {
				return err
			}

			switch f {
			case keyFieldKty:
				k.Type, err = readKeyTypeCBOR(in)
				if err != nil {
					return err
				}

			case keyFieldCrv:
				k.Curve, err = readCurveCBOR(in)
				if err != nil {
					return err
				}

			case keyFieldX:
				x := bytes.NewBuffer(nil)
				err = cbor.ReadBytes(
					in,
					func(indefinite bool, length uint64) error {
						x.Grow(int(length))
						return nil
					},
					x,
				)
				if err != nil {
					return err
				}
				k.X = x.Bytes()

			case keyFieldY:
				y := bytes.NewBuffer(nil)
				err = cbor.ReadBytes(
					in,
					func(indefinite bool, length uint64) error {
						y.Grow(int(length))
						return nil
					},
					y,
				)
				if err != nil {
					return err
				}
				k.Y = y.Bytes()

			default:
				err = cbor.ReadOver(in)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	return err
}

func (k *Key) WriteCBOR(out io.Writer) error {
	l := uint64(0)
	if k.Type != KeyTypeReserved {
		l++
	}
	if k.Curve != CurveReserved {
		l++
	}
	if len(k.X) != 0 {
		l++
	}
	if len(k.Y) != 0 {
		l++
	}

	_, err := cbor.WriteMapHeader(out, l)
	if err != nil {
		return err
	}

	if k.Type != KeyTypeReserved {
		_, err = cbor.WriteSigned(out, int8(keyFieldKty))
		if err != nil {
			return err
		}
		_, err = cbor.WriteSigned(out, int8(k.Type))
		if err != nil {
			return err
		}
	}
	if k.Curve != CurveReserved {
		_, err = cbor.WriteSigned(out, int8(keyFieldCrv))
		if err != nil {
			return err
		}
		_, err = cbor.WriteSigned(out, int8(k.Curve))
		if err != nil {
			return err
		}
	}
	if len(k.X) != 0 {
		_, err = cbor.WriteSigned(out, int8(keyFieldX))
		if err != nil {
			return err
		}
		_, err = cbor.WriteBytes(out, k.X)
		if err != nil {
			return err
		}
	}
	if len(k.Y) != 0 {
		_, err = cbor.WriteSigned(out, int8(keyFieldY))
		if err != nil {
			return err
		}
		_, err = cbor.WriteBytes(out, k.Y)
		if err != nil {
			return err
		}
	}
	return nil
}

func readKeyTypeCBOR(in io.Reader) (KeyType, error) {
	kty, err := cbor.ReadSigned[int8](in)
	if err != nil {
		return KeyTypeReserved, err
	}
	switch t := KeyType(kty); t {
	case KeyTypeEC2, KeyTypeOKP:
		return t, nil
	default:
		return KeyTypeReserved, nil
	}
}

func readCurveCBOR(in io.Reader) (Curve, error) {
	crv, err := cbor.ReadSigned[int8](in)
	if err != nil {
		return CurveReserved, err
	}
	switch c := Curve(crv); c {
	case CurveP256, CurveP384, CurveP521, CurveX25519, CurveX448, CurveEd25519, CurveEd448:
		return c, nil
	default:
		return CurveReserved, nil
	}
}

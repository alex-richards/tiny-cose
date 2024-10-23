package cose

import (
	"bytes"
	"errors"
	"io"

	cbor "github.com/alex-richards/tiny-cbor"
)

func (h *Headers) ReadCBOR(in io.Reader) error {
	return cbor.ReadMap(
		in,
		func(indefinite bool, length uint64) error {
			return nil
		},
		func(in io.Reader) error {
			var err error

			label, err := readHeaderLabelCBOR(in)
			if err != nil {
				return err
			}

			switch label {
			case HeaderLabelAlg:
				h.Alg, err = readAlgCBOR(in)
				if err != nil {
					return err
				}

			case HeaderLabelCrit:
				h.Crit, err = readCritCBOR(in)
				if err != nil {
					return err
				}

			case HeaderLabelX5Chain:
				h.X5Chain, err = readX5ChainCBOR(in)
				if err != nil {
					return err
				}

			default:
				if err = cbor.ReadOver(in); err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func (h *Headers) WriteCBOR(out io.Writer) error {
	length := uint64(0)
	if h.Alg != AlgReserved {
		length++
	}

	writeCrit := len(h.Crit) > 0
	if writeCrit {
		length++
	}

	writeX5Chain := len(h.X5Chain) > 0
	if writeX5Chain {
		length++
	}

	_, err := cbor.WriteMapHeader(out, length)
	if err != nil {
		return err
	}

	if h.Alg != AlgReserved {
		if err = writeHeaderLabelCBOR(HeaderLabelAlg, out); err != nil {
			return err
		}

		if err = writeAlgCBOR(h.Alg, out); err != nil {
			return err
		}
	}

	if writeCrit {
		if err = writeHeaderLabelCBOR(HeaderLabelCrit, out); err != nil {
			return err
		}

		if err = writeCritCBOR(h.Crit, out); err != nil {
			return err
		}
	}

	if writeX5Chain {
		if err = writeHeaderLabelCBOR(HeaderLabelX5Chain, out); err != nil {
			return err
		}

		if err = writeX5ChainCBOR(h.X5Chain, out); err != nil {
			return err
		}
	}

	return nil
}

func readHeaderLabelCBOR(in io.Reader) (HeaderLabel, error) {
	val, err := cbor.ReadSigned[int32](in)
	if err != nil {
		return HeaderLabelReserved, err
	}

	return HeaderLabel(val), nil
}

func writeHeaderLabelCBOR(label HeaderLabel, out io.Writer) error {
	_, err := cbor.WriteSigned(out, int32(label))
	return err
}

func readAlgCBOR(in io.Reader) (Alg, error) {
	val, err := cbor.ReadSigned[int32](in)
	if err != nil {
		return AlgReserved, err
	}

	return Alg(val), nil
}

func writeAlgCBOR(alg Alg, out io.Writer) error {
	_, err := cbor.WriteSigned(out, int32(alg))
	return err
}

func readCritCBOR(in io.Reader) ([]HeaderLabel, error) {
	var out []HeaderLabel
	err := cbor.ReadArray(
		in,
		func(indefinite bool, length uint64) error {
			out = make([]HeaderLabel, 0, length)
			return nil
		},
		func(i uint64, in io.Reader) error {
			label, err := readHeaderLabelCBOR(in)
			if err != nil {
				return err
			}
			out = append(out, label)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func writeCritCBOR(crit []HeaderLabel, out io.Writer) error {
	_, err := cbor.WriteArrayHeader(out, uint64(len(crit)))
	if err != nil {
		return err
	}
	for _, label := range crit {
		err = writeHeaderLabelCBOR(label, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func readX5ChainCBOR(in io.Reader) ([][]byte, error) {
	out := bytes.NewBuffer(nil)
	err := cbor.ReadRaw(in, out)
	if err != nil {
		return nil, err
	}
	b, err := out.ReadByte()
	if err != nil {
		return nil, err
	}
	err = out.UnreadByte()
	if err != nil {
		return nil, err
	}

	switch cbor.MajorType(b & cbor.MajorTypeMask) {
	case cbor.MajorTypeArray:
		var chain [][]byte
		err = cbor.ReadArray(
			out,
			func(indefinite bool, length uint64) error {
				chain = make([][]byte, 0, length)
				return nil
			},
			func(i uint64, in io.Reader) error {
				cert := bytes.NewBuffer(nil)
				err = cbor.ReadBytes(out,
					func(indefinite bool, length uint64) error {
						cert.Grow(int(length))
						return nil
					},
					cert,
				)
				if err != nil {
					return err
				}
				chain = append(chain, cert.Bytes())
				return nil
			},
		)
		if err != nil {
			return nil, err
		}
		return chain, nil

	case cbor.MajorTypeBstr:
		cert := bytes.NewBuffer(nil)
		err = cbor.ReadBytes(out,
			func(indefinite bool, length uint64) error {
				cert.Grow(int(length))
				return nil
			},
			cert,
		)
		if err != nil {
			return nil, err
		}
		return [][]byte{cert.Bytes()}, nil

	default:
		return nil, errors.New("TODO = can't read chain")
	}
}

func writeX5ChainCBOR(chain [][]byte, out io.Writer) error {
	chainLen := len(chain)
	if chainLen == 0 {
		return nil
	}

	if chainLen == 1 {
		_, err := cbor.WriteBytes(out, chain[0])
		if err != nil {
			return err
		}
	}

	_, err := cbor.WriteArrayHeader(out, uint64(chainLen))
	if err != nil {
		return err
	}

	for _, cert := range chain {
		_, err = cbor.WriteBytes(out, cert)
		if err != nil {
			return err
		}
	}

	return nil
}

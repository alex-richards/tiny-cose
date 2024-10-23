package cose

type KeyType int8

const (
	KeyTypeReserved KeyType = 0
	KeyTypeOKP      KeyType = 1
	KeyTypeEC2      KeyType = 2
)

type Curve int8

const (
	CurveReserved Curve = 0
	CurveP256     Curve = 1
	CurveP384     Curve = 2
	CurveP521     Curve = 3
	CurveX25519   Curve = 4
	CurveX448     Curve = 5
	CurveEd25519  Curve = 6
	CurveEd448    Curve = 7
)

type Key struct {
	Type  KeyType
	Curve Curve
	X     []byte
	Y     []byte
}

package cose

type HeaderLabel int32 // tstr is also valid, but no one(?) uses it

const (
	HeaderLabelReserved HeaderLabel = 0
	HeaderLabelAlg      HeaderLabel = 1
	HeaderLabelCrit     HeaderLabel = 2
	HeaderLabelX5Chain  HeaderLabel = 25
)

type Alg int32

const (
	AlgReserved   Alg = 0
	AlgES256      Alg = -7
	AlgES384      Alg = -35
	AlgES512      Alg = -36
	AlgEdDSA      Alg = -8
	AlgHMAC256_64 Alg = 4
	AlgHMAC256    Alg = 5
	AlgHMAC384    Alg = 6
	AlgHMAC512    Alg = 7
)

type Headers struct {
	Alg     Alg
	X5Chain [][]byte
	Crit    []HeaderLabel
}

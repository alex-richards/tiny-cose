package cose

type Sign1 struct {
	Protected   []byte
	Unprotected Headers
	Payload     []byte
	Signature   []byte
}

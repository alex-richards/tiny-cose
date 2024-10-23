package cose

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func Test_ReadKey(t *testing.T) {
	enc := "a4010220012142123422425678"
	dec, err := hex.DecodeString(enc)
	if err != nil {
		t.Fatal(err)
	}
	var key Key
	err = key.ReadCBOR(bytes.NewReader(dec))
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("key = %#v\n", key)
}

func Test_WriteKey(t *testing.T) {
	key := Key{
		Type:  KeyTypeEC2,
		Curve: CurveP256,
		X:     []byte{12, 34},
		Y:     []byte{56, 78},
	}
	out := bytes.NewBuffer(nil)
	err := key.WriteCBOR(out)
	if err != nil {
		t.Fatal(err)
	}
	println(hex.EncodeToString(out.Bytes()))
}

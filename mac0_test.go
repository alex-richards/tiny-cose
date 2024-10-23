package cose

import (
	"crypto/hmac"
	"crypto/sha256"
	"testing"
)

type TestTagger struct{}

func (tt *TestTagger) Alg() Alg {
	return AlgHMAC256
}

func (tt *TestTagger) CreateTag(toBeMaced []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, []byte("key"))
	hm.Reset()
	n, err := hm.Write(toBeMaced)
	if n != len(toBeMaced) {
		return nil, err
	}
	return hm.Sum(nil), nil
}

func Test_MAC0_Tag(t *testing.T) {
	mac, err := NewMAC0(
		[]byte("payload"),
		[]byte("aad"),
		false,
		&TestTagger{},
	)
	if err != nil {
		t.Fatal(err)
	}

	_ = mac
}

type BenchmarkTagger struct{}

func (bt *BenchmarkTagger) Alg() Alg {
	return AlgHMAC256
}

func (bt *BenchmarkTagger) CreateTag(toBeMaced []byte) ([]byte, error) {
	return toBeMaced, nil
}

func Benchmark_NewMAC0(b *testing.B) {
	payload := []byte("payload")
	aad := []byte("aad")
	tagger := &BenchmarkTagger{}

	for range b.N {
		mac, err := NewMAC0(
			payload,
			aad,
			false,
			tagger,
		)
		if err != nil {
			b.Fatal(err)
		}

		_ = mac
	}
}

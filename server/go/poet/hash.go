package poet

import (
	"hash"
	"github.com/minio/sha256-simd"
)

type HashFunc interface {
	HashVals(vs ...[]byte) []byte
}

type Hasher struct {
	BaseHash hash.Hash
}

// TODO: Will need to Benchmark this along with Unit testing. Could have
// significant speed upgrades by using the sha256 package correctly. Need to
// research that package and use more completely.
func (h *Hasher) HashVals(vs ...[]byte) (b []byte) {
	h.BaseHash.Reset()
	for _, v := range vs {
		// Can safely ignore error as not supposed to return error. TODO: Check
		_, _ = h.BaseHash.Write(v)
	}
	b = h.BaseHash.Sum([]byte{})
	return b
}

func NewSHA256() HashFunc {
	h := new(Hasher)
	h.BaseHash = sha256.New()
	return h
}

// TODO: add new sha3 function. What implementation of Sha3 used? Std lib?

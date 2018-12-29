package verifier

import (
	"errors"
	"fmt"

	"github.com/SomniaStellarum/POET/server/go/poet"
)

// veriStore stores and retrieves labels based on BinaryID's used to verify
// the Prover based on the challenge and proof response. The single store
// (verifierStore) can only be used for a single challenge/proof.
type VerifierStore struct {
	challengeProof []byte
	binIDList      []*poet.BinaryID
}

func NewVeriStoreSingle(b *poet.BinaryID, challengeProof []byte) (v *VerifierStore, err error) {
	v = new(VerifierStore)
	sib, err := poet.Siblings(b, false)
	if err != nil {
		return nil, err
	}
	v.binIDList = append(v.binIDList, b)
	v.binIDList = append(v.binIDList, sib...)
	v.challengeProof = challengeProof
	return v, nil
}

func (v *VerifierStore) StoreLabel(b *poet.BinaryID, label []byte) error {
	v.challengeProof = append(v.challengeProof, label...)
	v.binIDList = append(v.binIDList, b)
	return nil
}

func (v *VerifierStore) GetLabel(b *poet.BinaryID) (label []byte, err error) {
	for i, b_check := range v.binIDList {
		if b.Equal(b_check) {
			idx1 := i * size
			idx2 := idx1 + size
			debugLog.Println(
				"Get Node ",
				string(b.Encode()),
				"\n",
				idx1, " ", idx2, "\n",
				v.challengeProof[idx1:idx2],
			)
			return v.challengeProof[idx1:idx2], nil
		}
	}

	return nil, errors.New(fmt.Sprintf("BinID not on list: %v", string(b.Encode())))
}

func (v *VerifierStore) LabelCalculated(b *poet.BinaryID) (bool, error) {
	for _, b_check := range v.binIDList {
		if b.Equal(b_check) {
			return true, nil
		}
	}
	return false, nil
}

package verifier

import (
	"errors"
	"fmt"

	"github.com/SomniaStellarum/POET/server/go/poet"
)

type VerifierStoreMulti struct {
	vStore           *VerifierStore
	challengeLists   map[*poet.BinaryID][]*poet.BinaryID
	currentChallenge *poet.BinaryID
}

func NewVerifierStoreMulti(challenge []byte, challengeProof [][]byte) (*VerifierStoreMulti, error) {
	v := new(VerifierStoreMulti)
	challengeBins, err := poet.GammaToBinaryIDs(challenge)
	if err != nil {
		return nil, err
	}
	if len(challengeBins) != len(challengeProof) {
		return nil, errors.New(
			fmt.Sprintf(
				"Challenge and Proof length don't match: %v vs %v",
				len(challengeBins),
				len(challengeProof),
			),
		)
	}
	v.vStore, err = NewVeriStoreSingle(challengeBins[0], challengeProof[0])
	if err != nil {
		return nil, err
	}
	v.currentChallenge = challengeBins[0]
	v.challengeLists[challengeBins[0]], err = poet.Siblings(challengeBins[0], false)
	if err != nil {
		return nil, err
	}

	for i, bin := range challengeBins[1:] {
		k := 0
		v.challengeLists[bin], err = poet.Siblings(bin, false)
		if err != nil {
			return nil, err
		}
		stored, err := v.vStore.LabelCalculated(bin)
		if err != nil {
			return nil, err
		}
		if !stored {
			idx1 := k * size
			idx2 := idx1 + size
			v.vStore.StoreLabel(bin, challengeProof[i][idx1:idx2])
			k++
		}
		for _, sib := range v.challengeLists[bin] {
			stored, err := v.vStore.LabelCalculated(sib)
			if err != nil {
				return nil, err
			}
			if !stored {
				idx1 := k * size
				idx2 := idx1 + size
				v.vStore.StoreLabel(sib, challengeProof[i][idx1:idx2])
				k++
			}
		}
	}
	return v, nil
}

func (v *VerifierStoreMulti) SetCurrentChallenge(b *poet.BinaryID) error {
	// I think there is some pointer errors here. What if the binaryID pointer
	// is created from a gRPC call. It could be representing the same BinaryID
	// but has a different pointer value. TBD: fix this issue
	_, ok := v.challengeLists[b]
	if ok {
		v.currentChallenge = b
		return nil
	}
	return errors.New("BinaryID not in challenge list")
}

func (v *VerifierStoreMulti) StoreLabel(b *poet.BinaryID, label []byte) error {
	stored, err := v.vStore.LabelCalculated(b)
	if err != nil {
		return err
	}
	if !stored {
		err = v.vStore.StoreLabel(b, label)
	}
	v.challengeLists[v.currentChallenge] = append(
		v.challengeLists[v.currentChallenge],
		b,
	)
	return err
}

func (v *VerifierStoreMulti) GetLabel(b *poet.BinaryID) (label []byte, err error) {
	return v.GetLabel(b)
}

func (v *VerifierStoreMulti) LabelCalculated(*poet.BinaryID) (bool, error) {

	return false, nil
}

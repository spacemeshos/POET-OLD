package verifier

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/spacemeshos/POET/server/go/poet"
)

type VerifierStoreMulti struct {
	vStore           *VerifierStore
	challengeLists   map[*poet.BinaryID][]*poet.BinaryID
	currentChallenge *poet.BinaryID
	n                int
}

func NewVerifierStoreMulti(challenge []byte, challengeProof [][]byte, n int) (*VerifierStoreMulti, error) {
	v := new(VerifierStoreMulti)
	v.n = n
	challengeBins, err := poet.GammaToBinaryIDs(challenge, v.n)
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
	//fmt.Println(poet.StringList(challengeBins))
	v.vStore, err = NewVeriStoreSingle(challengeBins[0], challengeProof[0])
	if err != nil {
		return nil, err
	}
	v.currentChallenge = challengeBins[0]
	v.challengeLists = make(map[*poet.BinaryID][]*poet.BinaryID)
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

func (v *VerifierStoreMulti) SetDAGSize(size int) {
	v.n = size
}

func (v *VerifierStoreMulti) SetCurrentChallenge(b *poet.BinaryID) error {
	for bin, _ := range v.challengeLists {
		if b.Equal(bin) {
			v.currentChallenge = bin
			return nil
		}
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
	label, err = v.vStore.GetLabel(b)
	debugLog.Printf(
		"Node: %v, Label: %v",
		string(b.Encode()),
		hex.EncodeToString(label),
	)
	return v.vStore.GetLabel(b)
}

func (v *VerifierStoreMulti) LabelCalculated(b *poet.BinaryID) (bool, error) {
	//fmt.Println(string(b.Encode()))
	for _, bin := range v.challengeLists[v.currentChallenge] {
		if b.Equal(bin) {
			return true, nil
		}
	}
	return false, nil
}

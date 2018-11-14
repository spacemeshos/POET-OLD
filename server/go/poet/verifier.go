package poet

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
)

type Verifier struct {
	Prover          io.ReadWriter
	commitment      []byte
	commitmentProof []byte
	challenge       []byte
	challengeProof  []byte
}

func NewVerifier(Prover io.ReadWriter) (v *Verifier) {
	v = new(Verifier)
	v.Prover = Prover
	return v
}

func (v *Verifier) Commit(statement []byte) error {
	v.commitment = statement
	_, err := v.Prover.Write(statement)
	return err
}

func (v *Verifier) GetCommitProof() (b []byte, err error) {
	b = make([]byte, size)
	_, err = v.Prover.Read(b)
	v.commitmentProof = b
	return b, err
}

// Will return error if proof no good. Determine if send args or include data
// in data structure of Verifier
// Note: Do we need this function? Spec seems to be that no verification after
// commitment proof returned, only when challenge proof returned.
func (v *Verifier) VerifyCommitProof() error {
	return errors.New("Verify Not Coded")
}

func (v *Verifier) SelectChallenge() (challenge []byte, err error) {
	// TODO: Write SelectChallenge
	challengeInt := rand.Intn(n)
	binID, err := NewBinaryID(uint(challengeInt), n)
	if err != nil {
		return nil, err
	}
	challenge = binID.Encode()
	v.challenge = challenge
	debugLog.Println("Challenge: ", string(challenge))
	return challenge, nil
}

func (v *Verifier) Challenge() error {
	_, err := v.Prover.Write(v.challenge)
	return err
}

func (v *Verifier) GetChallengeProof() (b []byte, err error) {
	size := 32 * (n + 1) // TODO: Determine size. Should be size of hash times n (size of DAG)
	b = make([]byte, size)
	_, err = v.Prover.Read(b)
	if err != nil {
		return nil, err
	}
	v.challengeProof = b
	return b, err
}

// Will return error if proof no good. Determine if send args or include data
// in data structure of Verifier
func (v *Verifier) VerifyChallengeProof() (err error) {
	// For a single leaf challenge: Calc Hash of Leaf then walk up the tree using
	// the sibling leaf hash's as going. Much of the Code (eg GetParents) should
	// be developped through Prover code. TODO: Complete this verify function
	cOpts := new(ComputeOpts)
	cOpts.hash = NewSHA256()
	cOpts.commitment = v.commitment
	cOpts.commitmentHash = cOpts.hash.HashVals(v.commitment)
	// If challenge is nil, this is a NIP proof. Must generate NIP challenge
	if v.challenge == nil {
		gammas := CalcNIPChallenge(v.commitmentProof, cOpts)
		v.challenge = gammas[0].Encode()
	}
	challengeID := NewBinaryIDBytes(v.challenge)
	debugLog.Println("Challenge: ", string(challengeID.Encode()))
	cOpts.store, err = NewVeriStoreSingle(challengeID, v.challengeProof)
	if err != nil {
		return err
	}
	siblings, err := Siblings(challengeID)
	if err != nil {
		return err
	}
	debugLog.Println("Challenge: ", string(challengeID.Encode()))
	root, err := NewBinaryID(0, 0)
	if err != nil {
		return err
	}
	siblings = siblings[1:]
	for _, sib := range siblings {
		debugLog.Println("Sib: ", string(sib.Encode()))
		sib.FlipBit(sib.Length)
		_ = ComputeLabel(sib, cOpts) // ComputeLabel stores the label so can ignore
	}
	rootCalc := ComputeLabel(root, cOpts)
	if err != nil {
		return err
	}
	if bytes.Equal(v.commitmentProof, rootCalc) {
		return nil
	}
	return errors.New("Verify Failed")
}

type verifierStore struct {
	challengeProof []byte
	binIDList      []*BinaryID
}

func NewVeriStoreSingle(b *BinaryID, challengeProof []byte) (v *verifierStore, err error) {
	v = new(verifierStore)
	sib, err := Siblings(b)
	if err != nil {
		return nil, err
	}
	v.binIDList = append(v.binIDList, b)
	v.binIDList = append(v.binIDList, sib...)
	v.challengeProof = challengeProof
	return v, nil
}

func (v *verifierStore) StoreLabel(b *BinaryID, label []byte) error {
	v.challengeProof = append(v.challengeProof, label...)
	v.binIDList = append(v.binIDList, b)
	debugLog.Println(v.challengeProof)
	return nil
}

func (v *verifierStore) GetLabel(b *BinaryID) (label []byte, err error) {
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

func (v *verifierStore) LabelCalculated(*BinaryID) (bool, error) {
	return true, nil
}

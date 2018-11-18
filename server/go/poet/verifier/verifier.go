package verifier

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"

	"github.com/SomniaStellarum/POET/server/go/poet"
)

type Verifier struct {
	Prover          io.ReadWriter
	n               int
	commitment      []byte
	commitmentProof []byte
	challenge       []byte
	challengeProof  []byte
}

func NewVerifier(Prover io.ReadWriter, n int) (v *Verifier) {
	v = new(Verifier)
	v.Prover = Prover
	v.n = n
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
	challengeInt := rand.Intn(v.n)
	binID, err := poet.NewBinaryID(uint(challengeInt), v.n)
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
	s := size * (v.n + 1) // TODO: Determine size. Should be size of hash times n (size of DAG)
	b = make([]byte, s)
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
	cOpts := new(poet.ComputeOpts)
	cOpts.Hash = poet.NewSHA256()
	cOpts.Commitment = v.commitment
	cOpts.CommitmentHash = cOpts.Hash.HashVals(v.commitment)
	// If challenge is nil, this is a NIP proof. Must generate NIP challenge
	if v.challenge == nil {
		gammas := poet.CalcNIPChallenge(v.commitmentProof, cOpts)
		v.challenge = gammas[0].Encode()
	}
	challengeID := poet.NewBinaryIDBytes(v.challenge)
	debugLog.Println("Challenge: ", string(challengeID.Encode()))
	cOpts.Store, err = NewVeriStoreSingle(challengeID, v.challengeProof)
	if err != nil {
		return err
	}
	siblings, err := poet.Siblings(challengeID)
	if err != nil {
		return err
	}
	debugLog.Println("Challenge: ", string(challengeID.Encode()))
	root, err := poet.NewBinaryID(0, 0)
	if err != nil {
		return err
	}
	siblings = siblings[1:]
	for _, sib := range siblings {
		debugLog.Println("Sib: ", string(sib.Encode()))
		sib.FlipBit(sib.Length)
		_ = poet.ComputeLabel(sib, cOpts) // ComputeLabel stores the label so can ignore
	}
	rootCalc := poet.ComputeLabel(root, cOpts)
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
	binIDList      []*poet.BinaryID
}

func NewVeriStoreSingle(b *poet.BinaryID, challengeProof []byte) (v *verifierStore, err error) {
	v = new(verifierStore)
	sib, err := poet.Siblings(b)
	if err != nil {
		return nil, err
	}
	v.binIDList = append(v.binIDList, b)
	v.binIDList = append(v.binIDList, sib...)
	v.challengeProof = challengeProof
	return v, nil
}

func (v *verifierStore) StoreLabel(b *poet.BinaryID, label []byte) error {
	v.challengeProof = append(v.challengeProof, label...)
	v.binIDList = append(v.binIDList, b)
	debugLog.Println(v.challengeProof)
	return nil
}

func (v *verifierStore) GetLabel(b *poet.BinaryID) (label []byte, err error) {
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

func (v *verifierStore) LabelCalculated(*poet.BinaryID) (bool, error) {
	return true, nil
}

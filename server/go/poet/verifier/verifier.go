package verifier

import (
	"bytes"
	"errors"
	"math/rand"

	"github.com/SomniaStellarum/POET/server/go/poet"
)

type Verifier struct {
	n               int
	hash            poet.HashFunc
	commitment      []byte
	commitmentProof []byte
	challenge       []byte
	challengeProof  [][]byte
}

func NewVerifier(n int) (v *Verifier) {
	v = new(Verifier)
	v.n = n
	return v
}

func (v *Verifier) SetDAGSize(n int) {
	v.n = n
}

func (v *Verifier) SetHash(hash poet.HashFunc) {
	v.hash = hash
}

func (v *Verifier) SetCommitment(b []byte) {
	v.commitment = b
}

func (v *Verifier) SetCommitmentProof(phi []byte) {
	v.commitmentProof = phi
}

func (v *Verifier) SetChallenge(challenge []byte) {
	v.challenge = challenge
}

func (v *Verifier) SetChallengeProof(challengeProof [][]byte) {
	v.challengeProof = challengeProof
}

func (v *Verifier) SelectRndChallenge() (challenge []byte, err error) {
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

// Will return error if proof no good. Determine if send args or include data
// in data structure of Verifier
func (v *Verifier) VerifyChallengeProof() (err error) {
	// For a single leaf challenge: Calc Hash of Leaf then walk up the tree using
	// the sibling leaf hash's as going. Much of the Code (eg GetParents) should
	// be developped through Prover code. TODO: Complete this verify function
	cOpts := new(poet.ComputeOpts)
	cOpts.Hash = v.hash
	cOpts.Commitment = v.commitment
	cOpts.CommitmentHash = cOpts.Hash.HashVals(v.commitment)
	// If challenge is nil, this is a NIP proof. Must generate NIP challenge
	if v.challenge == nil {
		gammas := poet.CalcNIPChallenge(v.commitmentProof, cOpts)
		for _, b := range gammas {
			v.challenge = append(v.challenge, b.Encode()...)
		}
	}
	challengeIDs, err := poet.GammaToBinaryIDs(v.challenge)
	if err != nil {
		return err
	}
	vStore, err := NewVerifierStoreMulti(v.challenge, v.challengeProof)
	if err != nil {
		return err
	}
	root, err := poet.NewBinaryID(0, 0)
	if err != nil {
		return err
	}
	for _, bin := range challengeIDs {
		vStore.SetCurrentChallenge(bin)
		cOpts.Store = vStore
		rootCalc := poet.ComputeLabel(root, cOpts)
		poet.PrintDAG(root, cOpts.Store, "Verifier")
		if !bytes.Equal(v.commitmentProof, rootCalc) {
			return errors.New("Verify Failed")
		}
	}
	return nil
}

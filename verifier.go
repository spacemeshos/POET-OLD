package poet

import (
	"errors"
	"io"
)

type Verifier struct {
	Prover io.ReadWriter
}

func NewVerifier(Prover io.ReadWriter) (v *Verifier) {
	v = new(Verifier)
	v.Prover = Prover
	return v
}

func (v *Verifier) Commit(statement []byte) error {
	_, err := v.Prover.Write(statement)
	return err
}

func (v *Verifier) GetCommitProof() (b []byte, err error) {
	size := 32 // TODO: Set const at init. If algo change, would need to update
	b = make([]byte, size)
	_, err = v.Prover.Read(b)
	return b, err
}

// Will return error if proof no good. Determine if send args or include data
// in data structure of Verifier
func (v *Verifier) VerifyCommitProof() error {
	return errors.New("Verify Not Coded")
}

func (v *Verifier) SelectChallenge() (challenge []byte) {
	// TODO: Write SelectChallenge
	return challenge
}

func (v *Verifier) Challenge(challenge []byte) error {
	_, err := v.Prover.Write(challenge)
	return err
}

func (v *Verifier) GetChallengeProof() (b []byte, err error) {
	size := 32 // TODO: Determine size. Should be size of hash times n (size of DAG)
	b = make([]byte, size)
	_, err = v.Prover.Read(b)
	return b, err
}

// Will return error if proof no good. Determine if send args or include data
// in data structure of Verifier
func (v *Verifier) VerifyChallengeProof() error {
	return errors.New("Verify Not Coded")
}

package poet

import "errors"

// This type will provide the inteface to the Prover. It implements the
// io.ReadWriter interface, which will allow it to sit behind an RPC Server
// or be linked directly to a verifier.
// CurrentState is used to
type Prover struct {
	CreateChallenge bool
	CurrentState    State
	// other types based on implementation. Eg leveldb client & DAG
}

// Satifying io.ReadWriter interface. In Start State it returns Proof from
// Commitment. In WaitingChalleng State, it returns Challenge Proof. Both
// commitment and challenge are encoded as a byte slice (b). To retrieve
// the proof, the verifier calls Read.
func (p *Prover) Write(b []byte) (n int, err error) {
	if p.CurrentState == Start {
		err = p.CalcCommitProof(b)
		p.CurrentState = Commited
	} else if p.CurrentState == WaitingChallenge {
		err = p.CalcChallengeProof(b)
		p.CurrentState = ProofDone
	} else {
		return 0, errors.New("Prover in Wrong State for Write")
	}
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (p *Prover) Read(b []byte) (n int, err error) {
	// TODO: Check size of b. Read only supposed to send len(b) bytes. If
	// not big enough, need to return error
	if p.CurrentState == Commited {
		b, err = p.SendCommitProof()
		p.CurrentState = WaitingChallenge
	} else if p.CurrentState == ProofDone {
		b, err = p.SendChallengeProof()
		p.CurrentState = Start // For now, this just loops back.
		// TODO: What is the logic to change state? May have other functions to
		// reset prover to Start state (clear DAG and ready for new statement)
	} else {
		return 0, errors.New("Prover in Wrong State for Read")
	}
	return 0, nil
}

func (p *Prover) CalcCommitProof([]byte) error {
	return nil
}

func (p *Prover) CalcChallengeProof([]byte) error {
	return nil
}

func (p *Prover) SendCommitProof() (b []byte, err error) {
	return b, nil
}

func (p *Prover) SendChallengeProof() (b []byte, err error) {
	return b, nil
}

type State int

const (
	Start State = iota
	Commited
	WaitingChallenge
	ProofDone
)

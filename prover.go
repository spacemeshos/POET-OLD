package poet

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
// commitment and challenge are encoded as a byte slice (b)
func (p *Prover) Write(b []byte) (n int, err error) {
	if p.CurrentState == Start {
		err = p.CalcCommitProof(b)
	} else {
		err = p.CalcChallengeProof(b)
	}
	if err != nil {
		return 0, nil
	}
	return len(b), nil
}

func (p *Prover) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (p *Prover) CalcCommitProof([]byte) error {
	return nil
}

func (p *Prover) CalcChallengeProof([]byte) error {
	return nil
}

type State int

const (
	Start State = iota
	Commited
	WaitingChallenge
	ProofDone
)

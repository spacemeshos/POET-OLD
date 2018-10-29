package poet

type Commiter interface {
	Commit([]byte) error
	Proof() ([]byte, error)
	Verified() (bool, error)
}

type RandomVeriLocal struct {
	Ver      *Verifier
	proof    []byte
	verified bool
	err      error
}

func NewRandomVeriLocal() *RandomVeriLocal {
	rvl := new(RandomVeriLocal)
	p := NewProver(false)
	rvl.Ver = NewVerifier(p)
	return rvl
}

func (rvl *RandomVeriLocal) Commit(b []byte) error {
	err := rvl.Commit(b)
	go rvl.run()
	return err
}

// Goroutine to coordinate challenge and verification of the proof
func (rvl *RandomVeriLocal) run() {
	proof, err := rvl.Ver.GetCommitProof()
	if err != nil {
		rvl.err = err
	}
	rvl.proof = proof
	// TODO: complete
}

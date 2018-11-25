package poet

import (
	"encoding/hex"
	"errors"
	"log"
)

func CalcNIPChallenge(rootHash []byte, cOpts *ComputeOpts) (b_list []*BinaryID) {
	i := 0
	// Only generate one gamma for now, but later will modify to take up to t
	// gamma's as challenge.
	//for i := 0; i < t; i++ {
	b := NewBinaryIDInt(uint(i))
	v := cOpts.Hash.HashVals(cOpts.CommitmentHash, rootHash, b.Encode())
	//v := cOpts.Hash.HashVals(cOpts.Commitment, rootHash, b.Encode())
	v = v[:n]
	gamma := NewBinaryIDBytes(v)
	b_list = append(b_list, gamma)
	//}
	return b_list
}

// // This type will provide the inteface to the Prover. It implements the
// // io.ReadWriter interface, which will allow it to sit behind an RPC Server
// // or be linked directly to a verifier.
// // CurrentState is used to
type Prover struct {
	CreateNIPChallenge bool
	CurrentState       State
	rootHash           []byte
	challengeProof     []byte
	commitment         []byte
	commitmentHash     []byte
	store              StorageIO
	hash               HashFunc
}

func NewProver(CreateChallenge bool) *Prover {
	p := new(Prover)
	p.CreateNIPChallenge = CreateChallenge
	p.store = NewFileIO()
	p.hash = NewSHA256()
	return p
}

// // Satifying io.ReadWriter interface. In Start State it returns Proof from
// // Commitment. In WaitingChalleng State, it returns Challenge Proof. Both
// // commitment and challenge are encoded as a byte slice (b). To retrieve
// // the proof, the verifier calls Read.
func (p *Prover) Write(b []byte) (n int, err error) {
	if p.CurrentState == Start {
		err = p.CalcCommitProof(b)
		if err != nil {
			return 0, err
		}
		p.CurrentState = Commited
	} else if p.CurrentState == WaitingChallenge {
		err = p.CalcChallengeProof(b)
		if err != nil {
			return 0, err
		}
		p.CurrentState = ProofDone
	} else {
		return 0, errors.New("Prover in Wrong State for Write")
	}
	return len(b), nil
}

func (p *Prover) Read(b []byte) (n int, err error) {
	// TODO: Check size of b. Read only supposed to send len(b) bytes. If
	// not big enough, need to return error
	if p.CurrentState == Commited {
		proof, err := p.SendCommitProof()
		if err != nil {
			return 0, err
		}
		debugLog.Println(hex.EncodeToString(proof))
		copy(b, proof)
		if p.CreateNIPChallenge {
			err = p.CalcNIPCommitProof()
			if err != nil {
				return 0, err
			}
			p.CurrentState = ProofDone
		} else {
			p.CurrentState = WaitingChallenge
		}
	} else if p.CurrentState == ProofDone {
		proof, err := p.SendChallengeProof()
		if err != nil {
			return 0, err
		}
		copy(b, proof)
		p.CurrentState = Start // For now, this just loops back.
		// TODO: What is the logic to change state? May have other functions to
		// reset prover to Start state (clear DAG and ready for new statement)
	} else {
		return 0, errors.New("Prover in Wrong State for Read")
	}
	return 0, nil
}

// CalcCommitProof calculates the proof of seqeuntial work
func (p *Prover) CalcCommitProof(commitment []byte) error {
	cOpts := new(ComputeOpts)
	cOpts.Hash = p.hash
	cOpts.Store = p.store
	cOpts.Commitment = commitment
	cOpts.CommitmentHash = cOpts.Hash.HashVals(commitment)
	debugLog.Println("CommitmentHash: ", hex.EncodeToString(cOpts.CommitmentHash))
	p.commitment = commitment
	p.commitmentHash = cOpts.CommitmentHash
	node, err := NewBinaryID(0, 0)
	if err != nil {
		log.Panic("Error creating BinaryID: ", err)
	}
	phi := ComputeLabel(node, cOpts)
	p.rootHash = phi
	return nil
}

// SendCommitProof send the phi (root Hash) to the verifier
func (p *Prover) SendCommitProof() (b []byte, err error) {
	return p.rootHash, nil
}

// CalcNIPCommitProof proof created by computing openH for the challenge
// TODO: modify so that each Hx(phi, i) is used to calc challenge (first n bits)
func (p *Prover) CalcNIPCommitProof() error {
	cOpts := new(ComputeOpts)
	cOpts.Hash = p.hash
	cOpts.Commitment = p.commitment
	cOpts.CommitmentHash = p.commitmentHash
	cOpts.Store = p.store
	gamma := CalcNIPChallenge(p.rootHash, cOpts)
	// TODO: When return multiple gamma's, need to modify this code to handle that
	p.CalcChallengeProof(gamma[0].Encode())
	return nil
}

// CalcChallengeProof
func (p *Prover) CalcChallengeProof(gamma []byte) error {

	var proof []byte

	gammaBinID := NewBinaryIDBytes(gamma)
	siblings, err := Siblings(gammaBinID, false)
	if err != nil {
		return nil
	}
	// Should check if label was calculated?
	label_gamma, err := p.store.GetLabel(gammaBinID)
	if err != nil {
		return err
	}
	proof = append(proof, label_gamma...)
	debugLog.Println("GammaID: ", string(gammaBinID.Encode()))
	for _, sib := range siblings {
		// Should check if label was calculated?
		label, err := p.store.GetLabel(sib)
		if err != nil {
			return err
		}
		debugLog.Println("Appending label for ", string(sib.Encode()))
		proof = append(proof, label...)
	}

	p.challengeProof = proof
	return nil
}

// SendChallengeProof
func (p *Prover) SendChallengeProof() (b []byte, err error) {
	return p.challengeProof, nil
}

type State int

const (
	Start State = iota
	Commited
	WaitingChallenge
	ProofDone
)

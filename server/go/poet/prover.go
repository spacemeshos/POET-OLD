package poet

import (
	"encoding/hex"
	"log"
)

func CalcNIPChallenge(rootHash []byte, cOpts *ComputeOpts) (b_list []*BinaryID) {
	for i := 0; i < int(cOpts.T); i++ {
		b := NewBinaryIDInt(uint(i))
		v := cOpts.Hash.HashVals(cOpts.Commitment, rootHash, b.Val)
		v = v[:cOpts.N]
		gamma := NewBinaryIDBytes(v)
		b_list = append(b_list, gamma)
	}
	return b_list
}

// This type will provide the inteface to the Prover. It implements the
// io.ReadWriter interface, which will allow it to sit behind an RPC Server
// or be linked directly to a verifier.
// CurrentState is used to
type Prover struct {
	rootHash       []byte
	challengeProof [][]byte
	commitment     []byte
	commitmentHash []byte
	store          StorageIO
	hash           HashFunc
	started        bool
	t              uint
	n              int
}

func NewProver() *Prover {
	p := new(Prover)
	p.store = NewFileIO()
	p.hash = NewSHA256()
	p.t = uint(t)
	p.n = n
	return p
}

// CalcCommitProof calculates the proof of seqeuntial work
func (p *Prover) CalcCommitProof(commitment []byte) error {
	p.started = true
	cOpts := new(ComputeOpts)
	cOpts.Hash = p.hash
	cOpts.Store = p.store
	cOpts.Commitment = commitment
	cOpts.CommitmentHash = cOpts.Hash.HashVals(commitment)
	cOpts.T = p.t
	cOpts.N = p.n
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

// CommitProof returns phi (root Hash)
func (p *Prover) CommitProof() (b []byte, err error) {
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
	cOpts.T = p.t
	cOpts.N = p.n
	gamma := CalcNIPChallenge(p.rootHash, cOpts)
	// TODO: When return multiple gamma's, need to modify this code to handle that
	var gammaBytes []byte
	for _, g := range gamma {
		gammaBytes = append(gammaBytes, g.Encode()...)
	}
	p.CalcChallengeProof(gammaBytes)
	return nil
}

// CalcChallengeProof
func (p *Prover) CalcChallengeProof(gamma []byte) error {
	var proof [][]byte
	var BinIDs []*BinaryID
	gammaBinIDs, err := GammaToBinaryIDs(gamma, p.n)
	if err != nil {
		return err
	}
	for _, g := range gammaBinIDs {
		var labels []byte
		var added bool
		BinIDs, added = CheckAndAdd(BinIDs, g)
		if added {
			// Should check if label was calculated?
			label_gamma, err := p.store.GetLabel(g)
			if err != nil {
				return err
			}
			labels = append(labels, label_gamma...)
		}
		// debugLog.Println("GammaID: ", string(g.Encode()))
		siblings, err := Siblings(g, false)
		if err != nil {
			return nil
		}
		for _, sib := range siblings {
			BinIDs, added = CheckAndAdd(BinIDs, sib)
			if added {
				// Should check if label was calculated?
				label, err := p.store.GetLabel(sib)
				if err != nil {
					return err
				}
				// debugLog.Println("Appending label for ", string(sib.Encode()))
				labels = append(labels, label...)
			}
		}
		proof = append(proof, labels)
	}

	p.challengeProof = proof
	return nil
}

// SendChallengeProof
func (p *Prover) ChallengeProof() (b [][]byte, err error) {
	return p.challengeProof, nil
}

// Can only call this before the DAG has started to be calculated. Need to add
// a flag and check for this.
// TODO: might be cleaner to have n as a field in the Prover struct. Also would
// need to have it as a field in BinaryID as well in that case.
func (p *Prover) ChangeDAGSize(size int) {
	if !p.started {
		p.n = size
		n = size
		f, ok := p.store.(SetDAGSizer)
		if ok {
			f.SetDAGSize(size)
		}
	}
}

func (p *Prover) ChangeHashFunc(hfunc HashFunc) {
	if !p.started {
		p.hash = hfunc
	}
}

func (p *Prover) ChangeT(t uint) {
	if !p.started {
		p.t = t
	}
}

func (p *Prover) Clean() {
	p.store = NewFileIO()
	p.challengeProof = nil
	p.commitment = nil
	p.commitmentHash = nil
	p.rootHash = nil
	p.started = false
}

func (p *Prover) ShowDAG() {
	root, err := NewBinaryID(0, 0)
	if err != nil {
		return
	}
	PrintDAG(root, p.n, p.store, "Prover")
}

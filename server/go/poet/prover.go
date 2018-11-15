package poet

import (
	"encoding/hex"
	"errors"
	"log"
)

// Siblings returns the list of siblings along the path to the root
//
// Takes in an instance of class BinaryString and returns a list of the
// siblings of the nodes of the path to to root of a binary tree. Also
// returns the node itself, so there are N+1 items in the list for a
// tree with length N.
//
func Siblings(node *BinaryID) ([]*BinaryID, error) {

	var siblings []*BinaryID
	// Do we really need the node on the siblings list?
	//siblings = append(siblings, node)
	newBinaryID := NewBinaryIDCopy(node)
	for i := 0; i < node.Length; i++ {
		if i == node.Length-1 {
			newBinaryID.FlipBit(newBinaryID.Length)
			siblings = append(siblings, newBinaryID)
		} else {
			id := NewBinaryIDCopy(newBinaryID)
			id.FlipBit(id.Length)
			siblings = append(siblings, id)
			newBinaryID.TruncateLastBit()
		}
	}

	return siblings, nil
}

func LeftSiblings(node *BinaryID) ([]*BinaryID, error) {

	var siblings []*BinaryID
	// Do we really need the node on the siblings list?
	//siblings = append(siblings, node)
	//debugLog.Printf("Getting Left Siblings for node: %v\n", string(node.Encode()))
	newBinaryID := NewBinaryIDCopy(node)
	for i := 0; i < node.Length; i++ {
		if i == node.Length-1 {
			newBinaryID.FlipBit(newBinaryID.Length)
			// TODO: Add error check
			bit, _ := newBinaryID.GetBit(newBinaryID.Length)
			if bit == 0 {
				//debugLog.Printf("Adding %v to Left Siblings", string(newBinaryID.Encode()))
				siblings = append(siblings, newBinaryID)
			}
		} else {
			id := NewBinaryIDCopy(newBinaryID)
			id.FlipBit(id.Length)
			// TODO: Add error check
			bit, _ := id.GetBit(id.Length)
			if bit == 0 {
				//debugLog.Printf("Adding %v to Left Siblings", string(id.Encode()))
				siblings = append(siblings, id)
			}
			newBinaryID.TruncateLastBit()
		}
	}

	return siblings, nil
}

// GetParents get parents of a node
func GetParents(node *BinaryID) ([]*BinaryID, error) {
	var parents []*BinaryID
	parents = make([]*BinaryID, 0, n-1)

	if node.Length == n {
		left, err := LeftSiblings(node)
		if err != nil {
			return nil, err
		}
		parents = append(parents, left...)
	} else {
		id0 := NewBinaryIDCopy(node)
		id0.AddBit(0)
		parents = append(parents, id0)

		id1 := NewBinaryIDCopy(node)
		id1.AddBit(1)
		parents = append(parents, id1)
	}
	return parents, nil
}

type ComputeOpts struct {
	commitment     []byte
	commitmentHash []byte
	hash           HashFunc
	store          StorageIO
}

// ComputeLabel of a node id
func ComputeLabel(node *BinaryID, cOpts *ComputeOpts) []byte {
	parents, _ := GetParents(node)
	var parentLabels []byte
	// Loop through the parents and try to calculate their labels
	// if doesn't exist in computed
	for _, parent := range parents {
		// check if the label exists
		exists, err := cOpts.store.LabelCalculated(parent)
		if err != nil {
			log.Panic("Error Checking Label: ", err)
		}
		if exists {
			pLabel, err := cOpts.store.GetLabel(parent)
			if err != nil {
				log.Panic("Error Getting Label: ", err)
			}
			parentLabels = append(parentLabels, pLabel...)
		} else {
			// compute the label
			label := ComputeLabel(parent, cOpts)
			parentLabels = append(parentLabels, label...)
		}
	}

	debugLog.Printf(
		"Inputs: %v %v %v\n",
		hex.EncodeToString(cOpts.commitmentHash),
		hex.EncodeToString(node.Encode()),
		hex.EncodeToString(parentLabels),
	)

	result := cOpts.hash.HashVals(
		cOpts.commitment,
		node.Encode(),
		parentLabels)

	debugLog.Println(
		"Hash for node ",
		string(node.Encode()),
		" calculated: ",
		hex.EncodeToString(result),
	)

	err := cOpts.store.StoreLabel(node, result)
	if err != nil {
		log.Panic("Error Storing Label: ", err)
	}
	return result
}

func CalcNIPChallenge(rootHash []byte, cOpts *ComputeOpts) (b_list []*BinaryID) {
	i := 0
	// Only generate one gamma for now, but later will modify to take up to t
	// gamma's as challenge.
	//for i := 0; i < t; i++ {
	b := NewBinaryIDInt(uint(i))
	v := cOpts.hash.HashVals(cOpts.commitmentHash, rootHash, b.Encode())
	//v := cOpts.hash.HashVals(cOpts.commitment, rootHash, b.Encode())
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
	cOpts.hash = p.hash
	cOpts.store = p.store
	cOpts.commitment = commitment
	cOpts.commitmentHash = cOpts.hash.HashVals(commitment)
	debugLog.Println("CommitmentHash: ", hex.EncodeToString(cOpts.commitmentHash))
	p.commitment = commitment
	p.commitmentHash = cOpts.commitmentHash
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
	cOpts.hash = p.hash
	cOpts.commitment = p.commitment
	cOpts.commitmentHash = p.commitmentHash
	cOpts.store = p.store
	gamma := CalcNIPChallenge(p.rootHash, cOpts)
	// TODO: When return multiple gamma's, need to modify this code to handle that
	p.CalcChallengeProof(gamma[0].Encode())
	return nil
}

// CalcChallengeProof
func (p *Prover) CalcChallengeProof(gamma []byte) error {

	var proof []byte

	gammaBinID := NewBinaryIDBytes(gamma)
	siblings, err := Siblings(gammaBinID)
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

package poet

import (
	"errors"
	"math"
	"sort"
	"os"
	"fmt"
)

// // This type will provide the inteface to the Prover. It implements the
// // io.ReadWriter interface, which will allow it to sit behind an RPC Server
// // or be linked directly to a verifier.
// // CurrentState is used to
type Prover struct {
	CreateNIPChallenge bool
	CurrentState    State
	// other types based on implementation. Eg leveldb client & DAG
}


func NewProver(CreateChallenge bool) *Prover {
	p := new(Prover)
	p.CreateNIPChallenge = CreateChallenge
	return p
}

// // Satifying io.ReadWriter interface. In Start State it returns Proof from
// // Commitment. In WaitingChalleng State, it returns Challenge Proof. Both
// // commitment and challenge are encoded as a byte slice (b). To retrieve
// // the proof, the verifier calls Read.
func (p *Prover) Write(b []byte) (n int, err error) {
	if p.CurrentState == Start {

		err = p.CalcCommitProof(struct{b []byte})
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

// WriteToFile write the labels at height m to file
func (p *Prover) WriteToFile(data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	result, _ := file.Write(data)
	return nil
}

// GetParents get parents of a node
func (p *Prover) GetParents(node *BinaryID) ([]*BinaryID, error) {

	parents := make([]*BinaryID, 0)
	bitlist := node.BitList()
	length := len(bitlist)

	if length == n {
		for i := 1; i < (length+1); i++ {
			if bitlist[i-1] == 1 {
				data := append(bitlist[:i - 1], 0)
				id, _  := NewBinaryID(uint(i), BitsToInt(data))
				parents = append(parents, id)
			}
		}
	} else {
		data0 := append(bitlist, 0)
		id0, _ := NewBinaryID(uint(length+1), BitsToInt(data0))
		parents = append(parents, id0)

		data1 := append(bitlist,1)
		id1, _ := NewBinaryID(uint(length+1), BitsToInt(data1))
		parents = append(parents, id1)
	}

	// sort the parent ids
	sort.Slice(parents, func(a, b int) bool {
		return parents[a].GreaterThan(parents[b])
	})
	
	// get the byte values of the parents
	return parents, nil
}

// ComputeLabel of a node id
func (p *Prover) ComputeLabel(commitment []byte, node *BinaryID, hash HashFunc) []byte {
	parents, _ := p.GetParents(node)

	// should contain the concatenated byte array 
	// of parent labels
	var parentLabels []byte
	// maps the string encoding of a node id 
	// to its label bytes
	var computed map[string][]byte
	
	// Loop through the parents and try to calculate their labels
	// if doesn't exist in computed
	for i := 0; i < len(parents); i++ {
		// convert the byte array to a string representation
		str := fmt.Sprintf("%s", parents[i].Encode())
		// check if the label exists in computed
		if val, ok := computed[str]; ok {
			parentLabels = append(parentLabels, computed[str]...)
		} else {
			// compute the label
			label := p.ComputeLabel(commitment, node, hash)
			// store it in computed
			computed[str] = label
			parentLabels = append(parentLabels, label...)
		}
	}

	result := hash.HashVals(commitment, node.Val, parentLabels)
	return result
}

// ConstructDag create dag from the 
// n is the time parameter
// m is the number of layers to store
func (p *Prover) ConstructDag(commitment []byte, hash HashFunc) ([]byte, error) {
	// was told no need to use a graph anymore
	// can just compute the edges using an algorithm
	// god help us
	var l0, l1, lroot []byte

	// for height from 0 to m
	for height := 0; height < (m+1); height++ {
		// compute number of nodes for each sub tree
		numberOfNodes := int(math.Pow(float64(2), float64(height)))

		/**
		* Improvement: Can use a single loop and write offsets file
		* File offsets seems not quite easy to do cos of unknown
		* buffer length
		*/
		
		// left sub tree
		// perform left sub tree calculation
		for level := 0; level < numberOfNodes; level++ {
			leftId, _ := NewBinaryID(uint(height), level)
			leftId.AddBit(0)
			leftLabel := p.ComputeLabel(commitment, leftId, hash)
			if height ==  1 {
				l0 = leftLabel
			}
			p.WriteToFile(leftLabel)
		}

		// right sub tree
		// pefrom right sub tree calculation
		for level := 0; level < numberOfNodes; level++ {
			rightId, _ := NewBinaryID(uint(height), level)
			rightId.AddBit(1)
			rightLabel := p.ComputeLabel(commitment, rightId, hash)
			if height ==  1 {
				l1 = rightLabel
			}
			p.WriteToFile(rightLabel)
		}
	}

	rootHash := hash.HashVals(commitment, l0, l1)
	return rootHash, nil
}

type CommitProofParam struct {
	commitment []byte
	hash *HashFunc
}

// CalcCommitProof calculates the proof of seqeuntial work
func (p *Prover) CalcCommitProof(param CommitProofParam) error {
	var hashFunction *HashFunc

	hashFunction = param.hash
	if hashFunction == nil {
		hashFunction = NewSHA256()
	}
	
	graph, _ := p.ConstructDag(param.commitment, hashFunction)
	return nil
}

func (p *Prover) CalcNIPCommitProof([]byte) (error) {
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

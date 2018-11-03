package poet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
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
	siblings = append(siblings, node)

	uintValue, _ := binary.Uvarint(node.Val)
	newBinaryID, _ := NewBinaryID(uint(uintValue), node.Length)

	for i := 0; i < node.Length; i++ {
		newBinaryID.FlipBit(0)
		siblings = append(siblings, newBinaryID)
		newBinaryID.TruncateLastBit()
	}

	return siblings, nil
}

// GetParents get parents of a node
func GetParents(node *BinaryID) ([]*BinaryID, error) {
	var parents []*BinaryID
	parents = make([]*BinaryID, 0, 3)

	if node.Length == n {
		for i := 0; i < node.Length; i++ {
			j, err := node.GetBit(i)
			if err != nil {
				return nil, err
			}
			fmt.Println(node.Val, j, i)
			if j == 1 {
				id := NewBinaryIDCopy(node)
				for k := 0; k > i; k++ {
					id.TruncateLastBit()
				}
				id.FlipBit(id.Length)
				//fmt.Println("Appending\n", StringList(parents), "\n", id)
				parents = append(parents, id)
			}
		}
	} else {
		id0 := NewBinaryIDCopy(node)
		id0.AddBit(0)
		parents = append(parents, id0)

		id1 := NewBinaryIDCopy(node)
		id1.AddBit(1)
		parents = append(parents, id1)
	}

	//fmt.Println(StringList(parents))
	// if len(parents) > 1 {
	// 	// sort the parent ids
	// 	sort.Slice(parents, func(a, b int) bool {
	// 		return parents[a].GreaterThan(parents[b])
	// 	})
	//}

	// get the byte values of the parents
	return parents, nil
}

// ComputeLabel of a node id
func ComputeLabel(commitment []byte, node *BinaryID, hash HashFunc) []byte {
	parents, _ := GetParents(node)

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
		if _, ok := computed[str]; ok {
			parentLabels = append(parentLabels, computed[str]...)
		} else {
			// compute the label
			label := ComputeLabel(commitment, node, hash)
			// store it in computed
			computed[str] = label
			parentLabels = append(parentLabels, label...)
		}
	}

	result := hash.HashVals(commitment, node.Val, parentLabels)
	return result
}

// ConstructDag create dag
// returns the root hash of the dag as []byte
func ConstructDag(commitment []byte, hash HashFunc) ([]byte, error) {
	// was told no need to use a graph anymore
	// can just compute the edges using an algorithm
	var l0, l1 []byte

	// for height from 0 to m
	for height := 0; height < (m + 1); height++ {
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
			leftLabel := ComputeLabel(commitment, leftId, hash)
			if height == 1 {
				l0 = leftLabel
			}
			WriteToFile(leftLabel)
		}

		// right sub tree
		// pefrom right sub tree calculation
		for level := 0; level < numberOfNodes; level++ {
			rightId, _ := NewBinaryID(uint(height), level)
			rightId.AddBit(1)
			rightLabel := ComputeLabel(commitment, rightId, hash)
			if height == 1 {
				l1 = rightLabel
			}
			WriteToFile(rightLabel)
		}
	}

	rootHash := hash.HashVals(commitment, l0, l1)
	return rootHash, nil
}

// LabelIndex returns the index of a node id
// in the binary file
func LabelIndex(height, nodeSubtreeLen int) int {
	index := (int(math.Pow(float64(2), float64(height+1))) - 1) + nodeSubtreeLen
	return index
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
		var commitParam CommitProofParam
		commitParam.commitment = b
		err = p.CalcCommitProof(commitParam)
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

type CommitProofParam struct {
	commitment []byte
	hash       HashFunc
}

// CalcCommitProof calculates the proof of seqeuntial work
func (p *Prover) CalcCommitProof(param CommitProofParam) error {
	var hashFunction HashFunc

	hashFunction = param.hash
	if hashFunction == nil {
		hashFunction = NewSHA256()
	}
	//
	phi, _ := ConstructDag(param.commitment, hashFunction)
	p.rootHash = phi
	return nil
}

// SendCommitProof send the phi (root Hash) to the verifier
func (p *Prover) SendCommitProof() (b []byte, err error) {
	return p.rootHash, nil
}

// CalcNIPCommitProof proof created by computing openH for the challenge
func (p *Prover) CalcNIPCommitProof(commitment []byte, phi []byte) error {
	var proof []byte
	proof = make([]byte, 32)

	hash := NewSHA256()

	for i := 0; i < t; i++ {
		scParam := make([]byte, binary.MaxVarintLen64)
		binary.BigEndian.PutUint64(scParam, uint64(i))
		proof = append(proof, hash.HashVals(phi, commitment, scParam)...)
	}
	p.challengeProof = proof
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

	gammaIndex := LabelIndex(gammaBinID.Length, len(siblings))

	label_gamma, err := ReadLabelFromFile(gammaIndex)
	if err != nil {
		return err
	}

	proof = append(proof, label_gamma...)

	for i := 0; i < len(siblings); i++ {
		nodeID := siblings[i]
		nodeSiblings, err := Siblings(nodeID)
		if err != nil {
			return err
		}

		siblingIndex := LabelIndex(nodeID.Length, len(nodeSiblings))
		label, err := ReadLabelFromFile(siblingIndex)
		if err != nil {
			return err
		}
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

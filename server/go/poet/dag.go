package poet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

// Siblings returns the list of siblings along the path to the root
//
// Takes in an instance of class BinaryString and returns a list of the
// siblings of the nodes of the path to to root of a binary tree. Also
// returns the node itself, so there are N+1 items in the list for a
// tree with length N.
//
func Siblings(node *BinaryID, left bool) ([]*BinaryID, error) {

	var siblings []*BinaryID
	// Do we really need the node on the siblings list?
	//siblings = append(siblings, node)
	newBinaryID := NewBinaryIDCopy(node)
	for i := 0; i < node.Length; i++ {
		if i == node.Length-1 {
			newBinaryID.FlipBit(newBinaryID.Length)
			// TODO: Add error check
			bit, _ := newBinaryID.GetBit(newBinaryID.Length)
			if (bit == 0) || !(left) {
				siblings = append(siblings, newBinaryID)
			}
		} else {
			id := NewBinaryIDCopy(newBinaryID)
			id.FlipBit(id.Length)
			// TODO: Add error check
			bit, _ := id.GetBit(id.Length)
			if (bit == 0) || !(left) {
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
		left, err := Siblings(node, true)
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

func GammaToBinaryIDs(gamma []byte) ([]*BinaryID, error) {
	var gammas []*BinaryID
	if (len(gamma) % n) != 0 {
		return nil, errors.New(fmt.Sprintf("Gamma wrong length: %v", len(gamma)))
	}
	list_length := len(gamma) / n
	for i := 0; i < list_length; i++ {
		gammas = append(gammas, NewBinaryIDBytes(gamma[i*n:((i+1)*n)]))
	}
	return gammas, nil
}

// CheckAndAdd checks if the BinID is already in the list then adds it if it
// wasn't already there. Also returns true if the BinID was added to the list.
func CheckAndAdd(BinIDs []*BinaryID, BinID *BinaryID) ([]*BinaryID, bool) {
	for _, b := range BinIDs {
		if BinID.Equal(b) {
			return BinIDs, false
		}
	}
	BinIDs = append(BinIDs, BinID)
	return BinIDs, true
}

type ComputeOpts struct {
	Commitment     []byte
	CommitmentHash []byte
	Hash           HashFunc
	Store          StorageIO
}

// ComputeLabel of a node id
func ComputeLabel(node *BinaryID, cOpts *ComputeOpts) []byte {
	parents, _ := GetParents(node)
	var parentLabels []byte
	// Loop through the parents and try to calculate their labels
	// if doesn't exist in computed
	for _, parent := range parents {
		// check if the label exists
		exists, err := cOpts.Store.LabelCalculated(parent)
		if err != nil {
			log.Panic("Error Checking Label: ", err)
		}
		if exists {
			pLabel, err := cOpts.Store.GetLabel(parent)
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
		string(cOpts.Commitment),
		hex.EncodeToString(node.Encode()),
		hex.EncodeToString(parentLabels),
	)

	result := cOpts.Hash.HashVals(
		cOpts.Commitment,
		node.Encode(),
		parentLabels)

	debugLog.Println(
		"Hash for node ",
		string(node.Encode()),
		" calculated: ",
		hex.EncodeToString(result),
	)

	err := cOpts.Store.StoreLabel(node, result)
	if err != nil {
		log.Panic("Error Storing Label: ", err)
	}
	PrintDAG(node, cOpts.Store, "Compute")
	return result
}

func PrintDAG(b *BinaryID, store StorageIO, pre string) {
	if b.Length != n {
		parents, err := GetParents(b)
		if err != nil {
			return
		}
		for _, p := range parents {
			PrintDAG(p, store, pre)
		}
	}
	exists, err := store.LabelCalculated(b)
	if err != nil {
		return
	}
	if exists {
		label, err := store.GetLabel(b)
		if err != nil {
			return
		}
		infoLog.Printf(
			"%v: Node: %v Label: %v",
			pre,
			string(b.Encode()),
			hex.EncodeToString(label),
		)
	}
}

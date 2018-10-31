package poet

import (
	"bytes"
	"testing"
)

var Size int = 32

// This test will test the proper function of Prover with a challenge.
// TODO: Complete Test and provide correct test case. Also need to ensure
// using the correct encoding when converting from string literal to byte slice
// Right now it's using byte literals.
func TestProverWithChallenge(t *testing.T) {
	p := NewProver(false)
	b := []byte{'a', 'b'}
	_, err := p.Write(b)
	if err != nil {
		t.Error("Error Writing Commitment: ", err)
	}
	res := make([]byte, Size)
	_, err = p.Read(res)
	if err != nil {
		t.Error("Error Reading Commitment Proof: ", err)
	}
	expected := []byte{'a', 'b'}
	if !bytes.Equal(res, expected) {
		t.Error("Commitment Proof Not Correct.\nResult: ", res, "\nExpected: ", expected)
	}
}


func TestDagConstructions(t *testing.T) {
	// _ := NewProver(false)
	
}
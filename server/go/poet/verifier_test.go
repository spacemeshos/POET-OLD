package poet

import (
	"fmt"
	"testing"
)

func TestVerifier(t *testing.T) {
	p := NewProver(false)
	v := NewVerifier(p)
	b := []byte{'a', 'b'}
	err := v.Commit(b)
	if err != nil {
		t.Error("Error Sending Commitment: ", err)
	}
	_, err = v.GetCommitProof()
	if err != nil {
		t.Error("Error Getting Commit Proof: ", err)
	}
	_, err = v.SelectChallenge()
	if err != nil {
		t.Error("Error Selecting Challenge: ", err)
	}
	err = v.Challenge()
	if err != nil {
		t.Error("Error Getting Challenge: ", err)
	}
	_, err = v.GetChallengeProof()
	if err != nil {
		t.Error("Error Getting Challenge Proof: ", err)
	}
	fmt.Println("Verifying Challenge Proof", v.challengeProof)
	err = v.VerifyChallengeProof()
	if err != nil {
		t.Error("Error Verifying Challenge Proof: ", err)
	}
}

func TestNIPVerifier(t *testing.T) {
	p := NewProver(true)
	v := NewVerifier(p)
	b := []byte{'a', 'b'}
	err := v.Commit(b)
	if err != nil {
		t.Error("Error Sending Commitment: ", err)
	}
	_, err = v.GetCommitProof()
	if err != nil {
		t.Error("Error Getting Commit Proof: ", err)
	}
	_, err = v.GetChallengeProof()
	if err != nil {
		t.Error("Error Getting Challenge Proof: ", err)
	}
	fmt.Println("Verifying Challenge Proof", v.challengeProof)
	err = v.VerifyChallengeProof()
	if err != nil {
		t.Error("Error Verifying Challenge Proof: ", err)
	}
}

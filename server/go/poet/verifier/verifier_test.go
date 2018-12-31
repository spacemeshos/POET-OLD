package verifier

import (
	"testing"

	"github.com/SomniaStellarum/POET/server/go/poet"
)

func TestVerifier(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("skipping testing in short mode")
	// }
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
	test_n := 4
	p := poet.NewProver() // Should declare Prover with n
	v := NewVerifier(test_n)
	v.SetHash(poet.NewSHA256())
	b := []byte{'a', 'b'}
	err := p.CalcCommitProof(b)
	if err != nil {
		t.Error("Error Sending Commitment: ", err)
	}
	roothash, err := p.CommitProof()
	if err != nil {
		t.Error("Error Getting Commit Proof: ", err)
	}
	v.SetCommitmentProof(roothash)
	challenge, err := v.SelectRndChallenge()
	if err != nil {
		t.Error("Error Selecting Challenge: ", err)
	}
	v.SetChallenge(challenge)
	err = p.CalcChallengeProof(challenge)
	if err != nil {
		t.Error("Error Sending Challenge: ", err)
	}
	challengeProof, err := p.ChallengeProof()
	if err != nil {
		t.Error("Error Getting Challenge Proof: ", err)
	}
	v.SetChallengeProof(challengeProof)
	err = v.VerifyChallengeProof()
	if err != nil {
		t.Error("Error Verifying Challenge Proof: ", err)
	}
}

// func TestNIPVerifier(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping testing in short mode")
// 	}
// 	// debugLog.SetOutput(os.Stdout)
// 	// defer debugLog.SetOutput(ioutil.Discard)
// 	p := poet.NewProver()
// 	test_n := 4
// 	v := NewVerifier(p, test_n)
// 	b := []byte{'a', 'b'}
// 	err := v.Commit(b)
// 	if err != nil {
// 		t.Error("Error Sending Commitment: ", err)
// 	}
// 	_, err = v.GetCommitProof()
// 	if err != nil {
// 		t.Error("Error Getting Commit Proof: ", err)
// 	}
// 	_, err = v.GetChallengeProof()
// 	if err != nil {
// 		t.Error("Error Getting Challenge Proof: ", err)
// 	}
// 	//fmt.Println("Verifying Challenge Proof", v.challengeProof)
// 	err = v.VerifyChallengeProof()
// 	if err != nil {
// 		t.Error("Error Verifying Challenge Proof: ", err)
// 	}
// }

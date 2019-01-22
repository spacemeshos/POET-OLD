package verifier

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/spacemeshos/POET/server/go/poet"
)

func TestVerifier(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	// debugLog.SetOutput(os.Stdout)
	// poet.DebugOutput()
	// defer debugLog.SetOutput(ioutil.Discard)
	// defer poet.DebugSupress()
	test_n := 4
	p := poet.NewProver() // Should declare Prover with n
	p.ChangeDAGSize(test_n)
	v := NewVerifier(test_n)
	v.SetHash(poet.NewSHA256())
	b := []byte{'a', 'b'}
	v.SetCommitment(b)
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

var verifierTests = []struct {
	n              int
	commitment     []byte
	challenge      []byte
	challengeProof []byte
	expectedRoot   []byte
}{
	{
		5,
		[]byte("this is a commitment"),
		[]byte("00101"),
		[]byte("eee82defeda4a10a839dd07e333de474968850f3e1c70b6edec428086089b631760694739fa4b844d96db52e7e1b719d31801476a5eb49f11480e3fd06050462e6e6ee64bcd8f73351a3c067a2d394bae994ed9b708fe26ed08c07fbf8f6f78beceaafe8393d22e8ad5a29983ea47e5ee1006d7fab210dea4940e9807c9994329c7c51af1a7f237363ec1221391bd3f2a4f3c1a7b60cf75fe87e267497418ac5"),
		[]byte("f1418ee0a1c3cd9b8a248334f2549a78bb967a4796efd638b870a0434b479254"),
	},
}

func TestVerifierToChallenge(t *testing.T) {
	var err error
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	// debugLog.SetOutput(os.Stdout)
	// poet.DebugOutput()
	// defer debugLog.SetOutput(ioutil.Discard)
	// defer poet.DebugSupress()
	p := poet.NewProver() // Should declare Prover with n
	v := NewVerifier(4)
	v.SetHash(poet.NewSHA256())
	for _, vt := range verifierTests {
		v.SetDAGSize(vt.n)
		p.ChangeDAGSize(vt.n)
		v.SetCommitment(vt.commitment)
		phi := make([]byte, hex.DecodedLen(len(vt.expectedRoot)))
		_, err = hex.Decode(phi, vt.expectedRoot)
		if err != nil {
			log.Fatal(err)
		}
		v.SetCommitmentProof(phi)
		v.SetChallenge(vt.challenge)
		proof := make([][]byte, 0)
		proof = append(proof, vt.challengeProof)
		v.SetChallengeProof(proof)
		err = v.VerifyChallengeProof()
		if err != nil {
			t.Errorf("Error Verifying Challenge: %v", err)
		}
	}
}

func TestNIPVerifier(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	// debugLog.SetOutput(os.Stdout)
	// poet.DebugOutput()
	// defer debugLog.SetOutput(ioutil.Discard)
	// defer poet.DebugSupress()
	test_n := 4
	p := poet.NewProver() // Should declare Prover with n
	p.ChangeDAGSize(test_n)
	v := NewVerifier(test_n)
	v.SetHash(poet.NewSHA256())
	b := []byte{'a', 'b'}
	v.SetCommitment(b)
	err := p.CalcCommitProof(b)
	if err != nil {
		t.Error("Error Sending Commitment: ", err)
	}
	roothash, err := p.CommitProof()
	if err != nil {
		t.Error("Error Getting Commit Proof: ", err)
	}
	v.SetCommitmentProof(roothash)
	err = p.CalcNIPCommitProof()
	if err != nil {
		t.Error("Error Calculating NIP Challenge: ", err)
	}
	challengeProof, err := p.ChallengeProof()
	if err != nil {
		t.Error("Error Getting Challenge Proof: ", err)
	}
	v.SetChallengeProof(challengeProof)
	p.ShowDAG()
	err = v.VerifyChallengeProof()
	if err != nil {
		t.Error("Error Verifying Challenge Proof: ", err)
	}
}

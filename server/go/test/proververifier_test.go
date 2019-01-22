package poet

import (
	"encoding/hex"
	"fmt"
	"github.com/spacemeshos/POET/server/go/poet"
	"github.com/spacemeshos/POET/server/go/poet/verifier"
	"testing"
)

var proverVerifierTest = struct {
	n          int
	commitment []byte
}{n: 8, commitment: []byte("this is a commitment")}

func TestProverVerifier(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	n := proverVerifierTest.n
	commitment := proverVerifierTest.commitment

	p := poet.NewProver()
	p.ChangeDAGSize(n)

	v := verifier.NewVerifier(n)
	v.SetHash(poet.NewSHA256())
	v.SetCommitment(commitment)

	// calculating phi

	err := p.CalcCommitProof(commitment)
	if err != nil {
		t.Error("Error calculating commit proof: ", err)
	}

	phi, err := p.CommitProof()
	if err != nil {
		t.Error("Error returning commit proof: ", err)
	}

	fmt.Print("phi: ", hex.EncodeToString(phi), "\n")

	// generating a random challenge

	v.SetCommitmentProof(phi)

	challenge, err := v.SelectRndChallenge()
	if err != nil {
		t.Error("Error generating random challenge")
	}
	if (len(challenge) % n) != 0 {
		t.Error("Random challenge wrong size")
	}

	// calculating a proof for the random challenge

	err = p.CalcChallengeProof(challenge)
	if err != nil {
		t.Error("Error calculating challenge proof: ", err)
	}
	proof, err := p.ChallengeProof()
	if err != nil {
		t.Error("Error returning challenge proof: ", err)
	}

	// verifying the proof

	v.SetChallengeProof(proof)
	err = v.VerifyChallengeProof()
	if err != nil {
		t.Error("Error verifying challenge proof: ", err)
	}
}

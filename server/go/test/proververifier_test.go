package poet

import (
	"encoding/hex"
	"fmt"
	"github.com/spacemeshos/POET/server/go/poet"
	"github.com/spacemeshos/POET/server/go/poet/verifier"
	"testing"
	"time"
)

var proverVerifierTest = struct {
	n          int
	commitment []byte
}{n: 25, commitment: []byte("this is a commitment")}

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

	fmt.Printf("Generated proof for n = %d\n", n)

	// calculating phi
	t1 := time.Now()

	err := p.CalcCommitProof(commitment)
	if err != nil {
		t.Error("Error calculating commit proof: ", err)
	}
	e := time.Since(t1)

	fmt.Printf("Generated proof in %s (%f seconds)\n", e, e.Seconds())

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

	t1 = time.Now()
	err = v.VerifyChallengeProof()
	e = time.Since(t1)

	if err != nil {
		t.Error("Error verifying challenge proof: ", err)
	}

	fmt.Printf("Verified proof in %s (%f seconds)\n", e, e.Seconds())

}

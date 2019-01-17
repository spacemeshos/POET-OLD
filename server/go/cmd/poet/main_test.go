package main

import (
	"context"
	"errors"
	"fmt"
	api "github.com/spacemeshos/poet-core-api"
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"testing"
	"time"
)

const (
	BlackboxRPCPort          = "50052"
	BlackboxRPCHostPort      = "34.73.172.126:" + BlackboxRPCPort
)

var mainTests = []struct {
	commitment []byte
	n          uint32
	h          string
}{
	{[]byte("this is a commitment"), 4, "sha256"},
}

func TestPoetNIPMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	prover, cleanup := api.NewProverClient(":8888")
	defer cleanup()
	verifier, cleanup := api.NewVerifierClient(":8888")
	defer cleanup()

	if err := nipProofTests(prover, verifier); err != nil {
		t.Fatal(err)
	}
}

func TestPoetChallengeMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	prover, cleanup := api.NewProverClient(":8888")
	defer cleanup()
	verifier, cleanup := api.NewVerifierClient(":8888")
	defer cleanup()

	if err := challengeProofTests(prover, verifier); err != nil {
		t.Fatal(err)
	}
}

// Testing Verifier against the black box implementation
func TestPoetMainNIPVeriferRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	prover, cleanup := api.NewProverClient(":8888")
	defer cleanup()
	verifier, cleanup := api.NewVerifierClient(BlackboxRPCHostPort)
	defer cleanup()

	if err := nipProofTests(prover, verifier); err != nil {
		t.Fatal(err)
	}
}

// Testing Prover against the black box implementation
func TestPoetMainNIPProverRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	prover, cleanup := api.NewProverClient(BlackboxRPCHostPort)
	defer cleanup()
	verifier, cleanup := api.NewVerifierClient(":8888")
	defer cleanup()

	if err := nipProofTests(prover, verifier); err != nil {
		t.Fatal(err)
	}
}

// Testing Verifier against the black box implementation
func TestPoetMainChallengeVeriferRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	prover, cleanup := api.NewProverClient(":8888")
	defer cleanup()
	verifier, cleanup := api.NewVerifierClient(BlackboxRPCHostPort)
	defer cleanup()

	if err := challengeProofTests(prover, verifier); err != nil {
		t.Fatal(err)
	}
}

// Testing Prover against the black box implementation
func TestPoetMainChallengeProverRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	prover, cleanup := api.NewProverClient(BlackboxRPCHostPort)
	defer cleanup()
	verifier, cleanup := api.NewVerifierClient(":8888")
	defer cleanup()

	if err := challengeProofTests(prover, verifier); err != nil {
		t.Fatal(err)
	}
}

func nipProofTests(prover pcrpc.PoetCoreProverClient, verifier pcrpc.PoetVerifierClient) error {
	for _, mTest := range mainTests {
		ctx := context.Background()
		dag := &pcrpc.DagParams{X: mTest.commitment, N: mTest.n, H: mTest.h}
		_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: dag})
		defer prover.Shutdown(ctx, &pcrpc.ShutdownRequest{})
		if err != nil {
			return err
		}

		// verify NIP

		nipRes, err := prover.GetNIP(ctx, &pcrpc.GetNIPRequest{})
		if err != nil {
			return err
		}

		verifyNIPRes, err := verifier.VerifyNIP(ctx, &pcrpc.VerifyNIPRequest{D: dag, P: nipRes.Proof})
		if err != nil {
			return err
		}
		if !verifyNIPRes.Verified {
			// Should test the rest of the tests then return the error?
			return errors.New("NIP wasn't verified.")
		}
		prover.Clean(ctx, &pcrpc.CleanRequest{})
	}
	return nil
}

func challengeProofTests(prover pcrpc.PoetCoreProverClient, verifier pcrpc.PoetVerifierClient) error {
	for _, mTest := range mainTests {
		ctx := context.Background()
		dag := &pcrpc.DagParams{X: mTest.commitment, N: mTest.n, H: mTest.h}
		_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: dag})
		defer prover.Shutdown(ctx, &pcrpc.ShutdownRequest{})
		if err != nil {
			return err
		}

		// Get Challenge
		fmt.Println("Getting Challenge")

		challenge, err := verifier.GetRndChallenge(ctx, &pcrpc.GetRndChallengeRequest{D: dag})
		if err != nil {
			return err
		}

		fmt.Println(challenge)
		fmt.Println("Getting Proof")

		res, err := prover.GetProof(ctx, &pcrpc.GetProofRequest{C: challenge.C})
		if err != nil {
			return err
		}

		fmt.Println("Verifying Proof")

		verifyRes, err := verifier.VerifyProof(ctx, &pcrpc.VerifyProofRequest{D: dag, P: res.Proof, C: challenge.C})
		if err != nil {
			return err
		}
		if !verifyRes.Verified {
			// Should test the rest of the tests then return the error?
			return errors.New("Verifier Challenge wasn't verified.")
		}
		prover.Clean(ctx, &pcrpc.CleanRequest{})
	}
	return nil
}

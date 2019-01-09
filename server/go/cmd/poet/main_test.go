package main

import (
	"context"
	"testing"
	"time"

	"github.com/spacemeshos/poet-core-api/pcrpc"
	"google.golang.org/grpc"
)

var mainTests = []struct {
	commitment []byte
	n          uint32
	h          string
}{
	{[]byte("this is a commitment"), 4, "sha256"},
}

func TestPoetMain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error Dialing: %v", err)
	}
	defer conn.Close()

	prover := pcrpc.NewPoetCoreProverClient(conn)
	verifier := pcrpc.NewPoetVerifierClient(conn)

	for _, mTest := range mainTests {
		ctx := context.Background()
		dag := &pcrpc.DagParams{X: mTest.commitment, N: mTest.n, H: mTest.h}
		_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: dag})
		defer prover.Clean(ctx, &pcrpc.CleanRequest{})
		if err != nil {
			t.Fatal(err)
		}

		// verify NIP

		nipRes, err := prover.GetNIP(ctx, &pcrpc.GetNIPRequest{})
		if err != nil {
			t.Fatal(err)
		}

		verifyNIPRes, err := verifier.VerifyNIP(ctx, &pcrpc.VerifyNIPRequest{D: dag, P: nipRes.Proof})
		if err != nil {
			t.Fatal(err)
		}
		if !verifyNIPRes.Verified {
			t.Fatal("NIP wasn't verified.")
		}
	}
}

// Testing Verifier against the black box implementation
func TestPoetMainVeriferRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	connVerifier, err := grpc.Dial("35.196.137.245:50052", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error Dialing: %v", err)
	}
	defer connVerifier.Close()

	connProver, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error Dialing: %v", err)
	}
	defer connProver.Close()

	prover := pcrpc.NewPoetCoreProverClient(connProver)
	verifier := pcrpc.NewPoetVerifierClient(connVerifier)

	for _, mTest := range mainTests {
		ctx := context.Background()
		dag := &pcrpc.DagParams{X: mTest.commitment, N: mTest.n, H: mTest.h}
		_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: dag})
		defer prover.Clean(ctx, &pcrpc.CleanRequest{})
		if err != nil {
			t.Fatal(err)
		}

		// verify NIP

		nipRes, err := prover.GetNIP(ctx, &pcrpc.GetNIPRequest{})
		if err != nil {
			t.Fatal(err)
		}

		verifyNIPRes, err := verifier.VerifyNIP(ctx, &pcrpc.VerifyNIPRequest{D: dag, P: nipRes.Proof})
		if err != nil {
			t.Fatal(err)
		}
		if !verifyNIPRes.Verified {
			t.Fatal("NIP wasn't verified.")
		}
	}
}

// Testing Prover against the black box implementation
func TestPoetMainProverRPC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	go poetMain()

	// Might need a pause to let main finish setup and start listening. To test.
	time.Sleep(5 * time.Second)

	connProver, err := grpc.Dial("35.196.137.245:50052", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error Dialing: %v", err)
	}
	defer connProver.Close()

	connVerifier, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Error Dialing: %v", err)
	}
	defer connVerifier.Close()

	prover := pcrpc.NewPoetCoreProverClient(connProver)
	verifier := pcrpc.NewPoetVerifierClient(connVerifier)

	for _, mTest := range mainTests {
		ctx := context.Background()
		dag := &pcrpc.DagParams{X: mTest.commitment, N: mTest.n, H: mTest.h}
		_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: dag})
		defer prover.Clean(ctx, &pcrpc.CleanRequest{})
		if err != nil {
			t.Fatal(err)
		}

		// verify NIP

		nipRes, err := prover.GetNIP(ctx, &pcrpc.GetNIPRequest{})
		if err != nil {
			t.Fatal(err)
		}

		verifyNIPRes, err := verifier.VerifyNIP(ctx, &pcrpc.VerifyNIPRequest{D: dag, P: nipRes.Proof})
		if err != nil {
			t.Fatal(err)
		}
		if !verifyNIPRes.Verified {
			t.Fatal("NIP wasn't verified.")
		}
	}
}

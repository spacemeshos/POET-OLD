package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/SomniaStellarum/POET/server/go/poet"
	"github.com/SomniaStellarum/POET/server/go/poet/verifier"
	"github.com/spacemeshos/poet-core-api/pcrpc"
)

type VerifierServer struct {
	verifier *verifier.Verifier
}

func NewVerifierServer(v *verifier.Verifier) (vs *VerifierServer) {
	return &VerifierServer{verifier: v}
}

func (vs *VerifierServer) VerifyProof(ctx context.Context, proofRequest *pcrpc.VerifyProofRequest) (*pcrpc.VerifyProofResponse, error) {
	res := new(pcrpc.VerifyProofResponse)
	vs.verifier.SetDAGSize(int(proofRequest.D.N))
	switch proofRequest.D.H {
	case "sha256":
		vs.verifier.SetHash(poet.NewSHA256())
		// TODO: Add scrypt and other Hash functions as needed
	default:
		return nil, errors.New(
			fmt.Sprintf("Verifier does not implement Hash Function %v", proofRequest.D.H),
		)
	}
	vs.verifier.SetCommitment(proofRequest.D.X)
	var challenge []byte
	for _, s := range proofRequest.C {
		challenge = append(challenge, []byte(s)...)
	}
	vs.verifier.SetChallenge(challenge)
	vs.verifier.SetCommitmentProof(proofRequest.P.Phi)
	var proof [][]byte
	for _, l := range proofRequest.P.L {
		var cProof []byte
		for _, b := range l.Labels {
			cProof = append(cProof, b...)
		}
		proof = append(proof, cProof)
	}
	vs.verifier.SetChallengeProof(proof)
	err := vs.verifier.VerifyChallengeProof()
	// TODO: check for other errs before returning nil error
	if err != nil {
		res.Verified = false
		return res, nil
	}
	res.Verified = true
	return res, nil
}

func (vs *VerifierServer) VerifyNIP(ctx context.Context, proofRequest *pcrpc.VerifyNIPRequest) (*pcrpc.VerifyNIPResponse, error) {
	res := new(pcrpc.VerifyNIPResponse)
	vs.verifier.SetDAGSize(int(proofRequest.D.N))
	switch proofRequest.D.H {
	case "sha256":
		vs.verifier.SetHash(poet.NewSHA256())
		// TODO: Add scrypt and other Hash functions as needed
	default:
		return nil, errors.New(
			fmt.Sprintf("Verifier does not implement Hash Function %v", proofRequest.D.H),
		)
	}
	vs.verifier.SetCommitment(proofRequest.D.X)
	vs.verifier.SetCommitmentProof(proofRequest.P.Phi)
	var proof [][]byte
	for _, l := range proofRequest.P.L {
		var cProof []byte
		for _, b := range l.Labels {
			cProof = append(cProof, b...)
		}
		proof = append(proof, cProof)
	}
	vs.verifier.SetChallengeProof(proof)
	err := vs.verifier.VerifyChallengeProof()
	// TODO: check for other errs before returning nil error
	if err != nil {
		res.Verified = false
		return res, nil
	}
	res.Verified = true
	return res, nil
}

func (vs *VerifierServer) GetRndChallenge(ctx context.Context, rndChallengeRequest *pcrpc.GetRndChallengeRequest) (*pcrpc.GetRndChallengeResponse, error) {
	res := new(pcrpc.GetRndChallengeResponse)
	n := int(rndChallengeRequest.D.N)
	vs.verifier.SetDAGSize(n)
	vs.verifier.SetCommitment(rndChallengeRequest.D.X)
	switch rndChallengeRequest.D.H {
	case "sha256":
		vs.verifier.SetHash(poet.NewSHA256())
		// TODO: Add scrypt and other Hash functions as needed
	default:
		return nil, errors.New(
			fmt.Sprintf("Verifier does not implement Hash Function %v", rndChallengeRequest.D.H),
		)
	}
	b, err := vs.verifier.SelectRndChallenge()
	if err != nil {
		return nil, err
	}
	if (len(b) % n) != 0 {
		return nil, errors.New("Random Challenge wrong size")
	}
	for i := 0; i < (len(b) / n); i++ {
		res.C = append(res.C, string(b[i*n:(i+1)*n-1]))
	}
	return res, nil
}

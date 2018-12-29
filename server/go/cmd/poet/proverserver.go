package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/SomniaStellarum/POET/server/go/poet"
	"github.com/spacemeshos/poet-core-api/pcrpc"
)

type ProverServer struct {
	prover *poet.Prover
}

func NewProverServer(p *poet.Prover) (ps *ProverServer) {
	return &ProverServer{prover: p}
}

func (ps *ProverServer) Compute(ctx context.Context, computeRequest *pcrpc.ComputeRequest) (*pcrpc.ComputeResponse, error) {
	ps.prover.ChangeDAGSize(int(computeRequest.D.N))
	switch computeRequest.D.H {
	case "sha256":
		ps.prover.ChangeHashFunc(poet.NewSHA256())
		// TODO: Add scrypt and other Hash functions as needed
	default:
		return nil, errors.New(
			fmt.Sprintf("Prover does not implement Hash Function %v", computeRequest.D.H),
		)
	}
	err := ps.prover.CalcCommitProof(computeRequest.D.X)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error Writing in Prover %v", err))
	}
	// res := make([]byte, poet.HashSize)
	res, err := ps.prover.CommitProof()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error Reading in Prover %v", err))
	}
	computeResponse := new(pcrpc.ComputeResponse)
	computeResponse.Phi = res
	return computeResponse, nil
}

func (ps *ProverServer) GetNIP(ctx context.Context, nipRequest *pcrpc.GetNIPRequest) (*pcrpc.GetNIPResponse, error) {
	err := ps.prover.CalcNIPCommitProof()
	if err != nil {
		return nil, err
	}
	b, err := ps.prover.ChallengeProof()
	if err != nil {
		return nil, err
	}
	nipResponse := new(pcrpc.GetNIPResponse)
	nipResponse.Proof.Phi, err = ps.prover.CommitProof()
	if err != nil {
		return nil, err
	}
	nipResponse.Proof.L, err = GetLabels(b)
	if err != nil {
		return nil, err
	}
	return nipResponse, nil
}

func (ps *ProverServer) GetProof(ctx context.Context, proofRequest *pcrpc.GetProofRequest) (*pcrpc.GetProofResponse, error) {
	var gamma []byte
	for _, s := range proofRequest.C {
		gamma = append(gamma, []byte(s)...)
	}
	err := ps.prover.CalcChallengeProof(gamma)
	if err != nil {
		return nil, err
	}
	b, err := ps.prover.ChallengeProof()
	if err != nil {
		return nil, err
	}
	proofResponse := new(pcrpc.GetProofResponse)
	proofResponse.Proof.Phi, err = ps.prover.CommitProof()
	if err != nil {
		return nil, err
	}
	proofResponse.Proof.L, err = GetLabels(b)
	if err != nil {
		return nil, err
	}
	return proofResponse, nil
}

func (ps *ProverServer) Clean(context.Context, *pcrpc.CleanRequest) (*pcrpc.CleanResponse, error) {
	ps.prover.Clean()
	res := new(pcrpc.CleanResponse)
	return res, nil
}

func (ps *ProverServer) Shutdown(context.Context, *pcrpc.ShutdownRequest) (*pcrpc.ShutdownResponse, error) {
	res := new(pcrpc.ShutdownResponse)
	return res, nil
}

func GetLabels(b [][]byte) ([]*pcrpc.Labels, error) {
	var res []*pcrpc.Labels
	for _, bi := range b {
		if (len(bi) % poet.HashSize) != 0 {
			return nil, errors.New("Byte slice not multiple of hash size. Cannot Send Proof")
		}
		num := len(bi) / poet.HashSize
		l := new(pcrpc.Labels)
		for i := 0; i < num; i++ {
			l.Labels = append(l.Labels, bi[i*poet.HashSize:((i+1)*poet.HashSize-1)])
		}
		res = append(res, l)
	}
	return res, nil
}

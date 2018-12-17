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
	_, err := ps.prover.Write(computeRequest.D.X)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error Writing in Prover %v", err))
	}
	res := make([]byte, poet.HashSize)
	_, err = ps.prover.Read(res)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error Reading in Prover %v", err))
	}
	computeResponse := new(pcrpc.ComputeResponse)
	computeResponse.Phi = res
	return computeResponse, nil
}

func (ps *ProverServer) GetNIP(context.Context, *pcrpc.GetNIPRequest) (*pcrpc.GetNIPResponse, error) {
	return nil, errors.New("Not implemented")
}

func (ps *ProverServer) GetProof(context.Context, *pcrpc.GetProofRequest) (*pcrpc.GetProofResponse, error) {
	return nil, errors.New("Not implemented")
}

func (ps *ProverServer) Clean(context.Context, *pcrpc.CleanRequest) (*pcrpc.CleanResponse, error) {
	return nil, errors.New("Not implemented")
}

func (ps *ProverServer) Shutdown(context.Context, *pcrpc.ShutdownRequest) (*pcrpc.ShutdownResponse, error) {
	return nil, errors.New("Not implemented")
}

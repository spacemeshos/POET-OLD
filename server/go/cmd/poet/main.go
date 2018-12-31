package main

import (
	"log"
	"net"

	"github.com/SomniaStellarum/POET/server/go/poet"
	"github.com/SomniaStellarum/POET/server/go/poet/verifier"
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	p := poet.NewProver()
	proverServer := NewProverServer(p)
	// TODO: Setup correct Default DAG Size
	v := verifier.NewVerifier(4)
	verifierServer := NewVerifierServer(v)
	server := grpc.NewServer()
	pcrpc.RegisterPoetCoreProverServer(server, proverServer)
	pcrpc.RegisterPoetVerifierServer(server, verifierServer)
	log.Fatal(server.Serve(lis))
}

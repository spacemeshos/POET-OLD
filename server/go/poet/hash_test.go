package poet

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/facundomedica/sha256-simd"
	miniosha "github.com/minio/sha256-simd"
)

func BenchmarkSha256(t *testing.B) {

	buff := bytes.Buffer{}
	buff.Write([]byte("Seed data goes here"))
	out := [32]byte{}
	n := uint64(math.Pow(10, 8))

	fmt.Printf("Computing %d serial sha-256...\n", n)

	t1 := time.Now().Unix()

	for i := uint64(0); i < n; i++ {
		out = miniosha.Sum256(buff.Bytes())
		buff.Reset()
		buff.Write(out[:])
	}

	d := time.Now().Unix() - t1
	r := n / uint64(d)
	fmt.Printf("Final hash: %x. Running time: %d secs. Hash-rate:%d hashes-per-sec\n", buff.Bytes(), d, r)
}

func BenchmarkSha256Ex(t *testing.B) {

	buff := bytes.Buffer{}
	buff.Write([]byte("Seed data goes here"))
	out := [32]byte{}
	n := uint64(math.Pow(10, 8))

	fmt.Printf("Computing %d serial sha-256...\n", n)

	t1 := time.Now().Unix()

	for i := uint64(0); i < n; i++ {
		out = sha256.Sum256(buff.Bytes())
		buff.Reset()
		buff.Write(out[:])
	}

	d := time.Now().Unix() - t1
	r := n / uint64(d)
	fmt.Printf("Final hash: %x. Running time: %d secs. Hash-rate:%d hashes-per-sec\n", buff.Bytes(), d, r)
}

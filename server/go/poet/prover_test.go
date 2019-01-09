package poet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

var Size int = 32

var proverTests = []struct {
	n            int
	commitment   []byte
	expectedRoot []byte
}{
	// {5, []byte("this is a commitment"), []byte("f1418ee0a1c3cd9b8a248334f2549a78bb967a4796efd638b870a0434b479254")},
	// {8, []byte("this is a commitment"), []byte("ff9aae0df334a9d17e6efb20fbf293872df5ce307b134582b70d139f2f722d80")},
	// {9, []byte("this is a commitment"), []byte("5e164427c1e77e00c6718b5672e02adfc24a5ad3b0c146957ec360c6f7420d7a")},
	// {10, []byte("this is a commitment"), []byte("504670627c97ec74c72bd814f94a2e191931bda81e046009e5f2f66f09309827")},
	{17, []byte("this is a commitment"), []byte("0000000000000000000000000000000000000000000000000000000000000000")},
}

// This test will test the proper function of Prover with a challenge.
// TODO: Complete Test and provide correct test case. Also need to ensure
// using the correct encoding when converting from string literal to byte slice
// Right now it's using byte literals.
func TestProverWithChallenge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	// debugLog.SetOutput(logFile)
	// defer debugLog.SetOutput(ioutil.Discard)
	for _, pTest := range proverTests {
		n = pTest.n
		p := NewProver()
		err := p.CalcCommitProof(pTest.commitment)
		if err != nil {
			t.Error("Error Calculating Commit Proof: ", err)
		}
		res, err := p.CommitProof()
		if err != nil {
			t.Error("Error Returning Commit Proof: ", err)
		}
		expected := make([]byte, hex.DecodedLen(len(pTest.expectedRoot)))
		_, err = hex.Decode(expected, pTest.expectedRoot)
		if err != nil {
			log.Fatal(err)
		}
		if !bytes.Equal(res, expected) {
			t.Errorf(
				"Commitment Proof Not Correct. n=%v\nResult: %v\nExpected: %v\n",
				pTest.n,
				hex.EncodeToString(res),
				hex.EncodeToString(expected),
			)
		}
		p.CalcChallengeProof([]byte("00101"))
		b, _ := p.ChallengeProof()
		fmt.Print("Proof: ", hex.EncodeToString(b[0]), "\n")
	}
}

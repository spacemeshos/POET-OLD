package poet

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"
)

var Size int = 32

// This test will test the proper function of Prover with a challenge.
// TODO: Complete Test and provide correct test case. Also need to ensure
// using the correct encoding when converting from string literal to byte slice
// Right now it's using byte literals.
func TestProverWithChallenge(t *testing.T) {
	n = 5
	p := NewProver(false)
	b := []byte("this is a commitment")
	_, err := p.Write(b)
	if err != nil {
		t.Error("Error Writing Commitment: ", err)
	}
	res := make([]byte, Size)
	_, err = p.Read(res)
	if err != nil {
		t.Error("Error Reading Commitment Proof: ", err)
	}
	src := []byte("f1418ee0a1c3cd9b8a248334f2549a78bb967a4796efd638b870a0434b479254")
	expected := make([]byte, hex.DecodedLen(len(src)))
	_, err = hex.Decode(expected, src)
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(res, expected) {
		t.Error("Commitment Proof Not Correct.\nResult: ", hex.EncodeToString(res), "\nExpected: ", hex.EncodeToString(expected))
	}
}

// func TestDagConstructions(t *testing.T) {
// 	// _ := NewProver(false)
// 	b, _ := NewBinaryID(0, 0)
// 	p1s, _ := GetParents(b)
// 	p2s, _ := GetParents(p1s[0])
// 	p3s, _ := GetParents(p2s[0])
// 	p4s, _ := GetParents(p3s[0])
// 	p5s, _ := GetParents(p4s[1])
// 	fmt.Println(b)
// 	fmt.Println(p1s)
// 	fmt.Println(p2s)
// 	fmt.Println(p3s)
// 	fmt.Println(p4s)
// 	fmt.Println(p5s)
//
// 	t.Errorf("Error creating DAG\n")
// }

var siblingsTests = []struct {
	in       *BinaryID
	expected []*BinaryID
}{
	{in: &BinaryID{Length: 3, Val: []byte{byte(7)}},
		expected: []*BinaryID{
			&BinaryID{Length: 3, Val: []byte{byte(6)}},
			&BinaryID{Length: 2, Val: []byte{byte(2)}},
			&BinaryID{Length: 1, Val: []byte{byte(0)}},
		},
	},
	{in: &BinaryID{Length: 4, Val: []byte{byte(15)}},
		expected: []*BinaryID{
			&BinaryID{Length: 4, Val: []byte{byte(14)}},
			&BinaryID{Length: 3, Val: []byte{byte(6)}},
			&BinaryID{Length: 2, Val: []byte{byte(2)}},
			&BinaryID{Length: 1, Val: []byte{byte(0)}},
		},
	},
}

func TestSiblings(t *testing.T) {
	// Set n to known value for test
	n = 4
	for _, s := range siblingsTests {
		actual, err := Siblings(s.in)
		if err != nil {
			t.Errorf("Error returned from Siblings. Error: %v\n", err)
		}
		if !(BinaryIDListEqual(actual, s.expected)) {
			t.Errorf(
				"Siblings Failed\nExpected:\n%v\nActual:\n%v\n",
				StringList(s.expected),
				StringList(actual),
			)
		}
	}
}

var parentsTests = []struct {
	in       *BinaryID
	expected []*BinaryID
}{
	{in: &BinaryID{Length: 3, Val: []byte{byte(7)}},
		expected: []*BinaryID{
			&BinaryID{Length: 4, Val: []byte{byte(14)}},
			&BinaryID{Length: 4, Val: []byte{byte(15)}},
		},
	},
	{in: &BinaryID{Length: 4, Val: []byte{byte(6)}},
		expected: []*BinaryID{
			&BinaryID{Length: 3, Val: []byte{byte(2)}},
			&BinaryID{Length: 2, Val: []byte{byte(0)}},
		},
	},
	{in: &BinaryID{Length: 0, Val: []byte{byte(0)}},
		expected: []*BinaryID{
			&BinaryID{Length: 1, Val: []byte{byte(0)}},
			&BinaryID{Length: 1, Val: []byte{byte(1)}},
		},
	},
	{in: &BinaryID{Length: 4, Val: []byte{byte(0)}},
		expected: []*BinaryID{},
	},
}

func TestGetParents(t *testing.T) {
	// Set n to known value for test
	n = 4
	for _, p := range parentsTests {
		actual, err := GetParents(p.in)
		if err != nil {
			t.Errorf("Error returned from GetParents. Error: %v\n", err)
		}
		if !(BinaryIDListEqual(actual, p.expected)) {
			t.Errorf(
				"GetParents Failed\nExpected:\n%v\nActual:\n%v\n",
				StringList(p.expected),
				StringList(actual),
			)
		}
	}
}

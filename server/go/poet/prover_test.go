package poet

import (
	"testing"
)

var Size int = 32

// This test will test the proper function of Prover with a challenge.
// TODO: Complete Test and provide correct test case. Also need to ensure
// using the correct encoding when converting from string literal to byte slice
// Right now it's using byte literals.
// func TestProverWithChallenge(t *testing.T) {
// 	p := NewProver(false)
// 	b := []byte{'a', 'b'}
// 	_, err := p.Write(b)
// 	if err != nil {
// 		t.Error("Error Writing Commitment: ", err)
// 	}
// 	res := make([]byte, Size)
// 	_, err = p.Read(res)
// 	if err != nil {
// 		t.Error("Error Reading Commitment Proof: ", err)
// 	}
// 	expected := []byte{'a', 'b'}
// 	if !bytes.Equal(res, expected) {
// 		t.Error("Commitment Proof Not Correct.\nResult: ", res, "\nExpected: ", expected)
// 	}
// }

func TestDagConstructions(t *testing.T) {
	// _ := NewProver(false)

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

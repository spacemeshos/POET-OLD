package poet

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var Size int = 32

var proverTests = []struct {
	n            int
	commitment   []byte
	expectedRoot []byte
}{
	//	{5, []byte("this is a commitment"), []byte("f1418ee0a1c3cd9b8a248334f2549a78bb967a4796efd638b870a0434b479254")},
	//	{8, []byte("this is a commitment"), []byte("ff9aae0df334a9d17e6efb20fbf293872df5ce307b134582b70d139f2f722d80")},
	{9, []byte("this is a commitment"), []byte("5e164427c1e77e00c6718b5672e02adfc24a5ad3b0c146957ec360c6f7420d7a")},
	//	{10, []byte("this is a commitment"), []byte("504670627c97ec74c72bd814f94a2e191931bda81e046009e5f2f66f09309827")},
	//	{7, []byte("this is a commitment"), []byte("0000000000000000000000000000000000000000000000000000000000000000")},
}

// This test will test the proper function of Prover with a challenge.
// TODO: Complete Test and provide correct test case. Also need to ensure
// using the correct encoding when converting from string literal to byte slice
// Right now it's using byte literals.
func TestProverWithChallenge(t *testing.T) {
	debugLog.SetOutput(logFile)
	defer debugLog.SetOutput(ioutil.Discard)
	for _, pTest := range proverTests {
		n = pTest.n
		p := NewProver(false)
		_, err := p.Write(pTest.commitment)
		if err != nil {
			t.Error("Error Writing Commitment: ", err)
		}

		res := make([]byte, Size)
		_, err = p.Read(res)
		if err != nil {
			t.Error("Error Reading Commitment Proof: ", err)
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
	}
}

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

type vals struct {
	v []byte
}

var leftsiblingsTests = []struct {
	n        int
	in       []byte
	expected []vals
}{
	{n: 5, in: []byte("11111"),
		expected: []vals{
			{v: []byte("11110")},
			{v: []byte("1110")},
			{v: []byte("110")},
			{v: []byte("10")},
			{v: []byte("0")},
		},
	},
}

func TestLeftSiblings(t *testing.T) {
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
	for _, s := range leftsiblingsTests {
		// Set n to known value for test
		n = s.n
		b := NewBinaryIDBytes(s.in)
		actual, err := LeftSiblings(b)
		if err != nil {
			t.Errorf("Error returned from LeftSiblings. Error: %v\n", err)
		}
		expectedBins := make([]*BinaryID, 0, len(s.expected))
		for _, vs := range s.expected {
			expectedBins = append(expectedBins, NewBinaryIDBytes(vs.v))
		}
		if !(BinaryIDListEqual(actual, expectedBins)) {
			t.Errorf(
				"Siblings Failed\nExpected:\n%v\nActual:\n%v\n",
				StringList(expectedBins),
				StringList(actual),
			)
		}
	}
}

var parentsTests = []struct {
	n        int
	in       []byte
	expected []vals
}{
	{n: 5, in: []byte("11111"), expected: []vals{
		{v: []byte("11110")},
		{v: []byte("1110")},
		{v: []byte("110")},
		{v: []byte("10")},
		{v: []byte("0")},
	}},
	{n: 9, in: []byte("00000001"), expected: []vals{
		{[]byte("000000010")},
		{[]byte("000000011")},
	}},
}

func TestGetParents(t *testing.T) {
	for _, p := range parentsTests {
		n = p.n
		b := NewBinaryIDBytes(p.in)
		actual, err := GetParents(b)
		if err != nil {
			t.Errorf("Error returned from GetParents. Error: %v\n", err)
		}
		expected := make([]*BinaryID, 0)
		for _, v := range p.expected {
			expected = append(expected, NewBinaryIDBytes(v.v))
		}
		if !(BinaryIDListEqual(actual, expected)) {
			t.Errorf(
				"GetParents Failed Expected: %v Actual: %v\n",
				StringList(expected),
				StringList(actual),
			)
		}
	}
}

func TestGetParentLargeDAG(t *testing.T) {
	debugLog.SetOutput(os.Stdout)
	defer debugLog.SetOutput(ioutil.Discard)
	n = 9
	b := NewBinaryIDBytes([]byte("00000001"))
	p, _ := GetParents(b)
	debugLog.Println(p, b)
	t.Error("Test Fail")
}

package poet

import (
	"testing"
)

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
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
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

// func TestGetParentLargeDAG(t *testing.T) {
// 	debugLog.SetOutput(os.Stdout)
// 	defer debugLog.SetOutput(ioutil.Discard)
// 	n = 9
// 	b := NewBinaryIDBytes([]byte("00000001"))
// 	p, _ := GetParents(b)
// 	debugLog.Println(p, b)
// 	t.Error("Test Fail")
// }

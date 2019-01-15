package poet

import (
	"bytes"
	"testing"
)

var newBinaryIDTests = []struct {
	in_val    uint
	in_length int
	expected  []byte
}{
	{65535, 16, []byte{byte(255), byte(255)}},
	{4095, 12, []byte{byte(15), byte(255)}},
	{511, 9, []byte{byte(1), byte(255)}},
}

func TestNewBinaryID(t *testing.T) {
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
	for _, bin := range newBinaryIDTests {
		b, err := NewBinaryID(bin.in_val, bin.in_length)
		if err != nil {
			t.Errorf("Error creating BinaryID: %v\n", err)
		}
		if b.Length != bin.in_length {
			t.Errorf("Wrong Length Error: %v\n", b.Length)
		}
		if len(b.Val) != len(bin.expected) {
			t.Errorf("Bytes Slice Wrong Length Error: %v\n", len(b.Val))
		}
		if !bytes.Equal(b.Val, bin.expected) {
			t.Errorf("Bytes Slice Wrong Value Error: %v\n", b.Val)
		}
	}
}

var newBinaryIDBytesTests = []struct {
	in       []byte
	length   int
	expected []byte
}{
	{[]byte("1111111"), 7, []byte{byte(127)}},
	{[]byte("11111111"), 8, []byte{byte(255)}},
	{[]byte("111111111"), 9, []byte{byte(1), byte(255)}},
}

func TestNewBinaryIDBytes(t *testing.T) {
	for _, bin := range newBinaryIDBytesTests {
		b := NewBinaryIDBytes(bin.in)
		if b.Length != bin.length {
			t.Errorf("Wrong Length Error: %v\n", b.Length)
		}
		if len(b.Val) != len(bin.expected) {
			t.Errorf("Bytes Slice Wrong Length Error: %v\n", len(b.Val))
		}
		if !bytes.Equal(b.Val, bin.expected) {
			t.Errorf("Bytes Slice Wrong Value Error: %v\n", b.Val)
		}
	}
}

var flipTests = []struct {
	n         int
	in_val    uint
	in_length int
	expected  []byte
}{
	{8, 255, 8, []byte{byte(254)}},
	{7, 255, 8, []byte{byte(253)}},
	{6, 255, 8, []byte{byte(251)}},
	{5, 255, 8, []byte{byte(247)}},
	{4, 255, 8, []byte{byte(239)}},
	{3, 255, 8, []byte{byte(223)}},
	{2, 255, 8, []byte{byte(191)}},
	{1, 255, 8, []byte{byte(127)}},
	{1, 15, 4, []byte{byte(7)}},
	{1, 65535, 16, []byte{byte(127), byte(255)}},
	{2, 65535, 16, []byte{byte(191), byte(255)}},
	{3, 65535, 16, []byte{byte(223), byte(255)}},
	{4, 65535, 16, []byte{byte(239), byte(255)}},
	{5, 65535, 16, []byte{byte(247), byte(255)}},
	{6, 65535, 16, []byte{byte(251), byte(255)}},
	{7, 65535, 16, []byte{byte(253), byte(255)}},
	{8, 65535, 16, []byte{byte(254), byte(255)}},
	{9, 65535, 16, []byte{byte(255), byte(127)}},
	{10, 65535, 16, []byte{byte(255), byte(191)}},
	{11, 65535, 16, []byte{byte(255), byte(223)}},
	{12, 65535, 16, []byte{byte(255), byte(239)}},
	{13, 65535, 16, []byte{byte(255), byte(247)}},
	{14, 65535, 16, []byte{byte(255), byte(251)}},
	{15, 65535, 16, []byte{byte(255), byte(253)}},
	{16, 65535, 16, []byte{byte(255), byte(254)}},
	{1, 4095, 12, []byte{byte(7), byte(255)}},
	{2, 4095, 12, []byte{byte(11), byte(255)}},
	{3, 4095, 12, []byte{byte(13), byte(255)}},
	{4, 4095, 12, []byte{byte(14), byte(255)}},
	{5, 4095, 12, []byte{byte(15), byte(127)}},
	{6, 4095, 12, []byte{byte(15), byte(191)}},
	{7, 4095, 12, []byte{byte(15), byte(223)}},
	{8, 4095, 12, []byte{byte(15), byte(239)}},
	{9, 4095, 12, []byte{byte(15), byte(247)}},
	{10, 4095, 12, []byte{byte(15), byte(251)}},
	{11, 4095, 12, []byte{byte(15), byte(253)}},
	{12, 4095, 12, []byte{byte(15), byte(254)}},
}

func TestFlipBit(t *testing.T) {
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
	for _, f := range flipTests {
		b, _ := NewBinaryID(f.in_val, f.in_length)
		b.FlipBit(f.n)
		if !bytes.Equal(b.Val, f.expected) {
			t.Errorf("Flip Bit Function Error: expected: %v, actual %v\n", f.expected, b.Val)
		}
	}
}

var truncateTests = []struct {
	in_val          uint
	in_length       int
	expected_val    uint
	expected_length int
}{
	{in_val: 15, in_length: 4, expected_val: 7, expected_length: 3},
	{in_val: 61431, in_length: 16, expected_val: 30715, expected_length: 15},
}

func TestTruncateLastBit(t *testing.T) {
	for _, tt := range truncateTests {
		b, _ := NewBinaryID(tt.in_val, tt.in_length)
		b_expected, _ := NewBinaryID(tt.expected_val, tt.expected_length)
		b.TruncateLastBit()
		if !(b.Equal(b_expected)) {
			t.Errorf(
				"Truncate Last Bit Not Correct\nExpected: %v\nActual: %v",
				b_expected,
				b,
			)
		}
	}
}

var getBitTests = []struct {
	n        int
	length   int
	val      uint
	expected int
}{
	{n: 1, length: 4, val: 7, expected: 0},
	{n: 2, length: 4, val: 7, expected: 1},
	{n: 3, length: 4, val: 7, expected: 1},
	{n: 4, length: 4, val: 7, expected: 1},
	{n: 1, length: 12, val: 3831, expected: 1},
	{n: 2, length: 12, val: 3831, expected: 1},
	{n: 3, length: 12, val: 3831, expected: 1},
	{n: 4, length: 12, val: 3831, expected: 0},
	{n: 5, length: 12, val: 3831, expected: 1},
	{n: 6, length: 12, val: 3831, expected: 1},
	{n: 7, length: 12, val: 3831, expected: 1},
	{n: 8, length: 12, val: 3831, expected: 1},
	{n: 9, length: 12, val: 3831, expected: 0},
	{n: 10, length: 12, val: 3831, expected: 1},
	{n: 11, length: 12, val: 3831, expected: 1},
	{n: 12, length: 12, val: 3831, expected: 1},
	{n: 1, length: 16, val: 61431, expected: 1},
	{n: 2, length: 16, val: 61431, expected: 1},
	{n: 3, length: 16, val: 61431, expected: 1},
	{n: 4, length: 16, val: 61431, expected: 0},
	{n: 5, length: 16, val: 61431, expected: 1},
	{n: 6, length: 16, val: 61431, expected: 1},
	{n: 7, length: 16, val: 61431, expected: 1},
	{n: 8, length: 16, val: 61431, expected: 1},
	{n: 9, length: 16, val: 61431, expected: 1},
	{n: 10, length: 16, val: 61431, expected: 1},
	{n: 11, length: 16, val: 61431, expected: 1},
	{n: 12, length: 16, val: 61431, expected: 1},
	{n: 13, length: 16, val: 61431, expected: 0},
	{n: 14, length: 16, val: 61431, expected: 1},
	{n: 15, length: 16, val: 61431, expected: 1},
	{n: 16, length: 16, val: 61431, expected: 1},
}

func TestGetBit(t *testing.T) {
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
	for _, g := range getBitTests {
		b, _ := NewBinaryID(g.val, g.length)
		i, err := b.GetBit(g.n)
		if err != nil {
			t.Errorf("Error getting bit: %v\n", err)
		}
		if i != g.expected {
			t.Errorf(
				"Bit returned not expected value Expected: %v Actual: %v n: %v\n",
				g.expected,
				i,
				g.n,
			)
		}
	}
}

var addBitTests = []struct {
	in       []byte
	addBit   int
	expected []byte
}{
	{[]byte(""), 1, []byte("1")},
	{[]byte("111"), 0, []byte("1110")},
	{[]byte("11111111"), 0, []byte("111111110")},
	{[]byte("111111110"), 1, []byte("1111111101")},
	{[]byte("00000001"), 0, []byte("000000010")},
	{[]byte("0000000100000000"), 0, []byte("00000001000000000")},
}

func TestAddBit(t *testing.T) {
	// debugLog.SetOutput(os.Stdout)
	// defer debugLog.SetOutput(ioutil.Discard)
	for i, a := range addBitTests {
		debugLog.Printf("\nTest %v\n", i)
		b := NewBinaryIDBytes(a.in)
		err := b.AddBit(a.addBit)
		if err != nil {
			t.Errorf("Error Adding Bit in Check for 1 or 0: %v\n", err)
		}
		expected := NewBinaryIDBytes(a.expected)
		if !b.Equal(expected) {
			t.Errorf(
				"Error Adding Bit Case: %v expected: %v\n",
				b,
				expected,
			)
		}
	}
}

func TestEncode(t *testing.T) {
	b1, _ := NewBinaryID(255, 8)
	b := []byte("11111111")
	if !bytes.Equal(b, b1.Encode()) {
		t.Errorf("Error encoding BinaryID bytes as utf8: %v\n", b1.Encode())
	}
}

var indexTests = []struct {
	n        int
	length   int
	val      uint
	expected int
}{
	{n: 4, length: 4, val: 0, expected: 0},
	{n: 4, length: 4, val: 1, expected: 1},
}

func TestIndex(t *testing.T) {
	for _, i := range indexTests {
		n = i.n
		b, _ := NewBinaryID(i.val, i.length)
		v := Index(b, i.n)
		if v != i.expected {
			t.Errorf(
				"Bit returned not expected value\nExpected: %v\nActual: %v\n",
				i.expected,
				v,
			)
		}
	}
}

package poet

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNewBinaryID(t *testing.T) {
	b, err := NewBinaryID(65535, 16)
	if err != nil {
		t.Errorf("Error creating BinaryID: %v\n", err)
	}
	if b.Length != 16 {
		t.Errorf("Wrong Length Error: %v\n", b.Length)
	}
	if len(b.Val) != 2 {
		t.Errorf("Bytes Slice Wrong Length Error: %v\n", len(b.Val))
	}
	if (b.Val[0] != byte(255)) || (b.Val[1] != byte(255)) {
		t.Errorf("Bytes Slice Wrong Value Error: %v, %v\n", b.Val[0], b.Val[1])
	}
}

func TestNewBinaryIDBytes(t *testing.T) {
	b, _ := NewBinaryID(255, 8)
	v := b.Encode()
	b2 := NewBinaryIDBytes(v)
	if b2.Length != 8 {
		t.Errorf("Wrong Length Error: %v\n", b2.Length)
	}
	if len(b2.Val) != 1 {
		t.Errorf("Bytes Slice Wrong Length Error: %v\n", len(b2.Val))
	}
	if b2.Val[0] != byte(255) {
		t.Errorf("Bytes Slice Wrong Value Error: %v\n", b2.Val[0])
	}
}

var flipTests = []struct {
	n         int
	in_val    uint
	in_length int
	expected  byte
}{
	{8, 255, 8, byte(254)},
	{7, 255, 8, byte(253)},
	{6, 255, 8, byte(251)},
	{5, 255, 8, byte(247)},
	{4, 255, 8, byte(239)},
	{3, 255, 8, byte(223)},
	{2, 255, 8, byte(191)},
	{1, 255, 8, byte(127)},
	{1, 15, 4, byte(7)},
}

func TestFlipBit(t *testing.T) {
	for _, f := range flipTests {
		b, _ := NewBinaryID(f.in_val, f.in_length)
		b.FlipBit(f.n)
		if b.Val[0] != f.expected {
			t.Errorf("Flip Bit Function Error: expected: %v, actual %v\n", f.expected, b.Val[0])
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
	for _, g := range getBitTests {
		b, _ := NewBinaryID(g.val, g.length)
		i, err := b.GetBit(g.n)
		if err != nil {
			t.Errorf("Error getting bit: %v\n", err)
		}
		if i != g.expected {
			t.Errorf(
				"Bit returned not expected value\nExpected: %v\nActual: %v\nn: %v\n",
				g.expected,
				i,
				g.n,
			)
		}
	}
}

func TestAddBit(t *testing.T) {
	// Case 2: 255 => 1,254
	b, _ := NewBinaryID(255, 8)
	err := b.AddBit(0)
	if err != nil {
		t.Errorf("Error Adding Bit in Check for 1 or 0: %v\n", err)
	}
	if (b.Val[0] != byte(1)) || (b.Val[1] != byte(254)) {
		t.Errorf("Error Adding Bit Case 1: %v\n", b)
	}
	// Case 2: 127 => 255
	b2, _ := NewBinaryID(127, 7)
	err = b2.AddBit(1)
	if err != nil {
		t.Errorf("Error Adding Bit in Check for 1 or 0: %v\n", err)
	}
	if b2.Val[0] != byte(255) {
		t.Errorf("Error Adding Bit Case 2: %v\n", b2)
	}
}

func TestEncode(t *testing.T) {
	b1, _ := NewBinaryID(255, 8)
	b := []byte{'1', '1', '1', '1', '1', '1', '1', '1'}
	if !bytes.Equal(b, b1.Encode()) {
		t.Errorf("Error encoding BinaryID bytes as utf8: %v\n", b1.Encode())
	}
}

func TestBitLength(t *testing.T) {
	// testing microphone
	b, _ := NewBinaryID(255, 8)
	b.BitList()
	b.Encode()

	test := make([]byte, 10)
	binary.PutUvarint(test, uint64(1000))
	//fmt.Println(test)

	//value, _ := binary.Uvarint(test)
	//fmt.Println(value)

	t.Error("failed")
}

var indexTests = []struct {
	length   int
	val      uint
	expected int
}{
	{length: 4, val: 0, expected: 1},
	{length: 4, val: 1, expected: 2},
}

func TestIndex(t *testing.T) {
	for _, i := range indexTests {
		b, _ := NewBinaryID(i.val, i.length)
		v := Index(b)
		if v != i.expected {
			t.Errorf(
				"Bit returned not expected value\nExpected: %v\nActual: %v\n",
				i.expected,
				v,
			)
		}
	}
}

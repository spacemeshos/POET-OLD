package poet

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	n        int
	expected byte
}{
	{1, byte(254)},
	{2, byte(253)},
	{3, byte(251)},
	{4, byte(247)},
	{5, byte(239)},
	{6, byte(223)},
	{7, byte(191)},
	{8, byte(127)},
}

func TestFlipBit(t *testing.T) {
	for _, f := range flipTests {
		b, _ := NewBinaryID(255, 8)
		b.FlipBit(f.n)
		if b.Val[0] != f.expected {
			t.Errorf("Flip Bit Function Error: expected: %v, actual %v\n", f.expected, b.Val[0])
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
}

func TestGetBit(t *testing.T) {
	for _, g := range getBitTests {
		b, _ := NewBinaryID(g.val, g.length)
		fmt.Println(b, g.n)
		i, err := b.GetBit(g.n)
		if err != nil {
			t.Errorf("Error getting bit: %v\n", err)
		}
		if i != g.expected {
			t.Errorf(
				"Bit returned not expected value\nExpected: %v\nActual: %v\n",
				g.expected,
				i,
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
	fmt.Println(test)

	value, _ := binary.Uvarint(test)
	fmt.Println(value)

	t.Error("failed")
}

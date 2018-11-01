package poet

import (
	"bytes"
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

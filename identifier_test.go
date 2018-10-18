package poet

import (
	"testing"
)

func TestNewBinaryID(t *testing.T) {
	b, err := NewBinaryID(65535, 16)
	if err != nil {
		t.Errorf("Error creating BinaryID: %v\n", err)
	}
	if b.length != 16 {
		t.Errorf("Wrong Length Error: %v\n", b.length)
	}
	if len(b.val) != 2 {
		t.Errorf("Bytes Slice Wrong Length Error: %v\n", len(b.val))
	}
	if (b.val[0] != byte(255)) || (b.val[1] != byte(255)) {
		t.Errorf("Bytes Slice Wrong Length Error: %v, %v\n", b.val[0], b.val[1])
	}
}

func TestFlipBit(t *testing.T) {
	b, _ := NewBinaryID(255, 8)
	b.FlipBit(4)
	if b.val[0] != byte(239) {
		t.Errorf("Flip Bit Function Error: %v\n", b.val[0])
	}
}

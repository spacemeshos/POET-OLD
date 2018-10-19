package poet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/bits"
)

// BinaryID is a binary representation of the ID of a node. The length is
// needed because we need to identify between 0 and 00 for example.
// Note: added json Marshalling. Currently, encoding is base64 (for []byte).
// This is default for []byte in Go, but can change as needed.
type BinaryID struct {
	Val    []byte `json:"Values"`
	Length int    `json:"Length"`
}

func NewBinaryID(val uint, length int) (*BinaryID, error) {
	if bits.Len(val) > length {
		return nil, errors.New("Length not long enough")
	}
	idx := length / 8
	if (length % 8) != 0 {
		idx = idx + 1
	}
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(v, uint64(val))
	b := new(BinaryID)
	b.Val = make([]byte, idx)
	for i := 0; i < idx; i++ {
		b.Val[idx-i-1] = v[7-i]
	}
	b.Length = length
	return b, nil
}

func NewBinaryIDInt(val uint) *BinaryID {
	b := new(BinaryID)
	l := bits.Len(val) / 8
	b.Length = l
	b.Val = make([]byte, l)
	binary.PutUvarint(b.Val, uint64(val))
	return b
}

func (b *BinaryID) Equal(b2 *BinaryID) bool {
	return (b.Length == b2.Length) && bytes.Equal(b.Val, b2.Val)
}

func (b *BinaryID) GreaterThan(b2 *BinaryID) bool {
	if b.Length > b2.Length {
		return true
	} else if b.Length < b2.Length {
		return false
	}
	// TODO: Check number of bytes read
	bn, _ := binary.Uvarint(b.Val)
	b2n, _ := binary.Uvarint(b2.Val)
	return bn > b2n
}

// Flip the n'th bit from 0 to 1 or 1 to 0. Does nothing if n > length
func (b *BinaryID) FlipBit(n int) {
	if n >= b.Length {
		return
	}
	shift := uint(n % 8)
	idx := n / 8
	if (b.Val[idx] * (1 << shift)) == 0 {
		b.Val[idx] = b.Val[idx] + 1<<shift
	} else {
		b.Val[idx] = b.Val[idx] - 1<<shift
	}
}

func (b *BinaryID) TruncateLastBit() {
	carry := 0
	for i := 0; i < b.Length; i++ {
		add := carry * 1 << 8
		carry = int(b.Val[i] * 1)
		b.Val[i] = b.Val[i] >> 1
		b.Val[i] = b.Val[i] + byte(add)
	}
}

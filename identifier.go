package poet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/bits"
)

// BinaryID is a binary representation of the ID of a node. The length is
// needed because we need to identify between 0 and 00 for example.
type BinaryID struct {
	val    []byte
	length int
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
	b.val = make([]byte, idx)
	for i := 0; i < idx; i++ {
		b.val[idx-i-1] = v[7-i]
	}
	b.length = length
	return b, nil
}

func NewBinaryIDInt(val uint) *BinaryID {
	b := new(BinaryID)
	l := bits.Len(val) / 8
	b.length = l
	b.val = make([]byte, l)
	binary.PutUvarint(b.val, uint64(val))
	return b
}

func (b *BinaryID) Equal(b2 *BinaryID) bool {
	return (b.length == b2.length) && bytes.Equal(b.val, b2.val)
}

func (b *BinaryID) GreaterThan(b2 *BinaryID) bool {
	if b.length > b2.length {
		return true
	} else if b.length < b2.length {
		return false
	}
	// TODO: Check number of bytes read
	bn, _ := binary.Uvarint(b.val)
	b2n, _ := binary.Uvarint(b2.val)
	return bn > b2n
}

// Flip the n'th bit from 0 to 1 or 1 to 0. Does nothing if n > length
func (b *BinaryID) FlipBit(n int) {
	if n >= b.length {
		return
	}
	shift := uint(n % 8)
	idx := n / 8
	if (b.val[idx] * (1 << shift)) == 0 {
		b.val[idx] = b.val[idx] + 1<<shift
	} else {
		b.val[idx] = b.val[idx] - 1<<shift
	}
}

func (b *BinaryID) TruncateLastBit() {
	carry := 0
	for i := 0; i < b.length; i++ {
		add := carry * 1 << 8
		carry = int(b.val[i] * 1)
		b.val[i] = b.val[i] >> 1
		b.val[i] = b.val[i] + byte(add)
	}
}

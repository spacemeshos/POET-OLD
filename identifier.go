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
	b := new(BinaryID)
	b.val = make([]byte, 32) // TODO: Need to calc size
	binary.PutUvarint(b.val, uint64(val))
	return b, nil
}

func NewBinaryIDInt(val uint) *BinaryID {
	return nil
}

func (b *BinaryID) Equal(b2 *BinaryID) bool {
	return (b.length == b2.length) && bytes.Equal(b.val, b2.val)
}

func (b *BinaryID) GreaterThan(b2 *BinaryID) bool {
	return false
}

package poet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math/bits"
	"strconv"
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
	b.Val = v[(8 - idx):8]
	b.Length = length
	return b, nil
}

func NewBinaryIDInt(val uint) *BinaryID {
	var l int
	if val == 0 {
		l = 1
	} else {
		l = bits.Len(val) / 8
	}
	b, err := NewBinaryID(val, l)
	if err != nil {
		log.Panic("Error Creating BinaryID from single uint")
	}
	return b
}

func NewBinaryIDBytes(v []byte) *BinaryID {
	b := new(BinaryID)
	b.Length = len(v)
	debugLog.Printf(
		"Creating BinaryID: %v\n",
		string(v),
	)
	l := b.Length / 8
	if (b.Length % 8) != 0 {
		l = l + 1
	}
	b.Val = make([]byte, l)
	for i := 0; i < b.Length; i++ {
		j := b.Length - i
		byteBit := v[j-1]
		if byteBit == byte('1') {
			b.FlipBit(j)
		}
	}
	return b
}

func NewBinaryIDCopy(b *BinaryID) (b2 *BinaryID) {
	b2 = new(BinaryID)
	b2.Length = b.Length
	b2.Val = make([]byte, len(b.Val))
	copy(b2.Val, b.Val)
	return b2
}

func BinaryIDListEqual(b1 []*BinaryID, b2 []*BinaryID) bool {
	if (b1 == nil) != (b2 == nil) {
		return false
	}
	if len(b1) != len(b2) {
		return false
	}
	for i := range b1 {
		if !(b1[i].Equal(b2[i])) {
			return false
		}
	}
	return true
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
	if n > b.Length {
		return
	}
	add := len(b.Val)*8 - b.Length
	shift := uint((b.Length - n) % 8)
	idx := (n + add - 1) / 8
	if (b.Val[idx] & (1 << shift)) == 0 {
		b.Val[idx] = b.Val[idx] + 1<<shift
	} else {
		b.Val[idx] = b.Val[idx] - 1<<shift
	}
}

func (b *BinaryID) TruncateLastBit() {
	carry := 0
	for i := 0; i < len(b.Val); i++ {
		add := carry & 1 << 7
		carry = int(b.Val[i] & 1)
		b.Val[i] = b.Val[i] >> 1
		b.Val[i] = b.Val[i] + byte(add)
	}
	b.Length = b.Length - 1
}

// Pretty printing function for debugging. Not for encoding.
func (b *BinaryID) String() string {
	return string(b.Encode())
}

// Pretty printing function for a list of BinaryID's. Primarily for Debugging
func StringList(bList []*BinaryID) string {
	var buf bytes.Buffer
	for _, b := range bList {
		buf.WriteString(b.String() + " ")
	}
	return buf.String()
}

// Returns if n'th bit is 0 or 1. Error if n > Length
// n'th bit from the left. So for 1011 the 1'st bit is the first 1 on the left
func (b *BinaryID) GetBit(n int) (int, error) {
	if (n > b.Length) || (n == 0) {
		return 0, errors.New("n wrong length for binaryID")
	}
	add := len(b.Val)*8 - b.Length
	shift := uint((b.Length - n) % 8)
	idx := (n + add - 1) / 8
	if (b.Val[idx] & (1 << shift)) == 0 {
		return 0, nil
	} else {
		return 1, nil
	}
}

// Adds 0 or 1 to lsb of BinaryID. Returns error if not 0 or 1
func (b *BinaryID) AddBit(n int) error {
	isZero := n == 0
	isOne := n == 1
	if !isZero && !isOne {
		return errors.New("Not 0 or 1. Cannot add bit.")
	}
	debugLog.Printf("Adding %v to %v", n, b)
	debugLog.Printf("Length: %v Val: %v", b.Length, b.Val)
	if b.Length == 0 {
		b.Length = 1
		b.Val = make([]byte, 1)
		b.Val[0] = byte(n)
		return nil
	}
	l := len(b.Val)
	if b.Length%8 == 0 {
		for i := 0; i < l+1; i++ {
			if i == 0 {
				b.Val = append(b.Val, b.Val[l-1]<<1+byte(n))
			} else if i == l {
				b.Val[l-i] = b.Val[l-i] & (1 << 7)
				b.Val[l-i] = b.Val[l-i] >> 7
			} else {
				carry := b.Val[l-i] & (1 << 7)
				carry = carry >> 7
				b.Val[l-i] = b.Val[l-i-1]<<1 + carry
			}
		}
	} else {
		carry := byte(n)
		for i := 0; i < l; i++ {
			idx := l - i - 1
			add := carry
			carry = b.Val[idx] & (1 << 7)
			carry = carry >> 7
			b.Val[idx] = b.Val[idx] << 1
			b.Val[idx] = b.Val[idx] + byte(add)
		}
	}
	b.Length = b.Length + 1
	return nil
}

// Encode outputs a []byte encoded in utf8
func (b *BinaryID) Encode() (v []byte) {
	v = make([]byte, 0, b.Length)
	// if b.Length == 0 {
	// 	v = append(v, []byte("\xe2")...)
	// }
	for i := 1; i <= b.Length; i++ {
		bit, err := b.GetBit(i)
		if err != nil {
			break
		}
		s := strconv.Itoa(bit)
		v = append(v, []byte(s)...)
	}
	return v
}

func TreeSize(b *BinaryID) (size int) {
	size = 1<<uint(n-b.Length+1) - 1
	return size
}

func Index(b *BinaryID) (index int) {
	index = TreeSize(b) - 1
	// TODO: Add error check
	ls, _ := Siblings(b, true)
	for _, l := range ls {
		index = index + TreeSize(l)
	}
	return index
}

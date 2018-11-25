package poet

import (
	// "strings"
	"encoding/binary"
	"testing"
)

func TestFileIO(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	file := NewFileIO()

	binID := NewBinaryIDBytes([]byte("hello"))

	scParam := make([]byte, binary.MaxVarintLen64)
	binary.BigEndian.PutUint64(scParam, uint64(10))

	err := file.StoreLabel(binID, scParam)
	if err != nil {
		t.Error("Failed to store data", err)

	}

	data, err := file.GetLabel(binID)
	if err != nil {
		t.Error("Failed to retrieve data", err)
	}
	t.Log(data)
}

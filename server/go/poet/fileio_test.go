package poet

import (
	// "strings"
	"testing"
	"encoding/binary"
)


func TestFileIO(t *testing.T){
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
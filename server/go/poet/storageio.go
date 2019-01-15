package poet

type StorageIO interface {
	StoreLabel(b *BinaryID, label []byte) error
	GetLabel(*BinaryID) (label []byte, err error)
	LabelCalculated(*BinaryID) (bool, error)
}

// Some storageIO implementations need to know about the DAG Size (n)
// Calls to change DAG size will check if the storageIO can change DAG size
// then will forward call this interface function.
type SetDAGSizer interface {
	SetDAGSize(size int)
}

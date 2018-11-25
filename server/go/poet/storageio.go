package poet

type StorageIO interface {
	StoreLabel(b *BinaryID, label []byte) error
	GetLabel(*BinaryID) (label []byte, err error)
	LabelCalculated(*BinaryID) (bool, error)
}

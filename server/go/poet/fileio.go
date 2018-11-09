package poet

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type FileIO struct {
	file          *os.File
	storeLabel    chan []byte
	storeError    chan error
	getLabel      chan *BinaryID
	retLabel      chan *retLabel
	labelComputed chan *BinaryID
	retComputed   chan *retComputed
}

type retLabel struct {
	label []byte
	err   error
}

type retComputed struct {
	computed bool
	err      error
}

func NewFileIO() (f *FileIO) {
	f = new(FileIO)
	wd, err := os.Getwd()
	if err != nil {
		log.Panic("Error Getting Working Directory", err)
	}
	f.file, err = ioutil.TempFile(wd, "labels")
	if err != nil {
		log.Panic("Error Creating File: ", err)
	}
	//f.file = file
	// Create all the channels needed
	// This is likely not optimal. TODO(later): find a better communication pattern
	f.storeLabel = make(chan []byte)
	f.storeError = make(chan error)
	f.getLabel = make(chan *BinaryID)
	f.retLabel = make(chan *retLabel)
	f.labelComputed = make(chan *BinaryID)
	f.retComputed = make(chan *retComputed)

	// Start f.run() goroutine which handles file io
	go f.run()
	return f
}

func (f *FileIO) run() {
	defer f.file.Close()
	for {
		// TODO: We migth need a close channel here. Otherwise the file will just
		// stay open even when the program exits
		select {
		case b := <-f.storeLabel:
			_, err := f.file.Seek(0, 2)
			if err != nil {
				f.storeError <- err
				break
			}
			_, err = f.file.Write(b)
			f.storeError <- err
		case b := <-f.getLabel:
			ret := new(retLabel)
			idx := int64(Index(b)) * int64(size)
			_, err := f.file.Seek(0, 0)
			if err != nil {
				ret.err = err
				f.retLabel <- ret
				break
			}
			ret.label = make([]byte, size)
			f.file.ReadAt(ret.label, idx)
			fmt.Println(
				"Fetched node ",
				string(b.Encode()),
				" hash: ",
				ret.label,
			)
			f.retLabel <- ret
		case b := <-f.labelComputed:
			ret := new(retComputed)
			stats, err := f.file.Stat()
			if err != nil {
				ret.err = err
				f.retComputed <- ret
				break
			}
			idx := int64(Index(b)+1) * int64(size)
			s := stats.Size()
			//fmt.Println("Node: ", string(b.Encode()))
			fmt.Println("Index: ", Index(b), "filesize", s)
			if idx <= s {
				ret.computed = true
			} else {
				ret.computed = false
			}
			f.retComputed <- ret
		}
	}

	
}

func (f *FileIO) StoreLabel(b *BinaryID, label []byte) error {
	f.storeLabel <- label
	err := <-f.storeError
	return err
}

func (f *FileIO) GetLabel(b *BinaryID) (label []byte, err error) {
	f.getLabel <- b
	ret := <-f.retLabel
	return ret.label, ret.err
}
func (f *FileIO) LabelCalculated(b *BinaryID) (bool, error) {
	f.labelComputed <- b
	ret := <-f.retComputed
	return ret.computed, ret.err
}

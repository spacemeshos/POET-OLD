package verifier

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	t        = int(150)
	w        = int(256)
	size     = int(32)
	debugLog = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime)
	infoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logFile  *os.File
)

func init() {

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Panic("Error Getting Working Directory", err)
	// }
	// logFile, err = ioutil.TempFile(wd, "log")
	// if err != nil {
	// 	log.Panic("Error Creating File: ", err)
	// }
}

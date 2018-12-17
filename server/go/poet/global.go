package poet

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	n        = *flag.Int("DAG_Size", 4, "Set DAG Size") // Low for testing. TODO: Set correct default
	m        = *flag.Int("DAG_Size_Store", 4, "Set DAG Size to Store")
	t        = *flag.Int("Security_Param_t", 150, "Set the Security Parameter t")
	w        = *flag.Int("Security_Param_w", 256, "Set the Security Parameter w")
	HashSize = *flag.Int("Hash_Size_Bytes", 32, "Set the hash size") // TODO: Set Dynamically based on hash library
	filepath = *flag.String("filepath", "./test.dat", "Set the prover file storage path")
	debugLog = log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime)
	infoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logFile  *os.File
)

func init() {
	flag.Parse()

	// wd, err := os.Getwd()
	// if err != nil {
	// 	log.Panic("Error Getting Working Directory", err)
	// }
	// logFile, err = ioutil.TempFile(wd, "log")
	// if err != nil {
	// 	log.Panic("Error Creating File: ", err)
	// }
}

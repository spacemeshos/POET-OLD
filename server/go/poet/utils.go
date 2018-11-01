package poet

import (
	"os"
	"bufio"
	"fmt"
)
// BitsToInt 
func BitsToInt(data []byte) int {
	return 0
}

func WriteToFile(data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	// write to file
	fmt.Fprintln(w, data)
	return w.Flush()
}

func ReadLabelFromFile(offset int) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	var data []byte
	for scanner.Scan() {
		if i != offset {
			i++
			continue
		}
		data = scanner.Bytes()
		break
	}
	return data, nil
}

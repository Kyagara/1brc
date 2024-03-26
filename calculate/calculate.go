package calculate

import (
	"bufio"
	"os"
)

func Run(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	read := 0
	for scanner.Scan() {
		read++
	}
	println(read)
	return nil
}

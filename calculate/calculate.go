package calculate

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"slices"
	"strconv"
	"unsafe"
)

type Station struct {
	Name string
	Min  float32
	Mean float32
	Max  float32
}

func Run(path string) ([]Station, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// 1gb
	max := 1_024_000_000
	fileSize := int(stat.Size())
	alloc := 0
	if max >= fileSize {
		alloc = fileSize
	} else {
		alloc = max
	}

	scanner := bufio.NewReaderSize(file, alloc)

	delimiter := byte('\n')
	separator := byte(';')

	stations := make(map[string][]float32, 10000)

	for {
		line, err := scanner.ReadSlice(delimiter)
		if err == io.EOF {
			break
		}

		i := bytes.IndexByte(line, separator)

		name := line[:i]
		station := *(*string)(unsafe.Pointer(&name))
		tempBytes := line[i+1 : len(line)-2]
		tempString := *(*string)(unsafe.Pointer(&tempBytes))
		temperature, err := strconv.ParseFloat(tempString, 32)
		if err != nil {
			return nil, err
		}

		if stations[station] == nil {
			stations[station] = make([]float32, 100)
		}

		stations[station] = append(stations[station], float32(temperature))
	}

	sortedNames := make([]string, 0, 10000)
	for key := range stations {
		sortedNames = append(sortedNames, key)
	}
	slices.Sort(sortedNames)

	calculated := make([]Station, 0, 10000)

	for _, name := range sortedNames {
		temps := stations[name]
		slices.Sort(temps)
		station := Station{
			Name: name,
			Min:  temps[0],
			Mean: 0,
			Max:  temps[len(temps)-1],
		}

		calculated = append(calculated, station)
	}

	return calculated, nil
}

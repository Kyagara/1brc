package calculate

import (
	"bufio"
	"bytes"
	"hash/fnv"
	"io"
	"os"
	"sort"
	"unsafe"
)

// Final result
type Station struct {
	Name string
	Min  float32
	Mean float32
	Max  float32
}

type Info struct {
	Name  string
	Temps []float32
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

	// Reading lines

	scanner := bufio.NewReaderSize(file, alloc)

	delimiter := byte('\n')
	separator := byte(';')

	stations := make(map[uint32]Info, 10000)

	for {
		line, err := scanner.ReadSlice(delimiter)
		if err == io.EOF {
			break
		}

		i := bytes.IndexByte(line, separator)

		nameBytes := line[:i]
		name := *(*string)(unsafe.Pointer(&nameBytes))
		hash := hash(nameBytes)

		tempBytes := line[i+1 : len(line)-2]
		temperature := parseFloat32(tempBytes)

		info, ok := stations[hash]
		if !ok {
			// Attempting to prealloc some floats based on the filesize, just guess idk
			max := 1024
			// These numbers were revealed to me in a dream
			revelation := fileSize / 1024 / 128
			alloc := 0
			if max >= revelation {
				alloc = revelation
			} else {
				alloc = max
			}

			info = Info{
				Name:  name,
				Temps: make([]float32, 0, alloc),
			}
		}

		info.Temps = append(info.Temps, temperature)
		stations[hash] = info
	}

	// Sorting and calculating

	var sortedKeys []uint32
	for name := range stations {
		sortedKeys = append(sortedKeys, name)
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		return stations[sortedKeys[i]].Name < stations[sortedKeys[j]].Name
	})

	calculated := make([]Station, 0, len(sortedKeys))

	for _, name := range sortedKeys {
		info := stations[name]
		count := float32(len(info.Temps))

		var sum float32
		var min float32
		var max float32
		for _, temp := range info.Temps {
			sum += temp

			if temp < min {
				min = temp
			}
			if temp > max {
				max = temp
			}
		}

		mean := sum / count

		station := Station{
			Name: info.Name,
			Min:  min,
			Mean: mean,
			Max:  max,
		}

		calculated = append(calculated, station)
	}

	return calculated, nil
}

func parseFloat32(bytes []byte) float32 {
	var result float32
	var power float32 = 1
	isNegative := false
	decimal := false

	for i, b := range bytes {
		if i == 0 && b == '-' {
			isNegative = true
			continue
		}

		if b >= '0' && b <= '9' {
			if decimal {
				power *= 10
			}

			result = result*10 + float32(b-'0')
		} else if b == '.' {
			decimal = true
		}
	}

	if isNegative {
		result *= -1
	}

	result /= power
	return result
}

func hash(s []byte) uint32 {
	h := fnv.New32a()
	h.Write(s)
	return h.Sum32()
}

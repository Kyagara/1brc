package calculate

import (
	"bufio"
	"io"
	"os"
	"sort"
)

const (
	bufferSize = 65536
	delimiter  = byte('\n')
	separator  = byte(';')
)

// Final result
type Station struct {
	Name [100]byte
	Min  float32
	Mean float32
	Max  float32
}

type Info struct {
	Count int
	Total float32
	Min   float32
	Max   float32
	Name  [100]byte
}

func Run(path string) ([]Station, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 64kb
	reader := bufio.NewReaderSize(file, bufferSize)

	stations := make(map[uint32]Info, 10000)

	// Reading lines
	for {
		line, separatorIndex, err := read(reader)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if separatorIndex == -1 {
			break
		}

		nameBytes := line[:separatorIndex]
		hash := hash(nameBytes)

		tempBytes := line[separatorIndex+1:]
		temp := parseFloat32(tempBytes)

		info, ok := stations[hash]
		if !ok {
			var name [100]byte
			copy(name[:], nameBytes)

			info = Info{
				Name:  name,
				Count: 0,
				Total: 0,
				Min:   0,
				Max:   0,
			}
		}

		info.Count++
		info.Total += temp
		if temp > info.Max {
			info.Max = temp
		}
		if temp < info.Min {
			info.Min = temp
		}

		stations[hash] = info
	}

	// Sorting and calculating

	sortedKeys := make([]uint32, 0, len(stations))
	for name := range stations {
		sortedKeys = append(sortedKeys, name)
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		name1 := stations[sortedKeys[i]].Name
		len1 := len(name1)
		name2 := stations[sortedKeys[j]].Name
		len2 := len(name2)

		for i, j := 0, 0; i < len1 && j < len2; i, j = i+1, j+1 {
			diff := int32(name1[i]) - int32(name2[j])
			if diff != 0 {
				return diff < 0
			}
		}
		return len1 < len2
	})

	calculated := make([]Station, 0, len(sortedKeys))

	for _, hash := range sortedKeys {
		info := stations[hash]
		mean := info.Total / float32(info.Count)

		station := Station{
			Name: info.Name,
			Min:  info.Min,
			Mean: mean,
			Max:  info.Max,
		}

		calculated = append(calculated, station)
	}

	return calculated, nil
}

func read(reader *bufio.Reader) ([]byte, int, error) {
	line, err := reader.ReadSlice(delimiter)
	if err != nil && err != bufio.ErrBufferFull {
		return nil, -1, err
	}

	separatorIndex := -1
	delimiterIndex := len(line) - 1

	// looping here is slow
	for i, b := range line {
		if b == separator {
			separatorIndex = i
			break
		}
	}

	return line[:delimiterIndex], separatorIndex, nil
}

func parseFloat32(data []byte) float32 {
	var result float32
	var power float32 = 10
	negative := false

	if data[0] == '-' {
		negative = true
	}

	for _, b := range data {
		if b >= '0' && b <= '9' {
			result = result*10 + float32(b-'0')
		}
	}

	if negative {
		result *= -1
	}

	result /= power
	return result
}

// Should be enough for a hash function
func hash(bytes []byte) uint32 {
	var result uint32
	for _, b := range bytes {
		result = result*31 + uint32(b)
	}
	return result
}

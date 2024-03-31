package calculate

import (
	"bytes"
	"io"
	"runtime"
	"sync"

	"golang.org/x/exp/mmap"
)

const (
	bufferSize = 65536
	delimiter  = byte('\n')
	separator  = byte(';')
	equal      = byte('=')
	slash      = byte('/')
	minus      = byte('-')
)

var (
	separators = []byte{',', ' '}
)

type Station struct {
	Hash  uint32
	Count float32
	Total float32
	Min   float32
	Max   float32
	Name  [100]byte
}

func Run(path string) ([]byte, uint32, error) {
	// Should be 10000, but the hash function has collisions
	// This number was revealed to me in a dream
	hashmap := NewHashMap(30008)

	reader, err := mmap.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	cpu := runtime.NumCPU()
	processChannel := make(chan []byte, bufferSize*cpu)

	var wg sync.WaitGroup
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go process(&wg, hashmap, processChannel)
	}

	read(reader, processChannel)
	wg.Wait()

	hashmap.Sort()

	// Writing the string output

	output := make([]byte, 0, 128*hashmap.Size)
	writer := bytes.NewBuffer(output)

	writer.WriteByte('{')
	num := make([]byte, 0, 5)

	for _, station := range hashmap.Entries {
		if station.Count == 0 {
			continue
		}

		mean := roundTowardsPositive(station.Total / station.Count)
		// Name=Min/Mean/Max,<space>
		writeStation(writer, station, mean, num)
	}

	// Removing last ,<space>
	writer.Truncate(writer.Len() - 2)
	writer.WriteByte('}')
	writer.WriteByte(delimiter)

	return writer.Bytes(), hashmap.Size, nil
}

func read(reader *mmap.ReaderAt, process chan []byte) {
	defer close(process)

	var size int64 = int64(reader.Len())
	var at int64

	buffer := make([]byte, bufferSize)
	for at < size {
		n, err := reader.ReadAt(buffer, at)
		if err != nil && err != io.EOF {
			return
		}

		if buffer[n-1] != delimiter {
			// Backtrack until delimiter is found
			for i := n - 1; i >= 0; i-- {
				if buffer[i] == delimiter {
					n = i + 1
					break
				}
			}
		}

		bufferCopy := make([]byte, n)
		copy(bufferCopy, buffer[:n])
		process <- bufferCopy

		at += int64(n)
	}
}

func process(wg *sync.WaitGroup, hashmap *HashMap, process <-chan []byte) {
	defer wg.Done()

	for buffer := range process {
		name := make([]byte, 0, 100)

		var temperature float32
		negative := false

		separatorIndex := -1
		read := 0

		for i := 0; i < len(buffer); i++ {
			b := buffer[i]

			switch b {
			// Getting the name and then looping until the delimiter for the temperature
			case separator:
				separatorIndex = i
				name = buffer[read:separatorIndex]

				// Getting the float

				numIndex := i + 1
				if buffer[numIndex] == minus {
					negative = true
					numIndex++
				}

				for ; buffer[numIndex] != delimiter; numIndex++ {
					num := buffer[numIndex]
					if num >= '0' && num <= '9' {
						temperature = temperature*10 + float32(num-'0')
					}
				}

				if negative {
					temperature *= -1
				}

				temperature /= 10

				i = numIndex - 1

			// Calculating the temperature and storing the station
			case delimiter:
				read = i + 1

				hash := hashmap.Hash(name)
				station, ok := hashmap.Get(hash)
				if !ok {
					var nameCopy [100]byte
					copy(nameCopy[:], name)

					station = Station{
						Name:  nameCopy,
						Hash:  hash,
						Count: 0,
						Total: 0,
						Min:   0,
						Max:   0,
					}
				}

				station.Count++
				station.Total += temperature
				if temperature > station.Max {
					station.Max = temperature
				}

				if temperature < station.Min {
					station.Min = temperature
				}

				hashmap.Set(hash, station)

				// Reset
				temperature = 0
				negative = false
			}
		}
	}
}

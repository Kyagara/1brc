package calculate

import (
	"bytes"
	"io"
	"runtime"
	"sync"

	"golang.org/x/exp/mmap"
)

const (
	// This size works well on my machine
	bufferSize = 1 << 21
	// Hash map has colisions, with this size my dataset of 413 stations works well
	// This number was revealed to me in a dream
	hashMapSize = 30008
)

const (
	delimiter = byte('\n')
	separator = byte(';')
	equal     = byte('=')
	slash     = byte('/')
	minus     = byte('-')
)

var (
	separators = []byte{',', ' '}
)

type Station struct {
	Name  []byte
	Count float32
	Total float32
	Min   float32
	Max   float32
	Hash  uint32
}

func Run(path string) ([]byte, int, error) {
	reader, err := mmap.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	cpu := runtime.NumCPU()
	hashmap := NewHashMap(cpu)

	processChannel := make(chan []byte, bufferSize)

	var wg sync.WaitGroup
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go process(&wg, i, hashmap, processChannel)
	}

	go read(reader, processChannel)
	wg.Wait()

	stations := hashmap.Sort()

	// Writing the string output

	output := make([]byte, 0, 128*len(stations))
	writer := bytes.NewBuffer(output)

	writer.WriteByte('{')
	num := make([]byte, 0, 5)

	for _, station := range stations {
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

	return writer.Bytes(), len(stations), nil
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

		// This copy makes me upset
		bufferCopy := make([]byte, n)
		copy(bufferCopy, buffer[:n])
		process <- bufferCopy

		at += int64(n)
	}
}

func process(wg *sync.WaitGroup, shard int, hashmap *HashMap, process <-chan []byte) {
	defer wg.Done()

	for buffer := range process {
		var temperature float32
		negative := false

		separatorIndex := -1
		read := 0

		for i := 0; i < len(buffer); i++ {
			b := buffer[i]

			switch b {
			// Getting the name and then looping for the temperature
			case separator:
				// The last character of the station name ends at i-1
				separatorIndex = i

				numIndex := i + 1
				for ; buffer[numIndex] != delimiter; numIndex++ {
					num := buffer[numIndex]
					if num == minus {
						negative = true
						continue
					}

					if num >= '0' && num <= '9' {
						temperature = temperature*10 + float32(num-'0')
					}
				}

				if negative {
					temperature *= -1
				}

				temperature /= 10

				// Adjusting the index because of the temperature loop
				i = numIndex - 1

			// Calculating the temperature and storing the station
			case delimiter:
				name := buffer[read:separatorIndex]
				read = i + 1

				hashmap.Set(shard, name, temperature)

				// Reset
				temperature = 0
				negative = false
			}
		}
	}
}

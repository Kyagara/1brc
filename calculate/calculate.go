package calculate

import (
	"bytes"
	"io"
	"os"
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

var (
	at       int64
	fileSize int64
	atMutex  = &sync.Mutex{}
)

func Run(path string) ([]byte, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	file.Close()
	fileSize = stat.Size()

	cpu := runtime.NumCPU()
	hashmap := NewHashMap(cpu)

	var wg sync.WaitGroup
	for i := 0; i < cpu; i++ {
		wg.Add(1)
		go read(&wg, i, hashmap, path)
	}

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

func read(wg *sync.WaitGroup, shard int, hashmap *HashMap, p string) {
	defer wg.Done()

	reader, err := mmap.Open(p)
	if err != nil {
		return
	}
	defer reader.Close()

	buffer := make([]byte, bufferSize)
	for {
		atMutex.Lock()
		if at >= fileSize {
			atMutex.Unlock()
			return
		}

		n, err := reader.ReadAt(buffer, at)
		if err != nil && err != io.EOF {
			atMutex.Unlock()
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

		at += int64(n)
		atMutex.Unlock()

		// This copy makes me upset
		b := make([]byte, n)
		copy(b, buffer[:n])
		process(shard, hashmap, b)
	}
}

func process(shard int, hashmap *HashMap, buffer []byte) {
	var temperature float32
	negative := false

	separatorIndex := 0
	nameIndex := 0

	for i, b := range buffer {
		switch b {
		case separator:
			separatorIndex = i

		// Calculating the temperature and storing the station
		case delimiter:
			name := buffer[nameIndex:separatorIndex]
			temp := buffer[separatorIndex+1 : i]

			for _, num := range temp {
				if num == minus {
					negative = true
					continue
				}

				if (num & 0xF0) == 0x30 {
					// Expensive
					temperature = temperature*10 + float32(num-'0')
				}
			}

			if negative {
				temperature *= -1
			}

			temperature /= 10

			// Really expensive (cache miss)
			hashmap.Set(shard, name, temperature)

			// Reset
			temperature = 0
			negative = false
			nameIndex = i + 1
		}
	}
}

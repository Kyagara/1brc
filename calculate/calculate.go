package calculate

import (
	"bytes"
	"io"
	"strconv"

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

type Info struct {
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

	err = read(reader, hashmap)
	if err != nil {
		return nil, 0, err
	}

	hashmap.Sort()

	output := make([]byte, 0, 128*hashmap.Size)
	writer := bytes.NewBuffer(output)

	writer.WriteByte('{')
	num := make([]byte, 0, 5)
	for _, info := range hashmap.Entries {
		if info.Count == 0 {
			continue
		}

		mean := roundTowardsPositive(info.Total / info.Count)

		// Name=Min/Mean/Max,<space>
		writer.Write(info.Name[:])
		writer.WriteByte(equal)
		num = strconv.AppendFloat(num, float64(info.Min), 'f', -1, 32)
		writer.Write(num)
		num = num[:0]
		writer.WriteByte(slash)
		num = strconv.AppendFloat(num, float64(mean), 'f', -1, 32)
		writer.Write(num)
		num = num[:0]
		writer.WriteByte(slash)
		num = strconv.AppendFloat(num, float64(info.Max), 'f', -1, 32)
		writer.Write(num)
		num = num[:0]
		writer.Write(separators)
	}

	// Removing last ,<space>
	writer.Truncate(writer.Len() - 2)
	writer.WriteByte('}')
	writer.WriteByte(delimiter)

	return writer.Bytes(), hashmap.Size, nil
}

func read(reader *mmap.ReaderAt, hashmap *HashMap) error {
	var size int64 = int64(reader.Len())
	var at int64

	buffer := make([]byte, bufferSize)

	for at < size {
		n, err := reader.ReadAt(buffer, at)
		if err != nil && err != io.EOF {
			return err
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

		process(buffer[:n], hashmap)
		at += int64(n)
	}

	return nil
}

func process(buffer []byte, hashmap *HashMap) {
	name := make([]byte, 0, 100)

	atNumber := false
	var temperature float32
	negative := false

	separatorIndex := -1
	read := 0
	for i, b := range buffer {
		if atNumber {
			if b >= '0' && b <= '9' {
				temperature = temperature*10 + float32(b-'0')
				continue
			}

			if b == delimiter {
				read = i + 1

				hash := hashmap.Hash(name)
				info, ok := hashmap.Get(hash)
				if !ok {
					var nameCopy [100]byte
					copy(nameCopy[:], name)

					info = Info{
						Name:  nameCopy,
						Hash:  hash,
						Count: 0,
						Total: 0,
						Min:   0,
						Max:   0,
					}
				}

				if negative {
					temperature *= -1
				}

				temperature /= 10

				info.Count++
				info.Total += temperature
				if temperature > info.Max {
					info.Max = temperature
				}

				if temperature < info.Min {
					info.Min = temperature
				}

				hashmap.Set(hash, info)

				// Reset
				temperature = 0
				negative = false
				atNumber = false
			}

			continue
		}

		if b == separator {
			separatorIndex = i
			name = buffer[read:separatorIndex]
			read = i

			// For the next iteration
			atNumber = true
			if buffer[i+1] == minus {
				negative = true
			}
		}
	}
}

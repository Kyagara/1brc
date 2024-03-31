package calculate

import (
	"bytes"
	"math"
	"strconv"
	"sync"
)

func writeStation(writer *bytes.Buffer, station Station, mean float32, num []byte) {
	writer.Write(station.Name[:])
	writer.WriteByte(equal)

	num = num[:0]

	num = strconv.AppendFloat(num, float64(station.Min), 'f', -1, 32)
	writer.Write(num)

	num = num[:0]
	writer.WriteByte(slash)
	num = strconv.AppendFloat(num, float64(mean), 'f', 1, 32)
	writer.Write(num)

	num = num[:0]
	writer.WriteByte(slash)
	num = strconv.AppendFloat(num, float64(station.Max), 'f', -1, 32)
	writer.Write(num)

	writer.Write(separators)
}

func roundTowardsPositive(f float32) float32 {
	pow := math.Pow(10, 1)
	return float32(math.Round(float64(f)*pow) / pow)
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, bufferSize)
	},
}

func getBuffer() []byte {
	buf := bufferPool.Get()
	if b, ok := buf.(*[]byte); ok && b != nil {
		return *b
	}
	return make([]byte, bufferSize)
}

func putBuffer(b []byte) {
	b = b[:0]
	bufferPool.Put(&b)
}

package calculate

import (
	"math"
	"sync"
)

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

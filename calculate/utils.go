package calculate

import (
	"math"
	"strconv"
)

func roundTowardsPositive(f float32) float32 {
	pow := math.Pow(10, 1)
	return float32(math.Round(float64(f)*pow) / pow)
}

var float32Bytes = make([]byte, 32)

func float32ToBytes(f float32) []byte {
	// Convert float32 to string
	str := strconv.FormatFloat(float64(f), 'f', 1, 32)

	// Copy string bytes to the pre-allocated buffer
	copy(float32Bytes, str)

	// Return the bytes up to the length of the string
	return float32Bytes[:len(str)]
}

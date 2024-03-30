package calculate

import (
	"math"
)

func roundTowardsPositive(f float32) float32 {
	pow := math.Pow(10, 1)
	return float32(math.Round(float64(f)*pow) / pow)
}

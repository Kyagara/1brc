package calculate

import (
	"bytes"
	"math"
	"strconv"
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

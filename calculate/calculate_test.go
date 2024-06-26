package calculate_test

import (
	"brc/calculate"
	"testing"
)

func Benchmark100M(b *testing.B) {
	path := "../data/" + "100m" + ".txt"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytes, stations, err := calculate.Run(path)
		if err != nil {
			b.Fatal(err)
		}
		if len(bytes) == 0 {
			b.Fatal(err)
		}
		if stations == 0 {
			b.Fatal(err)
		}
	}
}

func Benchmark1B(b *testing.B) {
	path := "../data/" + "1b" + ".txt"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytes, stations, err := calculate.Run(path)
		if err != nil {
			b.Fatal(err)
		}
		if len(bytes) == 0 {
			b.Fatal(err)
		}
		if stations == 0 {
			b.Fatal(err)
		}
	}
}

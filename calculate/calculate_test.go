package calculate_test

import (
	"brc/calculate"
	"testing"
)

func Benchmark1M(b *testing.B) {
	path := "../data/" + "1m" + ".txt"
	b.ResetTimer()
	err := calculate.Run(path)
	if err != nil {
		b.Fatal(err)
	}
}

func Benchmark100M(b *testing.B) {
	path := "../data/" + "100m" + ".txt"
	b.ResetTimer()
	err := calculate.Run(path)
	if err != nil {
		b.Fatal(err)
	}
}

func Benchmark1B(b *testing.B) {
	path := "../data/" + "1b" + ".txt"
	b.ResetTimer()
	err := calculate.Run(path)
	if err != nil {
		b.Fatal(err)
	}
}

package main

import (
	"brc/calculate"
	"fmt"
	"os"
	"slices"
)

func main() {
	dataset := []string{"1m", "100m", "1b"}

	if len(os.Args) != 2 {
		panic("Not enough arguments, example: ./main <1m | 100m | 1b>")
	}

	option := os.Args[1]
	if !slices.Contains(dataset, option) {
		panic(fmt.Errorf("invalid dataset. Must be one of: %v, got: %v", dataset, option))
	}

	path := "./data/" + option + ".txt"
	err := calculate.Run(path)
	if err != nil {
		panic(err)
	}
}

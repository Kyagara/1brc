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
	stations, err := calculate.Run(path)
	if err != nil {
		panic(err)
	}

	output := "{"
	for i, station := range stations {
		if i == len(stations)-1 {
			output += fmt.Sprintf("%s=%.1f/%.1f/%.1f}", station.Name, station.Min, station.Mean, station.Max)
		} else {
			output += fmt.Sprintf("%s=%.1f/%.1f/%.1f,", station.Name, station.Min, station.Mean, station.Max)
		}
	}
	print(output)
}

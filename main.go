package main

import (
	"brc/calculate"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"time"
)

func main() {
	stats := false
	now := time.Now()

	dataset := []string{"1m", "100m", "1b"}

	argsLength := len(os.Args)
	if argsLength < 2 {
		panic("Not enough arguments, example: ./main <1m | 100m | 1b>")
	}

	if argsLength > 2 && os.Args[2] == "-s" {
		stats = true
		now = time.Now()
	}

	option := os.Args[1]
	if !slices.Contains(dataset, option) {
		panic(fmt.Errorf("invalid dataset. Must be one of: %v, got: %v", dataset, option))
	}

	path := "./data/" + option + ".txt"
	b, stations, err := calculate.Run(path)
	if err != nil {
		panic(err)
	}

	if stats {
		mem := runtime.MemStats{}
		runtime.ReadMemStats(&mem)

		fmt.Printf("Time: %.2fs\tMemory: %dmb\tStations: %d\n", time.Since(now).Seconds(), mem.Sys/1024/1024, stations)
		fmt.Printf("Mallocs: %d\tFrees: %d\tGC cycles: %d\n", mem.Mallocs, mem.Frees, mem.NumGC)

		return
	}

	io.Copy(os.Stdout, bytes.NewReader(b))
}

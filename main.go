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
	debug := false
	now := time.Now()

	if len(os.Args) < 2 {
		panic("Not enough arguments, example: ./brc <file>")
	}

	if slices.Contains(os.Args, "-d") {
		debug = true
		now = time.Now()
	}

	if slices.Contains(os.Args, "-s") {
		stats = true
		now = time.Now()
	}

	output, stations, err := calculate.Run(os.Args[1])
	if err != nil {
		panic(err)
	}

	if stats {
		printStats(now, stations)
		return
	}

	io.Copy(os.Stdout, bytes.NewReader(output))

	if debug {
		printStats(now, stations)
	}
}

func printStats(now time.Time, stations uint32) {
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	fmt.Printf("Time: %.2fs\tMemory: %dmb\tStations: %d\n", time.Since(now).Seconds(), mem.Sys/1024/1024, stations)
	fmt.Printf("Mallocs: %d\tFrees: %d\tGC cycles: %d\n", mem.Mallocs, mem.Frees, mem.NumGC)
}

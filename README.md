# 1brc

Having fun with the [1brc](https://github.com/gunnarmorling/1brc) challenge in go. Im not attempting to upload this somewhere, just trying different things and see what I can do.

Development and benchmarking on Windows because I don't want to accidentaly kill my wsl instance with a OOM that is bound to occur at some point.

Files used to create the dataset of 1m, 100m and 1b inside of `tools` folder.

## Challenge excerpts

Copy paste of some of the rules and how the output should look like.

### Rules and limits

- No external library dependencies may be used

- The computation must happen at application _runtime_, i.e. you cannot process the measurements file at _build time_ and just bake the result into the binary

- Input value ranges are as follows:

  - Station name: non null UTF-8 string of min length 1 character and max length 100 bytes, containing neither `;` nor `\n` character (i.e. this could be 100 one-byte characters, or 50 two-byte characters, etc.)
  - Temperature value: non null double between -99.9 (inclusive) and 99.9 (inclusive), always with one fractional digit

- Implementations must not rely on specifics of a given data set, e.g. any valid station name as per the constraints above and any data distribution (number of measurements per station) must be supported

- The rounding of output values must be done using the semantics of IEEE 754 rounding-direction "roundTowardPositive"

### Output

The task is to write a ~Java~ go program which reads the file, calculates the min, mean, and max temperature value per weather station, and emits the results on stdout like this (i.e. sorted alphabetically by station name, and the result values per station in the format <min>/<mean>/<max>, rounded to one fractional digit):

```
{Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, Accra=-10.1/26.4/66.4, Addis Ababa=-23.7/16.0/67.0, Adelaide=-27.8/17.3/58.5, ...}
```

## Commands

Saving some commands I use here.

```go
go test ./calculate -bench=Benchmark1M -benchtime=1x -benchmem -memprofile mem.out -cpuprofile cpu.out
go tool pprof -http :3000 cpu.out
go run . 1m
```

## Benchmark Results

Benchmarking the process of reading, calculating and sorting the data, basically, everything that happens inside the `calculate.Run()` function.

### V1

Jesus christ. The amount of allocs upsets me, the conversion of []byte to float32 gives me anxiety, and the fact that I used unsafe in my first attempt is a sign of things to come.

```
Benchmark1M-16                 1         216147500 ns/op        23946216 B/op      10073 allocs/op
PASS
ok      brc/calculate   0.339s

Benchmark100M-16               1        24891124400 ns/op       2613286360 B/op    73236 allocs/op
PASS
ok      brc/calculate   25.289s

go test ./calculate -bench=Benchmark1B -benchtime=1x -benchmem -memprofile mem.out -cpuprofile cpu.out
runtime: VirtualAlloc of 524288 bytes failed with errno=1455
fatal error: out of memory
```

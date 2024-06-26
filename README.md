# 1brc

Having fun with the [1brc](https://github.com/gunnarmorling/1brc) challenge in go. Im not attempting to upload this somewhere, just trying different things and see what I can do.

~Development and benchmarking on Windows because I don't want to accidentaly kill my WSL instance with a OOM that is bound to occur at some point.~ I hate Windows.

Tests are not currently being done to validate the correctness of the solutions.

Only the `README.md` from the main branch will be up-to-date.

## Challenge excerpts

Copy and paste of some of the rules and how the output should look like.

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

```bash
go test ./calculate -bench=Benchmark100M -benchtime=1x -benchmem -memprofile mem.out -cpuprofile cpu.out # create pprof profiles
go tool pprof -http :3000 cpu.out
go run . path # add '-s' to print only stats, '-d' to print the output and stats
```

## Benchmark Results

Benchmarking the process of reading, calculating and sorting the data, basically, everything that happens inside the `calculate.Run()` function.

Allocs in the new way of getting results (just using runtime package) from v2 and above is not a 1:1 from the old way, I don't know how to get the same allocs/op from the go benchmark with the `-benchtime=1x -benchmem` flags.

V2 and before had `-memprofile mem.out -cpuprofile cpu.out` flags set in the benchmark command, usually adding around 40 allocs/op.

There was a 1gb buffer set before v2.1 when reading the file which is why the memory usage up until v2.1 are pretty similar.

### v4.1

Using more for ranges. Forgot to use 1 in AppendFloat for the output at the Min and Max fields. Merging read and process functions, this improved my results by around 2s.

For now this is it, bottlenecks looks like to be the hashmap (cache miss when looking up a station) and the temperature calculation/parsing. A smaller struct and maybe a different approach on the hashmap itself might be the next step. There is a lack of unsafe, which is good thing for me.

Other than those things and the use of the experimental mmap package, I am happy with the progress, learned more about go and other things and had a lot of fun optimizing.

```
go build

./brc ./data/1b.txt -d
{Abha=-33.3/18.0/68.9, Abidjan=-20.6/26.0/73.5, ...}
Time: 19.81s    Memory: 216mb   Stations: 413
Mallocs: 8460   Frees: 7154     GC cycles: 137

time ./brc ./data/1b.txt > /dev/null
real    0m19.577s
user    1m9.453s
sys     0m5.433s
```

### v4

Using goroutines and channels. Hashmap and process workers using shards to avoid using a lock. Changed buffer size, again, now at 2gb, works well for me.

How can I avoid creating a new buffer for every read? Increasing the buffer size helps with gc since there will be less allocations and less garbage to collect.

I will probably have to start using unsafe, temperature conversion and updating the station is taking some cpu time.

```
Time: 22.85s Memory: 237mb Stations: 413
Mallocs: 6980 Frees: 6554 GC cycles: 121
```

### v3

Using mmap, from the golang experimental package, that is technically a external library, I don't care, fow now.

Custom hash map, incomplete, theres collision issues, size had to be increased to avoid it, which is why the memory usage is 4mb higher.

Some functions either moved or rewritten to accommodate the above, this will also help when adding goroutines and channels next version.

The writing of the results (from creation of the bytes.Buffer to copying it to stdout), is probably really slow, it works for now.

```
Time: 77.33s Memory: 10mb Stations: 413
Mallocs: 162 Frees: 16 GC cycles: 1
```

### v2.1

Added a len(stations) to stats. Read buffer is now a fixed 65kb size. Some micro optimizations to functions. Rework on readLines and hash.

More importantly, made new datasets, all have 413 unique stations, I will only include 1b benchmark results from now on.

Using Linux from now on.

```
Time: 97.12s    Memory: 6mb     Stations: 413
Mallocs: 125    Frees: 3        GC cycles: 0
```

### v2

No unsafe, for now. Dealing with min, mean and max in separated `Station` fields, instead of having a single []float32, that was a horrible idea. Station.Name is now a [100]byte. When using []byte, there usually needs to have a copy of the data to avoid garbage.

First time 1B works!

Including the old way of benchmark results and the new one.

```
Results for '1m':
Time: 0.14s     System Memory: 10mb
Mallocs: 235    Frees: 7        GC cycles: 0

Results for '100m':
Time: 11.72s    System Memory: 987mb
Mallocs: 273    Frees: 45       GC cycles: 1

Results for '1b':
Time: 156.70s   System Memory: 987mb
Mallocs: 277    Frees: 52       GC cycles: 2

Benchmark1M-16       144884800 ns/op    3257648 B/op 68 allocs/op
0.355s

Benchmark100M-16   12598715300 ns/op 1027298936 B/op 78 allocs/op
12.805s

Benchmark1B-16    168614731300 ns/op 1027350960 B/op 96 allocs/op
168.857s
```

### v1.1

I was, for some reason, sorting the temperatures. Calculating the min, mean and max is now done in a single loop. I also added a custom parseFloat32 function, no more anxiety.

```
Benchmark1M-16     140032600 ns/op   23945352 B/op 10068 allocs/op
0.338s

Benchmark100M-16 14905148400 ns/op 2613198056 B/op 73213 allocs/op
15.203s

Benchmark1B-can you guess?
```

### v1

Jesus christ. The amount of allocs upsets me, the conversion of []byte -> string -> float32 gives me anxiety and the fact that I used unsafe in my first attempt is a sign of things to come.

```
Benchmark1M-16     216147500 ns/op   23946216 B/op 10073 allocs/op
0.339s

Benchmark100M-16 24891124400 ns/op 2613286360 B/op 73236 allocs/op
25.289s

go test ./calculate -bench=Benchmark1B -benchtime=1x -benchmem -memprofile mem.out -cpuprofile cpu.out
runtime: VirtualAlloc of 524288 bytes failed with errno=1455
fatal error: out of memory
```

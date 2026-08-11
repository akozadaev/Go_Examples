package main

import (
	"flag"
	"fmt"
	"runtime"
)

var live [][]byte

func main() {
	blocks := flag.Int("blocks", 64, "number of allocated blocks")
	megabytes := flag.Int("mb", 1, "size of each block in MiB")
	keep := flag.Bool("keep", false, "keep blocks reachable before the last measurement")
	forceGC := flag.Bool("gc", true, "request a GC before the last measurement")
	flag.Parse()

	printStats("start")
	for range *blocks {
		block := make([]byte, *megabytes<<20)
		block[0] = 1
		live = append(live, block)
	}
	printStats("after allocation")

	if !*keep {
		live = nil
	}
	if *forceGC {
		runtime.GC() // Только для управляемой демонстрации, а не как стратегия приложения.
	}
	printStats("finish")
}

func printStats(label string) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Printf(
		"%-16s heap_alloc=%6.1f MiB heap_objects=%d num_gc=%d next_gc=%6.1f MiB\n",
		label,
		float64(stats.HeapAlloc)/(1<<20),
		stats.HeapObjects,
		stats.NumGC,
		float64(stats.NextGC)/(1<<20),
	)
}

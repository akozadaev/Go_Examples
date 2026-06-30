package main

import (
	"fmt"
	"log"

	"example.com/functions/part1/namedresults/report"
)

func main() {
	summary, err := report.BuildSummary(88, 72, 95, 64)
	if err != nil {
		log.Fatal(err)
	}

	passed, failed := report.PassFail(70, 88, 72, 95, 64)

	fmt.Printf("summary: min=%.1f max=%.1f avg=%.1f\n", summary.Min, summary.Max, summary.Avg)
	fmt.Printf("passed=%d failed=%d\n", passed, failed)
}

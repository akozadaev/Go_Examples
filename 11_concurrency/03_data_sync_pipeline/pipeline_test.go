package main

import "testing"

func TestPipelineSum(t *testing.T) {
	jobs := generateJobs([]int{2, 3, 4})
	results := runWorkers(2, jobs)

	total := 0
	count := 0
	for result := range results {
		total += result.Result
		count++
	}

	if count != 3 {
		t.Fatalf("получено результатов: %d, ожидается 3", count)
	}
	if total != 29 {
		t.Fatalf("сумма результатов: %d, ожидается 29", total)
	}
}

package jsonenums_example

//go:generate jsonenums -type=Priority

// Priority представляет приоритет задачи
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

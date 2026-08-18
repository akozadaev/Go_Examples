package stringer_example

//go:generate stringer -type=Status

// Status представляет статус выполнения задачи
type Status int

const (
	StatusPending Status = iota
	StatusInProgress
	StatusCompleted
	StatusCancelled
)


package mapper_example

import "time"

// parseTime преобразует строку в time.Time
// Используется сгенерированными мапперами для преобразования строк в time.Time
func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

package report

import "errors"

// ErrNoScores сообщает, что расчет сводки получил пустой список оценок.
var ErrNoScores = errors.New("no scores")

// Summary содержит базовую статистику по списку оценок.
type Summary struct {
	Min float64
	Max float64
	Avg float64
}

// BuildSummary вычисляет минимальную, максимальную и среднюю оценку.
func BuildSummary(scores ...float64) (summary Summary, err error) {
	if len(scores) == 0 {
		return summary, ErrNoScores
	}

	summary.Min = scores[0]
	summary.Max = scores[0]

	total := 0.0
	for _, score := range scores {
		if score < summary.Min {
			summary.Min = score
		}
		if score > summary.Max {
			summary.Max = score
		}
		total += score
	}

	summary.Avg = total / float64(len(scores))
	return summary, nil
}

// PassFail возвращает количество оценок, которые больше или равны threshold, и количество остальных оценок.
func PassFail(threshold float64, scores ...float64) (passed int, failed int) {
	for _, score := range scores {
		if score >= threshold {
			passed++
			continue
		}
		failed++
	}
	return passed, failed
}

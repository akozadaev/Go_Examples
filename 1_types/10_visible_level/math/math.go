package math

const Pi = 3.14           // экспортируемая
const internalVersion = 2 // неэкспортируемая

func Add(a, b int) int      { return a + b } // экспортируемая
func multiply(a, b int) int { return a * b } // неэкспортируемая

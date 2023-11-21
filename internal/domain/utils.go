package domain

// Decimal есть только в виде библиотеки, поэтому пока что здесь порнография с int.
// Чтобы наверняка, можно использовать uint64 (нам нельзя иметь отрицательные суммы)

func MonetaryToFloat(monetary int) float64 {
	return float64(monetary) / 100
}

func FloatToMonetary(value float64) int {
	return int(value * 100)
}

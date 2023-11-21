package domain

// Decimal есть только в виде библиотеки, поэтому пока что здесь порнография с int.
// Чтобы наверняка, можно использовать uint64 (нам нельзя иметь отрицательные суммы)

func MonetaryToFloat(monetary int) float32 {
	return float32(monetary) / 100
}

func FloatToMonetary(value float32) int {
	return int(value * 100)
}

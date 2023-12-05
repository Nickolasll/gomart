package domain

func MonetaryToFloat(monetary int) float64 {
	return float64(monetary) / 100
}

func FloatToMonetary(value float64) int {
	return int(value * 100)
}

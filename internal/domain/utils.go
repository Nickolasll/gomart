package domain

import (
	"strconv"
	"strings"
)

// Decimal есть только в виде библиотеки, поэтому пока что здесь порнография с int.
// Чтобы наверняка, можно использовать uint64 (нам нельзя иметь отрицательные суммы)

func MonetaryToString(monetary int) string {
	exp := monetary % 100
	value := monetary / 100
	return strconv.Itoa(value) + "." + strconv.Itoa(exp)
}

func StringToMonetary(sValue string) int {
	values := strings.Split(sValue, ".")
	value, _ := strconv.Atoi(values[0])
	value *= 100
	exp, _ := strconv.Atoi(values[1])
	return value + exp
}

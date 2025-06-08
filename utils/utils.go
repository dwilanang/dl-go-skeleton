package utils

import (
	"math"
	"strconv"
)

func ConvertStringToInt(str string) int64 {
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return int64(res)
}

func RoundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

package utils

import "math"

func TotalPage(limit int, totalData int64) int {
	if totalData <= 0 {
		return 0
	}

	flimit := float64(limit)
	fdata := float64(totalData)

	res := math.Ceil(fdata / flimit)

	return int(res)
}

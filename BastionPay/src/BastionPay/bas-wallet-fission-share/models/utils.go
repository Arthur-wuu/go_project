package models

import "math"

const (
	Const_Level_Tourist    = 1
	Const_Level_Login_Free = 10
)

func AdjustFloatAcc(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

package controllers

import "time"

const (
	CONST_MIN_RED_DIV_BASE = 10
)

func GetTodayUnix() (int64, int64) {
	ll := time.FixedZone("UTC", 8*3600)
	temp := time.Now()
	tt := temp.In(ll)
	return time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, ll).Unix(), temp.Unix()
}

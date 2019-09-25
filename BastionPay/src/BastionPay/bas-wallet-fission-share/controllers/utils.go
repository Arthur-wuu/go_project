package controllers

import (
	"BastionPay/bas-wallet-fission-share/models"
	"github.com/shopspring/decimal"
	"math/rand"
)

func GenRedMoney(red *models.Red) (*decimal.Decimal, bool) {
	if *red.Mode == models.Const_Share_Mode_Fixed {
		return GetAverageMoney(red)
	} else {
		return GetRandomMoney(red)
	}
}

func GetRandomMoney(red *models.Red) (*decimal.Decimal, bool) {
	if *red.RemainRob == 0 {
		tmp := decimal.NewFromFloat(0)
		return &tmp, false
	}
	if *red.RemainRob == 1 {
		*red.RemainRob--
		tmp := red.RemainCoin
		zero := decimal.NewFromFloat(0)
		red.RemainCoin = &zero
		return tmp, true
	}
	average := red.RemainCoin.Div(decimal.New(*red.RemainRob, 0))
	min := red.TotalCoin.Div(decimal.New(*red.TotalRob, 0))
	min = min.Div(decimal.New(100, 0)).Truncate(*red.Precision)

	cmp := min.Mul(decimal.New(*red.RemainRob+1, 0)).Cmp(*red.RemainCoin)
	if cmp >= 0 {
		*red.RemainRob--
		*red.RemainCoin = red.RemainCoin.Sub(min)
		return &min, true
	}
	max := average.Mul(decimal.New(2, 0)).Truncate(*red.Precision)

	r := max.Sub(min).Mul(decimal.NewFromFloat(rand.Float64())).Add(min).Truncate(*red.Precision)
	*red.RemainRob--
	*red.RemainCoin = red.RemainCoin.Sub(r)
	return &r, true
}

func GetAverageMoney(red *models.Red) (*decimal.Decimal, bool) {
	if *red.RemainRob == 0 {
		tmp := decimal.NewFromFloat(0)
		return &tmp, false
	}
	if *red.RemainRob == 1 {
		*red.RemainRob--
		tmp := red.RemainCoin
		zero := decimal.NewFromFloat(0)
		red.RemainCoin = &zero
		return tmp, true
	}
	average := red.RemainCoin.Div(decimal.New(*red.RemainRob, 0)).Truncate(*red.Precision)

	*red.RemainRob--
	*red.RemainCoin = red.RemainCoin.Sub(average)
	return &average, true
}

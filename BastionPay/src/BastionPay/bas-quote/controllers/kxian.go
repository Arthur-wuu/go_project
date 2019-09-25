package controllers

import (
	"BastionPay/bas-api/apibackend"
	apiquote "BastionPay/bas-api/quote"
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote/quote"
	. "BastionPay/bas-quote/utils"
	"fmt"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func NewKXianCtl(q *quote.QuoteMgr) *KXianCtl {
	return &KXianCtl{
		mQuote: q,
	}
}

type KXianCtl struct {
	mQuote *quote.QuoteMgr
}

func (this *KXianCtl) GetKXian(ctx iris.Context) {
	defer PanicPrint()
	period, start, limit, fromArr, toArr, err := this.parseKXianParam(ctx)
	if err != nil {
		ZapLog().With(zap.Any("urlparams", ctx.URLParams())).Error("param err")
		ctx.JSON(*apiquote.NewResMsg(apibackend.BASERR_INVALID_PARAMETER.Code(), apibackend.BASERR_INVALID_PARAMETER.Desc()))
		return
	}
	resMsg := apiquote.NewResMsg(apibackend.BASERR_SUCCESS.Code(), "")

	for i := 0; i < len(fromArr); i++ {
		if len(fromArr[i]) == 0 {
			continue
		}
		codeInfo := this.mQuote.GetSymbol(fromArr[i])
		if codeInfo == nil {
			ZapLog().With(zap.String("from", fromArr[i])).Error("GetSymbol nofind err")
			continue
		}
		info := resMsg.GenHistoryDetailInfo()
		info.Id = codeInfo.Id
		info.Symbol = codeInfo.Symbol
		for j := 0; j < len(toArr); j++ {
			if len(toArr[j]) == 0 {
				continue
			}
			to := toArr[j]
			switch period {
			case "1DAY":
				if codeInfo.GetId() < 100000 {
					fmt.Println("go here 1 codeInfo.GetId() :", codeInfo.GetId())
					if to == "BTC" || to == "ETH" || to == "XRP" || to == "LTC" || to == "BCH" || to == "USD" {
						kxians, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), toArr[j], start, limit)
						if err != nil {
							ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
							break
						}
						kxianDetail := info.AddKXianDetail()
						kxianDetail.KXians = ToApiKXians(kxians)
						kxianDetail.SetSymbol(toArr[j])
						break
					}
					if to != "BTC" && to != "ETH" && to != "XRP" && to != "LTC" && to != "BCH" && to != "USD" {
						fmt.Println("go here 2 toArr[j] :", toArr[j])
						if toArr[j] == "USD" {
							kxian1, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), "BTC", start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}
							kxian2, err := this.mQuote.GetKxian1Day(1, "USD", start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}
							kxianArr := make([]*quote.KXian, 0)
							if len(kxian1) > 0 && len(kxian2) > 0 {
								for i := 0; i < min(len(kxian1), len(kxian2)); i++ {
									kxianTemp := new(quote.KXian)
									kxianTemp.SetClosePrice(kxian1[i].GetClosePrice() * kxian2[i].GetClosePrice())
									kxianTemp.SetHighPrice(kxian1[i].GetHighPrice() * kxian2[i].GetHighPrice())
									kxianTemp.SetLastPrice(kxian1[i].GetLastPrice() * kxian2[i].GetLastPrice())
									kxianTemp.SetLowPrice(kxian1[i].GetLowPrice() * kxian2[i].GetLowPrice())
									kxianTemp.SetOpenPrice(kxian1[i].GetOpenPrice() * kxian2[i].GetOpenPrice())
									kxianTemp.SetTimestamp(kxian2[i].GetTimestamp())

									fmt.Println("aaaaa  ", kxian1[i].GetTimestamp(), kxian2[i].GetTimestamp())
									kxianArr = append(kxianArr, kxianTemp)
								}
							}
							kxianDetail := info.AddKXianDetail()
							kxianDetail.KXians = ToApiKXians(kxianArr)
							kxianDetail.SetSymbol(toArr[j])
							break
						} else {
							//kxian1, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), "BTC", start, limit)
							//if err != nil {
							//	ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
							//	break
							//}
							kxian2, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), "USD", start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}
							kxian3, err := this.mQuote.GetKxian1DayUseUSD(toArr[j], start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}

							kxianArr := make([]*quote.KXian, 0)
							if len(kxian2) > 0 && len(kxian3) > 0 {
								for i := 0; i < min(len(kxian2), len(kxian3)); i++ {
									kxianTemp := new(quote.KXian)

									kxianTemp.SetClosePrice(kxian2[i].GetClosePrice() * kxian3[i].GetClosePrice())
									kxianTemp.SetHighPrice(kxian2[i].GetHighPrice() * kxian3[i].GetHighPrice())
									kxianTemp.SetLastPrice(kxian2[i].GetLastPrice() * kxian3[i].GetLastPrice())
									kxianTemp.SetLowPrice(kxian2[i].GetLowPrice() * kxian3[i].GetLowPrice())
									kxianTemp.SetOpenPrice(kxian2[i].GetOpenPrice() * kxian3[i].GetOpenPrice())
									kxianTemp.SetTimestamp(kxian2[i].GetTimestamp())

									kxianArr = append(kxianArr, kxianTemp)
								}
							}
							kxianDetail := info.AddKXianDetail()
							kxianDetail.KXians = ToApiKXians(kxianArr)
							kxianDetail.SetSymbol(toArr[j])
							break
						}

					}
				}

				if codeInfo.GetId() >= 100000 {
					fmt.Println("go here 1")
					if to == "BTC" || to == "ETH" || to == "XRP" || to == "LTC" || to == "BCH" {
						kxian1, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), "BTC", start, limit)
						if err != nil {
							ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
							break
						}
						kxian2, err := this.mQuote.GetKxian1Day(1, toArr[j], start, limit)
						if err != nil {
							ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
							break
						}
						kxianArr := make([]*quote.KXian, 0)
						if len(kxian1) > 0 && len(kxian2) > 0 {
							for i := 0; i < min(len(kxian1), len(kxian2)); i++ {
								kxianTemp := new(quote.KXian)
								kxianTemp.SetClosePrice(kxian1[i].GetClosePrice() * kxian2[i].GetClosePrice())
								kxianTemp.SetHighPrice(kxian1[i].GetHighPrice() * kxian2[i].GetHighPrice())
								kxianTemp.SetLastPrice(kxian1[i].GetLastPrice() * kxian2[i].GetLastPrice())
								kxianTemp.SetLowPrice(kxian1[i].GetLowPrice() * kxian2[i].GetLowPrice())
								kxianTemp.SetOpenPrice(kxian1[i].GetOpenPrice() * kxian2[i].GetOpenPrice())
								kxianTemp.SetTimestamp(kxian1[i].GetTimestamp())

								kxianArr = append(kxianArr, kxianTemp)
							}
						}
						kxianDetail := info.AddKXianDetail()
						kxianDetail.KXians = ToApiKXians(kxianArr)
						kxianDetail.SetSymbol(toArr[j])
						break
					}
					if to != "BTC" && to != "ETH" && to != "XRP" && to != "LTC" && to != "BCH" {
						if toArr[j] == "USD" {
							kxian1, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), "BTC", start, limit)
							if len(kxian1) == 0 {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}

							kxian2, err := this.mQuote.GetKxian1Day(1, "USD", start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}
							kxianArr := make([]*quote.KXian, 0)
							if len(kxian1) > 0 && len(kxian2) > 0 {
								for i := 0; i < min(len(kxian1), len(kxian2)); i++ {
									kxianTemp := new(quote.KXian)

									kxianTemp.SetClosePrice(kxian1[i].GetClosePrice() * kxian2[i].GetClosePrice())
									kxianTemp.SetHighPrice(kxian1[i].GetHighPrice() * kxian2[i].GetHighPrice())
									kxianTemp.SetLastPrice(kxian1[i].GetLastPrice() * kxian2[i].GetLastPrice())
									kxianTemp.SetLowPrice(kxian1[i].GetLowPrice() * kxian2[i].GetLowPrice())
									kxianTemp.SetOpenPrice(kxian1[i].GetOpenPrice() * kxian2[i].GetOpenPrice())
									kxianTemp.SetTimestamp(kxian2[i].GetTimestamp())

									kxianArr = append(kxianArr, kxianTemp)
								}
							}
							kxianDetail := info.AddKXianDetail()
							kxianDetail.KXians = ToApiKXians(kxianArr)
							kxianDetail.SetSymbol(toArr[j])
							break

						} else {
							kxian1, err := this.mQuote.GetKxian1Day(codeInfo.GetId(), "BTC", start, limit)
							if len(kxian1) == 0 {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}
							kxian2, err := this.mQuote.GetKxian1Day(1, "USD", start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}
							kxian3, err := this.mQuote.GetKxian1DayUseUSD(toArr[j], start, limit)
							if err != nil {
								ZapLog().With(zap.Error(err), zap.Int("id", codeInfo.GetId()), zap.String("to", toArr[j]), zap.String("from", fromArr[i])).Error("GetKxian1Day err")
								break
							}

							kxianArr := make([]*quote.KXian, 0)
							if len(kxian1) > 0 && len(kxian2) > 0 && len(kxian3) > 0 {
								for i := 0; i < min3(len(kxian1), len(kxian2), len(kxian3)); i++ {
									kxianTemp := new(quote.KXian)
									kxianTemp.SetClosePrice(kxian1[i].GetClosePrice() * kxian2[i].GetClosePrice() * kxian3[i].GetClosePrice())
									kxianTemp.SetHighPrice(kxian1[i].GetHighPrice() * kxian2[i].GetHighPrice() * kxian3[i].GetHighPrice())
									kxianTemp.SetLastPrice(kxian1[i].GetLastPrice() * kxian2[i].GetLastPrice() * kxian3[i].GetLastPrice())
									kxianTemp.SetLowPrice(kxian1[i].GetLowPrice() * kxian2[i].GetLowPrice() * kxian3[i].GetLowPrice())
									kxianTemp.SetOpenPrice(kxian1[i].GetOpenPrice() * kxian2[i].GetOpenPrice() * kxian3[i].GetOpenPrice())
									kxianTemp.SetTimestamp(kxian1[i].GetTimestamp())
									kxianArr = append(kxianArr, kxianTemp)
								}
							}
							kxianDetail := info.AddKXianDetail()
							kxianDetail.KXians = ToApiKXians(kxianArr)
							kxianDetail.SetSymbol(toArr[j])
							break
						}
					}
				}
			}
		}
	}
	ctx.JSON(resMsg)
}

func (this *KXianCtl) parseKXianParam(ctx iris.Context) (string, int, int, []string, []string, error) {
	period := strings.ToUpper(ctx.URLParam("period"))
	limitStr := strings.ToUpper(ctx.URLParam("limit"))
	startStr := strings.ToUpper(ctx.URLParam("start"))
	//from := strings.ToUpper(ctx.URLParam("from"))
	from := ctx.URLParam("from")
	to := strings.ToUpper(ctx.URLParam("to"))
	if len(period) == 0 {
		period = "1DAY"
	}
	if len(to) == 0 {
		to = "USD"
	}
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	from = strings.TrimRight(from, ",")
	to = strings.TrimRight(to, ",")
	fromArr := strings.Split(from, ",")
	toArr := strings.Split(to, ",")

	limit := 30
	if len(limitStr) != 0 {
		n, err := strconv.Atoi(limitStr)
		if err != nil {
			return "", 0, 0, nil, nil, err
		} else {
			limit = n
		}
	}
	start := 0
	if len(startStr) != 0 {
		n, err := strconv.Atoi(startStr)
		if err != nil {
			return "", 0, 0, nil, nil, err
		} else {
			start = n
		}
	}
	return period, start, limit, fromArr, toArr, nil
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
func min3(a, b, c int) int {
	if a >= b && a >= c {
		if b >= c {
			return c
		}
		return b
	}
	if b >= c && b >= a {
		if a >= c {
			return c
		}
		return a
	}
	if c >= b && c >= a {
		if b >= a {
			return a
		}
		return b
	}
	return 1
}

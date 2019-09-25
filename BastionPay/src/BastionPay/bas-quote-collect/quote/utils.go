package quote

import (
	"BastionPay/bas-quote-collect/collect"
	"BastionPay/bas-quote-collect/db"
	//"github.com/bugsnag/bugsnag-go/errors"
	. "BastionPay/bas-base/log/zap"
	"errors"
	"go.uber.org/zap"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"
)

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

func CodeInfoToCodeTable(info *collect.CodeInfo) *db.CodeTable {
	table := new(db.CodeTable)

	if info.Symbol != nil {
		table.Symbol = *info.Symbol
	}

	table.Name = info.Name
	if info.Id != nil {
		table.Code = uint(*info.Id)
	}

	table.WebsiteSlug = info.Website_slug

	if info.Timestamp != nil {
		table.UpdatedAt = info.Timestamp
	}
	table.Valid = info.Valid

	return table
}

func CodeTableToCodeInfo(table *db.CodeTable) *collect.CodeInfo {
	info := new(collect.CodeInfo)

	info.Symbol = &table.Symbol

	info.Id = new(int)
	*info.Id = int(table.Code)
	info.Name = table.Name

	info.Website_slug = table.WebsiteSlug
	info.Timestamp = table.UpdatedAt
	info.Valid = table.Valid

	return info
}

func NowTimestamp() int64 {
	return time.Now().Unix()
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func BubbleSort(values []int) []int {
	var arrlen int = len(values)
	for i := 0; i < arrlen-1; i++ {

		for j := 0; j < len(values)-1-i; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
			}
		}
	}
	return values
}

func StringArrToIntArr(strArr []string) ([]int, error) {

	intArr := make([]int, 0)
	for i := 0; i < len(strArr); i++ {
		i, err := strconv.Atoi(strArr[i])
		if err != nil {
			return nil, errors.New("StringArrToIntArr err")
		}
		intArr = append(intArr, i)
	}
	return intArr, nil
}

func ArrToString(idInt []int) string {
	param := ""
	for i := 0; i < len(idInt); i++ {
		if i == len(idInt)-1 {
			param = param + strconv.Itoa(idInt[i])
			break
		}
		param = param + strconv.Itoa(idInt[i]) + ","
	}
	return param
}

func TimeToTimestamp(t string) int64 {
	datetime := t //待转化为时间戳的字符串

	//日期转化为时间戳
	timeLayout := "2006-01-02T15:04:05.000Z" //转化所需模板
	loc, _ := time.LoadLocation("GMT")       //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix()
	return timestamp
}

package quote

import (
	"BastionPay/bas-quote/collect"
	"BastionPay/bas-quote/db"
	"math/rand"
	"time"
)

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

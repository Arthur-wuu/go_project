package main

import (
	"BastionPay/bas-api/quote"
	"BastionPay/bas-quote/common"
	"encoding/json"
	"fmt"
	"time"
)

func createXlsFile(filePath string, rowArr []interface{}, colNames []string) error {
	xlsObj, err := common.NewXlsx(rowArr, colNames, nil)
	if err != nil {
		return err
	}
	if err = xlsObj.Generate(); err != nil {
		return err
	}
	if err = xlsObj.File(filePath); err != nil {
		return err
	}
	return nil
}

func main() {
	codesBytes, err := common.HttpSend("http://quote.rkuan.com/api/v1/coin/code", nil, "GET", nil)
	if err != nil {
		fmt.Println("err=", err)
		return
	}

	codeList := new(quote.ResMsg)
	if err := json.Unmarshal(codesBytes, codeList); err != nil {
		fmt.Println("err=", err)
		return
	}
	if codeList.Err != 0 {
		fmt.Println("err=", codeList.Err)
		return
	}
	rowArr := make([]interface{}, 0)
	colNames := make([]string, 0)
	colNames = append(colNames, "symbol")
	colNames = append(colNames, "usd_price")
	colNames = append(colNames, "valid")
	colNames = append(colNames, "last_update_time")
	for i := 0; i < len(codeList.Codes); i++ {
		if codeList.Codes[i].Symbol == nil {
			continue
		}
		fmt.Println("start ", i, *codeList.Codes[i].Symbol)
		symbol := *codeList.Codes[i].Symbol
		valid := "valid"
		if codeList.Codes[i].Valid == nil || *codeList.Codes[i].Valid == 0 {
			valid = "in-valid"
		}
		quoteBytes, err := common.HttpSend("http://quote.rkuan.com/api/v1/coin/quote?from="+symbol+"&to=USD", nil, "GET", nil)
		if err != nil {
			fmt.Println("err=", err)
			return
		}
		codequote := new(quote.ResMsg)
		if err := json.Unmarshal(quoteBytes, codequote); err != nil {
			fmt.Println("err=", err)
			return
		}
		if codequote.Err != 0 {
			fmt.Println("err=", codeList.Err)
			continue
		}
		if len(codequote.Quotes) <= 0 {
			fmt.Println("err= noquote", symbol)
			continue
		}
		if len(codequote.Quotes[0].MoneyInfos) == 0 {
			fmt.Println("err= noquote", symbol)
			continue
		}
		if codequote.Quotes[0].MoneyInfos[0].Price == nil {
			fmt.Println("err= noquote", symbol)
			continue
		}
		UpdateTime := ""
		if codequote.Quotes[0].MoneyInfos[0].Last_updated != nil {
			UpdateTime = time.Unix(*codequote.Quotes[0].MoneyInfos[0].Last_updated, 0).In(time.FixedZone("UTC", 8*3600)).Format("2006-01-02 15:04:05")
		}

		rowArr = append(rowArr, []string{symbol, fmt.Sprintf("%v", *codequote.Quotes[0].MoneyInfos[0].Price)}, valid, UpdateTime)
	}
	fmt.Println("ok0")
	err = createXlsFile("out.xlsx", rowArr, colNames)

	fmt.Println("ok", err)
}

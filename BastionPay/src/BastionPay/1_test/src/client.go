package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"

	//"golang.org/x/net/websocket"
//"github.com/gorilla/websocket"
	"fmt"


	//"log"
)

var origin = "http://127.0.0.1:1234/"
var url = "ws://127.0.0.1:1234/"

func main() {

//	type NotifyRequest struct{
//		//Id      		*int       `json:"id,omitempty"  `
//		//ActivityUuid    *string    `json:"activity_uuid,omitempty"  `
//		//RedId      		*string       `json:"red_uuid,omitempty"  `
//		//AppId     		*string    `json:"app_id,omitempty"`
//		UserId     		*string    `json:"user_id,omitempty"`
//		//CountryCode 	*string    `json:"country_code,omitempty"`
//		//Phone      		*string    `json:"phone,omitempty" `
//		Symbol    	 	*string    `json:"symbol,omitempty" `
//		Coin     		*string   `json:"coin,omitempty"`
//		//SponsorAccount  *string    `json:"sponsor_account,omitempty" `
//		//ApiKey          *string    `json:"api_key,omitempty" `
//		//OffAt           *int64     `json:"off_at,omitempty" `
//		//Lang            *string     `json:"language,omitempty" `
//		//TransferFlag 	*int       `json:"transfer_flag,omitempty"`
//	}
//	params := new(NotifyRequest)
//	c := "3.4444444444"
//	params.Coin = &c
//
//
//	coin, _ := strconv.ParseFloat( *params.Coin, 64)
//
//	fmt.Println("coin",coin)
//s := RandString(3)
//
//	fmt.Println("s",s)
	url := "a.txt"
		fileName := url
		file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Open file error!", err)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			panic(err)
		}

		var size = stat.Size()
		fmt.Println("file size=", size)

		buf := bufio.NewReader(file)
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			fmt.Println(line)
			if err != nil {
				if err == io.EOF {
					fmt.Println("File read ok!")
					break
				} else {
					fmt.Println("Read file error!", err)
					return
				}
			}
		}
	}


func ReadCsv(name string) [][]string{
	file, err := os.Open(name)
	if err != nil {
		fmt.Println("open err[%v]", err)
		return nil
	}
	defer file.Close()
	reader := csv.NewReader(file)
	allRecord := make([][]string, 0)
	//var newRecord []string
	for i:=0;true;i++{
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Read %v err[%v]",i, err)
			return nil
		}
		for j:=0; j < len(record);j++ {
			record[j] = strings.Replace(record[j], " ", "", -1)
			record[j] = strings.Replace(record[j], "\r", "", -1)
			record[j] = strings.Replace(record[j], "\n", "", -1)
		}
		allRecord = append(allRecord, record)
	}
	return allRecord
}

//var r *rand.Rand
//func init() {
//	r = rand.New(rand.NewSource(time.Now().Unix()))
//}
//
//// RandString 生成随机字符串
//func RandString(len int) string {
//	bytes := make([]byte, len)
//	for i := 0; i < len; i++ {
//		b := r.Intn(26) + 65
//		bytes[i] = byte(b)
//	}
//	return string(bytes)
//}
//
//rands := RandString(4)
//requestNoStr = requestNoStr + rands
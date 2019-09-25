package src

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)



func main() {
	bastionPayUrl := "http://xxxxxxx/push_message?content=%s&createTime=%d&historyFlag=%t&id=%d&msgType=%d&redFlag=%t&sourceId=%s&subType=%s&title=%s&userId=%d"
	fmt.Printf("%s", []byte("Go语言"))
	fmt.Print("%s", []byte("Go语言"))
	content := "aaa"
	createTime := int32(time.Now().Unix())
	historyFlag := false
	id := 1
	msgType := 1
	redFlag := true
	sourceId := "11"
	subType := "4_okg"
	title := "title"
	userId, err := strconv.Atoi("123")
	if err != nil {
		//ZapLog().Error( "account userId string to int err", zap.Error(err))
		return
	}

	fmt.Printf(bastionPayUrl,content,createTime,historyFlag,id,msgType,redFlag,sourceId,subType,title,userId)





	file, err := os.Open("im.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(record) // record has the type []string
	}
}
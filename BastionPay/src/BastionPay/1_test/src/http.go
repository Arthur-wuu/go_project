package src

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://baidu.com"
	fmt.Println("URL:>", url)

	//登陆用户名
	usrId := "LaoWong"

	//登陆密码


	//json序列化

	post := "{\"UserId\":\"" + usrId + "\"}"



	var jsonStr = []byte(post)


	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

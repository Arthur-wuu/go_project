package collect

import (
	"BastionPay/bas-quote-collect/config"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	//huilvRelativePath = "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=1美元等于多少%s&co=&resource_id=4278&t=1540264697195&cardId=4278&ie=utf8&oe=gbk&cb=op_aladdin_callback&format=json&tn=baidu&cb=jQuery11020616690161198119_1540264601010&_=1540264601014"
	huilvRelativePath = "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php?query=1美元等于多少%s&co=&resource_id=6017&cardId=6017&ie=utf8&oe=gbk&cb=op_aladdin_callback&format=json&tn=baidu&cb=jQuery110208426151593621876_1540522865066"
)

type HuiLv struct {
}

func (this *HuiLv) Init() error {
	return nil
}

//参数是每个国家的中文名称,百度
func (this *HuiLv) GetHuiLv(countryName string) (*MoneyInfo, bool, error) {

	huiLvUrl := fmt.Sprintf(huilvRelativePath, countryName)

	resp, err := http.Get(huiLvUrl)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode == 404 {
		return nil, true, nil
	}
	if resp.StatusCode != 200 {
		return nil, false, errors.New("response code != 200 err")
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	str := string(content)
	index1 := strings.Index(str, "(")
	index2 := strings.LastIndex(str, ")")
	if index1 < 0 || index2 < 0 {
		return nil, false, errors.New("get reponse content out of range")
	}

	str = str[index1+1 : index2]

	moneyInfo := new(MoneyInfo)
	huilv := new(HuiLvInfo)
	err = json.Unmarshal([]byte(str), &huilv)
	if err != nil {
		return nil, false, err
	}

	if len(huilv.Data) <= 0 {
		return nil, false, errors.New("huilv data nil err")
	}

	v3, err := strconv.ParseFloat(huilv.Data[0].Money2_num, 64)
	v4, err := strconv.ParseInt(huilv.Data[0].UpdateTime, 10, 64)
	if err != nil {
		return nil, false, err
	}
	moneyInfo.Price = &v3
	moneyInfo.Last_updated = &v4

	return moneyInfo, false, nil
}

//sina 汇率
func (this *HuiLv) GetHuiLvSina() ([]*MoneyInfo, error) {

	huiLvUrl1 := genUrl1()
	resp1, err := http.Get(huiLvUrl1)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()
	if resp1.StatusCode != 200 {
		return nil, errors.New(resp1.Status)
	}

	huiLvUrl2 := genUrl2()
	resp2, err := http.Get(huiLvUrl2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		return nil, errors.New(resp1.Status)
	}

	content1, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return nil, err
	}
	content2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return nil, err
	}

	resContent := string(content1) + string(content2)
	stringArr := strings.Split(resContent, ";") //182个，最后一个是空的

	moneyInfos := make([]*MoneyInfo, 0)
	for i := 0; i < len(stringArr); i++ {
		if len(stringArr[i]) == 0 {
			continue
		}
		moneyInfo := new(MoneyInfo)
		//截取价格price
		douhaoIndex := strings.Split(stringArr[i], ",")
		if len(douhaoIndex) < 9 {
			continue
		}
		pricef, err := strconv.ParseFloat(douhaoIndex[8], 64)
		if err != nil {
			return nil, errors.New("string to float err")
		}
		//截取时间  分时秒和日期拼接 做成时间戳
		yinhaoIndex1 := strings.Index(douhaoIndex[0], "\"")
		timeMiao := douhaoIndex[0][yinhaoIndex1+1:]

		douhaoLast := strings.LastIndex(stringArr[i], ",")
		yinhaoLast := strings.LastIndex(stringArr[i], "\"")
		timeTian := stringArr[i][douhaoLast+1 : yinhaoLast]

		times := timeTian + "T" + timeMiao
		timeLayout := "2006-01-02T15:04:05"
		cstZone := time.FixedZone("UTC", 8*3600)
		tmp, _ := time.ParseInLocation(timeLayout, times, cstZone)
		timestamp := tmp.Unix()

		moneyInfo.Price = &pricef
		moneyInfo.Symbol = &config.GPreConfig.CountryCodeArr[i]
		moneyInfo.Last_updated = &timestamp
		moneyInfos = append(moneyInfos, moneyInfo)

		//fmt.Println("price time symbol",pricef,timestamp,config.GPreConfig.CountryCodeArr[i])
	}
	return moneyInfos, nil
}

//3p3ah&list=fx_susdaed,fx_susdafn
func genUrl1() string {
	randStr := GenRandomString(5)
	str := "http://hq.sinajs.cn/rn=" + randStr + "&list="
	fx := ""

	for i := 0; (i < len(config.GPreConfig.CountryCodeArr)) && (i < 100); i++ {
		if len(config.GPreConfig.CountryCodeArr[i]) == 0 {
			continue
		}
		fx = "fx_susd" + strings.ToLower(config.GPreConfig.CountryCodeArr[i]) + ","
		str = str + fx
	}
	str = str[0 : len(str)-1]
	return str
}

func genUrl2() string {
	randStr := GenRandomString(5)
	str := "http://hq.sinajs.cn/rn=" + randStr + "&list="
	fx := ""

	for i := 100; i < len(config.GPreConfig.CountryCodeArr); i++ {
		if len(config.GPreConfig.CountryCodeArr[i]) == 0 {
			continue
		}
		fx = "fx_susd" + strings.ToLower(config.GPreConfig.CountryCodeArr[i]) + ","
		str = str + fx
	}
	str = str[0 : len(str)-1]
	return str
}

//生成随机字符串
func GenRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

package collect

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-quote-collect/base"
	. "BastionPay/bas-quote-collect/config"
	"go.uber.org/zap"

	"encoding/json"
	"fmt"
	"github.com/bugsnag/bugsnag-go/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	tickerRelativePath     = "/v2/ticker?convert=%s&sort=id&structure=array&start=%d&limit=%d"
	tickerRelativeListPath = "/v1/cryptocurrency/listings/latest?start=%d&limit=%d"
	tickerRelativePartPath = "/v1/cryptocurrency/quotes/latest?id=%s"
	listRelativePath       = "/v2/listings/"
)

type CoinMarket struct {
	ApiIndex  int
	codeTable CodeTableIntf
}

func (this *CoinMarket) Init(codeTable CodeTableIntf) error {
	this.ApiIndex = 0
	this.codeTable = codeTable
	return nil
}

func (this *CoinMarket) GetIndex() int {
	return this.ApiIndex
}

func (this *CoinMarket) NextIndex() int {
	this.ApiIndex++
	if this.ApiIndex == len(GPreConfig.CoinmarketcapApiKeys) {
		this.ApiIndex = 0
	}
	return this.ApiIndex
}

func (this *CoinMarket) GetCodeTable() (*CodeListInfo, error) {
	tickerUrl := GConfig.Coinmarketcap.CoinMarket_url + listRelativePath
	resp, err := http.Get(tickerUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(CodeListInfo)
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//参数为法币名称
func (this *CoinMarket) GetAllTicker(coin string, start, limit int) (*TickerInfo, bool, bool, error) {
	tickerUrl := GConfig.Coinmarketcap.CoinMarket_url
	tickerUrl += fmt.Sprintf(tickerRelativePath, coin, start, limit)

	resp, err := http.Get(tickerUrl)
	if err != nil {
		return nil, false, false, err
	}
	if resp.StatusCode == 404 {
		return nil, true, false, nil
	}
	if resp.StatusCode == 429 {
		return nil, false, false, errors.New(resp.Status, 0)
	}
	if resp.StatusCode != 200 {
		return nil, false, false, errors.New(resp.Status, 0)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, false, err
	}

	//	fmt.Println(string(content))

	result := new(TickerInfo)
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, false, false, err
	}
	return result, false, false, nil
}

// coinmarketcap 新  全量
func (this *CoinMarket) GetTickerList(start, limit int) (*ResCoinMarketCapAll, bool, error) {
	tickerUrl := GConfig.Coinmarketcap.CoinMarket_new_url
	tickerUrl += fmt.Sprintf(tickerRelativeListPath, start, limit)
	//测试用
	//testUrl := "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=2001&limit=100"

	//配置中的api-key  每个可支持333次查询
	apiKey := strings.Replace(GConfig.Coinmarketcap.Api_key, " ", "", -1)
	apiKeyArr := strings.Split(apiKey, ",")
	res := new(ResCoinMarketCapAll)
	for {

		content, err := base.HttpSend(tickerUrl, nil, "GET", map[string]string{"X-CMC_PRO_API_KEY": apiKeyArr[this.GetIndex()]})
		if err != nil {
			ZapLog().Error("httpSend request err", zap.Error(err), zap.String("contentData", string(content)))
			return nil, false, err
		}

		if err = json.Unmarshal(content, res); err != nil {
			ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
			return nil, false, err
		}

		if res.Status.Error_code == 0 {
			return res, false, nil
		}

		if res.Status.Error_code != 0 {
			//  换下一个API-key   再请求
			this.NextIndex()
		}
	}
	return res, false, nil

}

// coinmarketcap 新  码表里有的币
func (this *CoinMarket) GetPartTicker() (*ResCoinMarketCapPart, string, error) {
	//准备数据库里有效的码表

	if len(GPreConfig.CoinmarketcapApiKeys) == 0 {
		return nil, "", fmt.Errorf("nil coinmarketcap apikey")
	}

	codeInfo := this.codeTable.ListSymbols()
	if len(codeInfo) == 0 {
		ZapLog().Error("getAllCode err:nil codeTable")
		return nil, "", fmt.Errorf("nil codeTable")
	}
	codeInfoMap := make(map[string]*CodeInfo)

	codes := ""
	//Allcodes := ""
	for i := 0; i < len(codeInfo); i++ {
		//Allcodes += fmt.Sprintf("[%d %d]", codeInfo[i].GetId(), codeInfo[i].GetValid())
		if codeInfo[i].GetId() == 0 || codeInfo[i].GetId() >= 100000 || codeInfo[i].GetValid() != 1 {
			continue
		}
		if _, ok := codeInfoMap[strconv.Itoa(codeInfo[i].GetId())]; ok {
			continue
		}
		codes = codes + strconv.Itoa(int(codeInfo[i].GetId())) + ","
		codeInfoMap[strconv.Itoa(codeInfo[i].GetId())] = &codeInfo[i]
	}
	if len(codes) <= 1 {
		//ZapLog().Error("nil valid codeTable err:"+Allcodes)
		return nil, "", fmt.Errorf("nil valid codeTable")
	}
	codes = codes[:len(codes)-1]
	ZapLog().Debug("codes=" + codes)

	tickerUrl := GConfig.Coinmarketcap.CoinMarket_new_url + fmt.Sprintf(tickerRelativePartPath, codes)

	//配置中的api-key  每个可支持333次查询
	res := new(ResCoinMarketCapPart)

	content, err := base.HttpSend(tickerUrl, nil, "GET", map[string]string{"X-CMC_PRO_API_KEY": GPreConfig.CoinmarketcapApiKeys[this.GetIndex()]})
	if err != nil {
		ZapLog().Error("httpSend request err", zap.Error(err), zap.String("contentData", string(content)))
		//			return nil, codes ,err
	}
	if len(content) == 0 {
		return nil, codes, fmt.Errorf("coinmarketcap response nil")
	}

	if err := json.Unmarshal(content, res); err != nil {
		ZapLog().Error("unMarshal err", zap.Error(err), zap.String("data", string(content)))
		return nil, codes, err
	}
	if res.Status.Error_code == 0 {
		return res, codes, nil
	}
	ZapLog().Sugar().Errorf("code[%d] errmasg[%s]", res.Status.Error_code, res.Status.Error_message)

	if res.Status.Error_code == 400 {
		errCodeAll := this.parseInValidId(res.Status.Error_message)
		for j := 0; j < len(errCodeAll); j++ {
			codeInfo, ok := codeInfoMap[errCodeAll[j]]
			if !ok {
				ZapLog().Error("it is a BUG")
				continue
			}
			codeInfo.SetValid(0)
			if err := this.codeTable.SetCodeTable(codeInfo); err != nil {
				ZapLog().Error("SetCodeTable err", zap.Error(err), zap.Any("codeInfo", codeInfo))
				continue
			}
		}
		//goto RequestAgain
		return res, codes, fmt.Errorf("%d-%s", res.Status.Error_code, res.Status.Error_message)
	}

	if res.Status.Error_code != 400 && res.Status.Error_code != 0 {
		//  换下一个API-key   再请求
		this.NextIndex()
	}
	//429是频率限制错误

	return res, codes, fmt.Errorf("%d-%s", res.Status.Error_code, res.Status.Error_message)
}

func (this *CoinMarket) parseInValidId(inStr string) []string {
	inStrArr := strings.Split(inStr, ":")
	if len(inStrArr) < 2 {
		return nil
	}
	inStr = strings.Replace(inStrArr[1], "\"", "", -1)
	inStr = strings.Replace(inStr, " ", "", -1)
	ids := strings.Split(inStr, ",")
	return ids
}

//coinmarketcap
//func ResToResCoinMarketCap(coinId string,res *ResCoinMarketCap) *ResToResCoinMarketCap {
//	 moneyInfo := new(MoneyInfo)
//	 price := res.Data[coinId].Quote.Price
//	 moneyInfo.Price = &price
//	 symbol := res.Data[coinId].Symbol
//	 moneyInfo.Symbol = &symbol
//	 times := res.Status.Timestamp
//	 timestamp :=TimeToTimestamp(times)
//	 moneyInfo.Last_updated =&timestamp
//
//	return moneyInfo
//}

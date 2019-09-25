package v1

// 模拟充值
type ReqRecharge struct {
	Coin  string  `json:"coin" doc:"币种"`
	Token string  `json:"token" doc:"token"`
	To    string  `json:"to" doc:"充值地址"`
	Value float64 `json:"value" doc:"数量，为币种单位，！！！不是10^-8了"`
}

// 模拟挖矿
type ReqGenerate struct {
	Coin  string `json:"coin" doc:"币种，支持btc"`
	Count int    `json:"count" doc:"块数量"`
}

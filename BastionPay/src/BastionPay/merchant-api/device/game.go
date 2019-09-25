package device

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/device/game/base"
	"encoding/json"
	"go.uber.org/zap"
)

type Game struct {
	wsconn base.WsCon
	id     string
}

func (this *Game) Init(addr, id string) error {
	this.id = id
	err := this.wsconn.Init(addr, SendPingHander, RecvPingHander, RecvPongHander)
	if err != nil {
		return err
	}

	this.wsconn.Start()
	return nil
}

func (this *Game) GetId() string {
	return ""
}

func (this *Game) Send(data interface{}) error {
	message2 := []byte("{\"type\":\"admin.coinup\",\"upnumber\":\"" + data.(string) + "\",\"devid\":\"" + this.id + "\"}")
	if err := this.wsconn.Send(1, message2); err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		msg, err := this.wsconn.Recv()
		if err != nil {
			ZapLog().Error("wsconn.Recv err", zap.Error(err))
			continue
		}
		Msg := new(MsgRcv)

		err = json.Unmarshal(msg.Data, Msg)
		if err != nil {
			ZapLog().Error("unmarshal recv data err", zap.Error(err))
			continue
		}

		if Msg.Result.Res == true {
			return nil
		}
	}
	return nil
}

func SendPingHander() []byte {

	return []byte("ping")
}

func RecvPingHander(str string) []byte {

	return []byte(str)
}

func RecvPongHander(str string) []byte {

	return []byte(str)
}

type MsgRcv struct {
	Type   string `yaml:"type"`
	Result Result `yaml:"result"`
}

type Result struct {
	Name  string `yaml:"name"`
	Res   bool   `yaml:"res"`
	Stat  int64  `yaml:"stat"`
	Token int64  `yaml:"token"`
	Cmd   int64  `yaml:"cmd"`
	Coins int64  `yaml:"coins"`
	Imie  string `yaml:"imie"`
}

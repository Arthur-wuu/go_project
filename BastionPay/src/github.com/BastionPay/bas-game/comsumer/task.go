package comsumer

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-game/base"
	"BastionPay/bas-game/db"
	"BastionPay/bas-game/type"
	"encoding/json"
	"go.uber.org/zap"
	"runtime/debug"
	"strconv"
	"time"
)

const(
	 Prefix_Game_Total_Finash = "Game_AllMoney"
)

var GTasker Tasker

type  Tasker struct{
	wsconn  *base.WsCon
}

func (this *Tasker) Init(wsconn  *base.WsCon) {
	this.wsconn = wsconn
}

func (this *Tasker) Start() {
	go func(){
		for {
			this.run()
		}
	}()
}

func (this *Tasker) run(){
	defer PanicPrint()
		var finashTotal float64
		var newestTotalStr string
		newestTotal, err := this.getNewestTotal()
		if err != nil {
			ZapLog().Sugar().Errorf("getNewestTotal err[%v]", err)
			goto TASKRUNHERE
		}

		finashTotal, err = this.getFinashFromRedis()
		if err != nil {
			ZapLog().Sugar().Errorf("getFinashFromRedis err[%v]", err)
			goto TASKRUNHERE
		}
	    //fmt.Println("****newestTotal <= finashTotal",newestTotal <= finashTotal ,newestTotal,finashTotal)
		if newestTotal <= finashTotal {

			ZapLog().Sugar().Debugf("nofind new coins")
			goto TASKRUNHERE
		}
		ZapLog().Sugar().Infof("has recv coins new[%f]old[%f]", newestTotal, finashTotal)
		if err :=this.reqUpCoin(); err != nil {
			ZapLog().Sugar().Errorf("reqUpCoin err[%v]", err)
			goto TASKRUNHERE
		}
		ZapLog().Sugar().Infof("Up coins ok")

	    newestTotalStr = strconv.FormatFloat(newestTotal, 'E', -1, 64)
		this.storeFinashToRedis(newestTotalStr)
TASKRUNHERE:
		time.Sleep(time.Second * time.Duration(6))

}

func (this *Tasker) reqUpCoin() error {
	message2 := []byte("{\"type\":\"admin.coinup\",\"upnumber\":\"2\",\"devid\":\"860344040771835\"}")
	if err := this.wsconn.Send(1, message2); err != nil {
		return err
	}

	for i:=0 ;i < 3; i++{
		msg, err := this.wsconn.Recv()
		if err != nil {
			ZapLog().Error("wsconn.Recv err", zap.Error(err))
			continue
		}
		Msg := new(_type.MsgRcv)
		json.Unmarshal(msg.Data,Msg)

		if Msg.Message == "success" && Msg.State == "coinup end" && Msg.Type == "coinup"{
			return nil
		}
	}
	return nil
}

func (this *Tasker) getNewestTotal() (float64, error) {
	result := &struct{
		Total  float64
	}{}
	err :=  db.GDbMgr.Get().Table("USER_ACC").Select("sum(BALANCE) as total").Where("USER_ID=?", 35).Group("USER_ID").Scan(result).Error
	if err != nil {
		return 0,err
	}
	return result.Total,nil
}


func (this *Tasker) storeFinashToRedis (amount string) error {
	_,  err := db.GRedis.Do("SET", Prefix_Game_Total_Finash, amount)
	if err != nil {
		return err
	}
	return nil
}

func (this *Tasker) getFinashFromRedis() (float64, error ){
	allMoney,  err := db.GRedis.Do("GET", Prefix_Game_Total_Finash )
	if err != nil {
		return 0, err
	}
	sumFloat1, err := strconv.ParseFloat(string(allMoney.([]byte)), 64)
	return sumFloat1,err
}

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

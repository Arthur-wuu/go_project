package controllers

import (
	//"go.uber.org/zap"

	. "BastionPay/bas-base/log/zap"
	"BastionPay/merchant-api/baspay"
	"BastionPay/merchant-api/device"
	"BastionPay/merchant-api/models"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"go.uber.org/zap"
	"io/ioutil"
	"net/url"

	//l4g "github.com/alecthomas/log4go"
	//"github.com/asaskevich/govalidator"
	"encoding/json"
	"errors"
	"github.com/kataras/iris"
	//"BastionPay/merchant-api/base"
	//"BastionPay/bas-api/apibackend"
	//. "BastionPay/bas-base/log/zap"
	//"go.uber.org/zap"
)

//var serverPubByte = []byte(`
//-----BEGIN PUBLIC KEY-----
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArdDMsRHlTzzf8rfryGo8
//2NDQ6VntnD07ax+7CMsKAAlICv28NxLHPoWRZAl9dRhM/uWGpgOPs2sKDayilyyR
//0gZ8NPIVU4AWmn4xnv5l4Vu5HND9DcIoyvHLCiel+Lj/6HcpUzlJ+GmJ6L0QO/PI
//CPq4KyR24ggCfknzAfLi8DQ+LUGFOhiSnu1ta3z4rVeOIyy72thlGoN7aTxXSMe6
//yTi1bshkmFLgHyOcM2vpx4Vhtfb7xfu77LkRQEwi2k4vIZozInp4s5UaVFstd/Zd
//IM/hMlwKP5zv4caLhI6Op3PrG+/6McLhx3j4tRxZhc6IdfSpvzEqO7icD+oRa5Sd
//DwIDAQAB
//-----END PUBLIC KEY-----
//`)

var proServerPubByte = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhQ8CGW4m9itPq0MCmOhLpbF
yfbeiGghqRWH2x3fUzRStTLCTEOPF1gczodgESHoPIPE/MxUnqA0kJ3NnQU6kFga4CO
bza72hx0qqI5KKXAb0ApBGFZQVUJwWnVvUjGKtvB8wJ3yFewzB0TA2X3CZW1YYYlY3E
A012Alj47XQ+16EqwJuYPxIhv+u352SGMrWkAaHIld/pfZRSUzqI7Jl+4SsVTEbdQMr
ECW8g1pIG2iBtGrr7jzS6T1FgUuUy5VpNTwbpBPKn7xlI6tDgNNASnmbTeJF0LSqnBx
gMjf6xn1vLv5WBGmLQxKnG6tppbFe0xg5i7YkJDwC5VBIphgE1wIDAQAB
-----END PUBLIC KEY-----
`)

type (
	CallBacker struct {
		Controllers
	}
)

//func (this *CallBacker) TradeComplete(ctx iris.Context) {
//	//param := new(models.TradeAdd)
//	//
//	//err := Tools.ShouldBindJSON(ctx, param)
//	//if err != nil {
//	//	this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
//	//	ZapLog().Error( "param err", zap.Error(err))
//	//	return
//	//}
//	//
//	//data,err := param.Add()
//	//if err != nil {
//	//	ZapLog().Error( "Add err", zap.Error(err))
//	//	ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
//	//	return
//	//}
//
//	err := new(models.Trade).UpdateByTradeNo("", "", 1)
//	if err != nil {
//
//	}
//
//	this.Response(ctx, nil)
//
//}

//func (this *CallBacker) TradeCancel(ctx iris.Context) {
//	param := new(models.TradeAdd)
//
//	err := Tools.ShouldBindJSON(ctx, param)
//	if err != nil {
//		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "param err")
//		ZapLog().Error( "param err", zap.Error(err))
//		return
//	}
//
//	data,err := param.Add()
//	if err != nil {
//		ZapLog().Error( "Add err", zap.Error(err))
//		ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
//		return
//	}
//
//	this.Response(ctx, data)
//
//}

type NotifyParam struct {
	TradeNo         *string `json:"tradeOrderNo,omitempty"`
	Assets          *string `json:"assets,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	MerchantOrderNo *string `json:"merchantOrderNo,omitempty"`
	Status          *string `json:"status,omitempty"`
	Signature       *string `json:"signature,omitempty"`
}

func (this *CallBacker) Notify(ctx iris.Context) {

	paramMap := make(map[string][]string, 6)

	bodyBytes, err := ioutil.ReadAll(ctx.Request().Body)
	paramMap, err = url.ParseQuery(string(bodyBytes))

	ZapLog().Info("notify params:", zap.Any("paramMap", paramMap))

	amount := paramMap["amount"][0]
	assets := paramMap["assets"][0]
	status := paramMap["status"][0]
	signature := paramMap["signature"][0]
	tradeOrderNo := paramMap["tradeOrderNo"][0]
	merchantOrderNo := paramMap["merchantOrderNo"][0]

	//fmt.Println("param**",amount,assets,status)
	ZapLog().Info("order status:", zap.String("status", status))
	this.Response(ctx, "succ")

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"amount":          amount,
		"assets":          assets,
		"status":          status,
		"tradeOrderNo":    tradeOrderNo,
		"merchantOrderNo": merchantOrderNo,
	})

	signStr := baspay.RequestBodyToSignStr(reqBodySign)
	decodeSign, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		ZapLog().Error("base64 decoding err, notify end", zap.Error(err))
		return
	}
	//这里对比两个签名
	err = VerifySign(proServerPubByte, []byte(signStr), decodeSign)
	if err != nil {
		ZapLog().Error("VerifySign err, notify end", zap.Error(err))
		return
	}
	ZapLog().Info("VerifySign succ")

	if status == "3" {
		//支付成功，根据订单号去更新表的状态，然好查询机器的id，币数量
		bool, err := new(models.Trade).UpdateRowsAffected(merchantOrderNo, models.EUM_TRANSFER_STATUS_SUCCESS)
		if err != nil {
			ZapLog().Error(" update by tradeNo err, notify end", zap.Error(err))
			ctx.JSON(Response{Code: 10010, Message: err.Error()})
			return
		}
		if bool == false {
			ZapLog().Error("update 0 rows, notify end", zap.Error(err))
			ctx.JSON(Response{Code: 10010, Message: "too many notifys"})
			return
		}
		//game := new(device.Game)
		//查数据库，把订单号对应的机器id放进来就可以了
		deviceId, gameCoin, err := new(models.Trade).GetDeviceId(nil, merchantOrderNo)

		if gameCoin == "" || len(gameCoin) == 0 {
			ZapLog().Info("game coin nill")
			return
		}
		ZapLog().Info("deviceId:", zap.Any("deviceId", deviceId))

		if err != nil {
			ZapLog().Error("get device id err, notify end", zap.Error(err))
			return
		}

		//获得机器id对应的机器
		device.GDeviceMgr.Init()

		devices := device.GDeviceMgr.Get(deviceId)
		err = devices.Send(gameCoin)
		if err != nil {
			ZapLog().Error("send coin to machine fail, notify end", zap.Error(err))
			//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
			return
		}
		ZapLog().Info("send coin finish···")
		//this.Response(ctx, "succ")
		return

	} else {
		//更新表里的状态
		err := new(models.Trade).UpdateTransferStatusByTradeNo(nil, merchantOrderNo, models.EUM_TRANSFER_STATUS_FAIL)
		if err != nil {
			ZapLog().Error("fail update by tradeNo err, notify end", zap.Error(err))
			//ctx.JSON(Response{Code: apibackend.BASERR_DATABASE_ERROR.Code(), Message: err.Error()})
			return
		}
	}

	//可用chan，redis来存储，方式消息处理过慢
	{ // 退款成功，退款失败。是不是 有退款限制啊？？？
	}

	{ // 提币成功，提币失败
	}

	this.Response(ctx, nil)

}

func VerifySign(signingPubKey, signStr, sign []byte) error {
	block, _ := pem.Decode(signingPubKey)
	if block == nil {
		ZapLog().Error("block err")
		return errors.New("block err")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		ZapLog().Error("pub key err", zap.Error(err))
		return err
	}

	rsaPubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		ZapLog().Error("pubkey get err", zap.Error(err))
		return errors.New("pubkey get err")
	}

	c := crypto.Hash.New(crypto.SHA1)
	c.Write(signStr)
	digest := c.Sum(nil)

	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA1, digest, sign)
	if err != nil {
		ZapLog().Error("send coin to machine fail, notify end", zap.Error(err))
		return err
	}
	return nil
}

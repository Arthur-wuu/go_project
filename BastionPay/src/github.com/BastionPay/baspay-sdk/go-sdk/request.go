package go_sdk

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/baspay-sdk/base"
	"BastionPay/baspay-sdk/config"
	"BastionPay/baspay-sdk/util"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

//转账请求
func (this *TransferParam) Send() (string, error) {
	signType := "RSA"
	this.SignType = &signType

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"amount":               this.Amount,
		"assets":               this.Assets,
		"payee_id":             this.PayeeId,
		"product_name":         this.ProductName,
		"merchant_transfer_no": this.MerchantTransNo,
		"merchant_id":          this.MerchantId,
		"timestamp":            this.Timestamp,
		"notify_url":           this.NotifyUrl,
	})

	signStr := RequestBodyToSignStr(reqBodySign)
	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"amount":               this.Amount,
		"assets":               this.Assets,
		"payee_id":             this.PayeeId,
		"merchant_transfer_no": this.MerchantTransNo,
		"merchant_id":          this.MerchantId,
		"product_name":         this.ProductName,
		"timestamp":            this.Timestamp,
		"notify_url":           this.NotifyUrl,
		"sign_type":            this.SignType,
		"signature":            finalSign,
	})
	//fmt.Println("**reqBody**",string(reqBody))
	url := "https://open-api.bastionpay.com/open-api/trade/transfer" //"https://open-api.bastionpay.com/"open-api/trade/transfer
	//if config.GConfig.EnvFlag.Flag == "pro" {
	//	url = commonUrl + transUrl
	//}else {
	//	url = testCommonUrl+ transUrl
	//}
	//fmt.Println("url:", url )
	result, err := base.HttpSend(url, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		fmt.Println("err", err)
		//ZapLog().Error( "request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	transRes := new(TransRes)
	err = json.Unmarshal(result, transRes)
	fmt.Println("err:", err)
	fmt.Println("status:", transRes.Data.Status)
	status := strconv.Itoa(transRes.Data.Status)

	return status, nil
}

//有效币种请求
func (this *AvailAssetsParam) Send() (interface{}, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"assets": this.Assets,
	})

	url := ""
	if config.GConfig.EnvFlag.Flag == "pro" {
		url = commonUrl + availAssetsUrl
	} else {
		url = testCommonUrl + availAssetsUrl
	}

	result, err := base.HttpSend(url, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err", zap.Error(err))
		return "", err
	}
	fmt.Println("**result**", string(result))

	transRes := new(AvailAssetsRes)
	err = json.Unmarshal(result, transRes)

	fmt.Println("array:", transRes.Data.Array)

	return transRes.Data.Array, nil
}

//创建二维码订单
func (this *QrTrade) Send() (interface{}, error) {
	signType := "RSA"
	expireTime := strconv.FormatInt(this.ExpireTime, 10)

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址
		"merchant_id":       this.MerchantId,      //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature,	     //⽤户请求的签名串
		"timestamp":  this.Timestamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	fmt.Println("**reqBody sign**", string(reqBodySign))

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   signType,        //签名算法类型
		"signature":   finalSign,       //⽤户请求的签名串
		"timestamp":   this.Timestamp,  //发送请求的时间
		"notify_url":  this.NotifyUrl,  //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	url := ""
	if config.GConfig.EnvFlag.Flag == "pro" {
		url = commonUrl + createQrUrl
	} else {
		url = testCommonUrl + createQrUrl
	}

	result, err := base.HttpSend(url, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	qrRes := new(TradeQrRes)
	err = json.Unmarshal(result, qrRes)

	fmt.Println("qrRes:", qrRes)

	return qrRes, nil
}

//创建Sdk订单
func (this *SdkTrade) Send() (interface{}, error) {
	signType := "RSA"
	expireTime := strconv.FormatInt(this.ExpireTime, 10)

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址
		"merchant_id":       this.MerchantId,      //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature,	     //⽤户请求的签名串
		"timestamp":  this.Timestamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	fmt.Println("**reqBody sign**", string(reqBodySign))

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   signType,        //签名算法类型
		"signature":   finalSign,       //⽤户请求的签名串
		"timestamp":   this.Timestamp,  //发送请求的时间
		"notify_url":  this.NotifyUrl,  //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(createSdkUrl, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	qrRes := new(TradeQrRes)
	err = json.Unmarshal(result, qrRes)

	fmt.Println("qrRes:", qrRes)

	return qrRes, nil
}

//创建wap订单
func (this *WapTrade) Send() (interface{}, error) {
	signType := "RSA"
	expireTime := strconv.FormatInt(this.ExpireTime, 10)

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址
		"merchant_id":       this.MerchantId,      //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature,	     //⽤户请求的签名串
		"timestamp":  this.Timestamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	fmt.Println("**reqBody sign**", string(reqBodySign))

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"payee_id":          this.PayeeId,         //商户BastionPay 用户Id
		"assets":            this.Assets,          //数字货币币种
		"amount":            this.Amount,          //数字货币数量
		"product_name":      this.ProductName,     //订单标题
		"product_detail":    this.ProductDetail,   //订单描述
		"expire_time":       expireTime,           //订单超时时间
		"remark":            this.Remark,          //交易易备注信息
		"return_url":        this.ReturnUrl,       //⽀付完成后的回调地址
		"show_url":          this.ShowUrl,         //取消⽀付的回调地址

		"merchant_id": this.MerchantId, //商户在BabstionPay 注册的商户ID
		"sign_type":   signType,        //签名算法类型
		"signature":   finalSign,       //⽤户请求的签名串
		"timestamp":   this.Timestamp,  //发送请求的时间
		"notify_url":  this.NotifyUrl,  //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(createWapUrl, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	qrRes := new(TradeQrRes)
	err = json.Unmarshal(result, qrRes)

	fmt.Println("qrRes:", qrRes)

	return qrRes, nil
}

//查询交易信息
func (this *TradeInfo) Send() (interface{}, error) {
	signType := "RSA"

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"merchant_id":       this.MerchantId,      //商户在BabstionPay 注册的商户ID
		"notify_url":        this.NotifyUrl,       //回调通知商户服务器器的地址
	})

	fmt.Println("**reqBody sign**", string(reqBodySign))

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_trade_no": this.MerchantTradeNo, //商户的订单号 UUID 生成的
		"merchant_id":       this.MerchantId,      //商户在BabstionPay 注册的商户ID
		"sign_type":         signType,             //签名算法类型
		"signature":         finalSign,            //⽤户请求的签名串
		"notify_url":        this.NotifyUrl,       //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(tradeInfoUrl, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	infoRes := new(InfoRes)
	err = json.Unmarshal(result, infoRes)

	fmt.Println("infoRes:", infoRes)

	return infoRes, nil
}

//创建pos订单
func (this *PosTrade) Send() (interface{}, error) {
	signType := "RSA"

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"payee_id":        this.PayeeId,       //商户BastionPay 用户Id
		"assets":          this.Assets,        //数字货币币种
		"amount":          this.Amount,        //数字货币数量
		"product_name":    this.ProductName,   //订单标题
		"product_detail":  this.ProductDetail, //订单描述
		"pay_voucher":     this.PayVoucher,    //订单超时时间
		"remark":          this.Remark,        //交易易备注信息
		"pos_machine_id":  this.PosMachineId,  //⽀付完成后的回调地址
		"merchant_id":     this.MerchantId,    //商户在BabstionPay 注册的商户ID
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature,	     //⽤户请求的签名串
		"timestamp":  this.Timestamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	fmt.Println("**reqBody sign**", string(reqBodySign))

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"payee_id":        this.PayeeId,       //商户BastionPay 用户Id
		"assets":          this.Assets,        //数字货币币种
		"amount":          this.Amount,        //数字货币数量
		"product_name":    this.ProductName,   //订单标题
		"product_detail":  this.ProductDetail, //订单描述
		"pay_voucher":     this.PayVoucher,    //订单超时时间
		"remark":          this.Remark,        //交易易备注信息
		"pos_machine_id":  this.PosMachineId,  //⽀付完成后的回调地址
		"merchant_id":     this.MerchantId,    //商户在BabstionPay 注册的商户ID
		"sign_type":       signType,           //签名算法类型
		"signature":       finalSign,          //⽤户请求的签名串
		"timestamp":       this.Timestamp,     //发送请求的时间
		"notify_url":      this.NotifyUrl,     //回调通知商户服务器器的地址
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(posUrl, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	posRes := new(PosRes)
	err = json.Unmarshal(result, posRes)

	fmt.Println("posRes:", posRes)

	return posRes, nil
}

//查询pos订单列表
func (this *PosOrderList) Send() (interface{}, error) {
	signType := "RSA"

	reqBodySign, _ := json.Marshal(map[string]interface{}{
		"pos_machine_id":  this.PosMachineId,
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"page":            this.Page,          //商户BastionPay 用户Id
		"page_size":       this.PageSize,      //数字货币币种
		"merchant_id":     this.MerchantId,    //商户在BabstionPay 注册的商户ID
		"begin_time":      this.BeginTime,     //订单标题
		"end_time":        this.EndTime,       //订单描述
		"trade_status":    this.TradeStatus,   //交易易备注信息
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature  	     //⽤户请求的签名串
		"timestamp":  this.TimeStamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
	})

	fmt.Println("**reqBody sign**", string(reqBodySign))

	signStr := RequestBodyToSignStr(reqBodySign)
	fmt.Println("**signStr", signStr)

	//签名 SHA1
	sha1 := new(utils.SHAwithRSA)
	sha1.SetPriKey(h5PrivateKey)
	finalSign, err := sha1.Sign(signStr)

	if err != nil {
		ZapLog().Error("sign err", zap.Error(err))
		return "", err
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"pos_machine_id":  this.PosMachineId,
		"merchant_pos_no": this.MerchantPosNo, //商户的订单号 UUID 生成的
		"page":            this.Page,          //商户BastionPay 用户Id
		"page_size":       this.PageSize,      //数字货币币种
		"merchant_id":     this.MerchantId,    //商户在BabstionPay 注册的商户ID
		"begin_time":      this.BeginTime,     //订单标题
		"end_time":        this.EndTime,       //订单描述
		"trade_status":    this.TradeStatus,   //交易易备注信息
		//"sign_type": this.Request.SignType,	     //签名算法类型
		//"signature": this.Request.Signature  	     //⽤户请求的签名串
		"timestamp":  this.TimeStamp, //发送请求的时间
		"notify_url": this.NotifyUrl, //回调通知商户服务器器的地址
		"sign_type":  signType,       //签名算法类型
		"signature":  finalSign,      //⽤户请求的签名串
	})
	fmt.Println("**reqBody**", string(reqBody))

	result, err := base.HttpSend(posOrderUrl, bytes.NewBuffer(reqBody), "POST", nil)
	if err != nil {
		ZapLog().Error("request err")
		return "", err
	}

	fmt.Println("**result**", string(result))

	posRes := new(PosListRes)
	err = json.Unmarshal(result, posRes)

	fmt.Println("posRes:", posRes)

	return posRes, nil
}

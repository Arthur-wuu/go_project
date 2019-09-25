package bastionpay

import (
	"testing"
)

const (
	apiKey              = "1a627ca4-66d4-460a-a344-4c5f9021a9b7"
	privateKey          = "./cert/user_private.pem"
	publicKey           = "./cert/user_public.pem"
	bastionPayPublicKey = "./cert/bestionpay_public.pem"
)

func TestSession(t *testing.T) {
	sess, err := New(NewConfig(apiKey,
		NewCredentials(privateKey, publicKey, bastionPayPublicKey)))
	if err != nil {
		t.Fatal(err.Error())
	}

	// 获取支持列表
	t.Log("start get supports...")
	symbols, err := sess.GetSupports()
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(symbols)
	t.Log("end get supports...\n")

	// 获取币种资料
	t.Log("start get currency info...")
	currency, err := sess.GetCurrencysInfo([]string{"eth", "btc"})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(currency)
	t.Log("end get currency info...\n")

	// 获取地址
	t.Log("start get address...")
	addresses, err := sess.GetAddress("eth", 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(addresses)
	t.Log("end get address...\n")

	// 块高查询
	t.Log("start withdrawal...")
	height, err := sess.GetHeight([]string{"eth"})
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(height)
	t.Log("end withdrawal...\n")

	// 查询余额
	t.Log("start get balance...")
	balance, err := sess.GetBalance([]string{"eth"})
	if err != nil {
		t.Log(err.Error())
	}
	t.Log(balance)
	t.Log("end get balance...\n")

	// 提币
	t.Log("start withdrawal...")
	orderId, err := sess.Withdrawal("235", "eth", "0xhaha", 30.834)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(orderId)
	t.Log("end withdrawal...\n")

	// 获取订单列表
	t.Log("start get order history...")
	history, err := sess.OrderHistory("eth", 0, 0, 0, 159087093048)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(history)
	t.Log("end get order history...\n")

	// 获取通知
	t.Log("start get order history...")
	msg, err := sess.GetMsgs(1, 2)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(msg)
	t.Log("end get order history...\n")

	// 测试充值回调
	// 从channel中取回回调数据
	//go func() {
	//	for {
	//		select {
	//		case msg := <-sess.MsgChannel:
	//			t.Log("recive callback message...")
	//			t.Log(msg)
	//		}
	//	}
	//}()
	//// 设置最后一次处理的ID, 通知ID如果小于此ID将会被忽略
	//sess.SetLastMsgID(1)
	//// 监听10000端口  POST :10000/callback
	//sess.EnableMonitoring("10000")

}

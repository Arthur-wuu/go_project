package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	"BastionPay/merchant-api/controllers"
)

func (this *WebServer) routes() {
	app := this.mIris
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(
			map[string]interface{}{
				"code": 0,
			})
	})

	app.Any("/bk_all_url", func(ctx iris.Context) {
		ctx.JSON([]string{
			//"/v2/bas-merchant/trade/create-web",
			//"/v2/bas-merchant/trade/create-qr",
			//"/v2/bas-merchant/trade/create-sdk",
			//"/v2/bas-merchant/callback/notify",

			"/v2/bas-merchant-bk/config-device/add",
			"/v2/bas-merchant-bk/config-device/update",
			"/v2/bas-merchant-bk/config-device/get",
			"/v2/bas-merchant-bk/config-device/del",
			"/v2/bas-merchant-bk/config-device/list",
			"/v2/bas-merchant-bk/config-device/coinlist",
			"/v2/bas-merchant-bk/config-device/price",
		})
	})

	v1 := app.Party("/v2", crs)
	{
		bas := v1.Party("/bas-merchant")
		{
			tradeParty := bas.Party("/trade")
			{
				trade := new(controllers.Trade)

				tradeParty.Post("/create-web", trade.CreateWeb)
				tradeParty.Post("/create-qr", trade.CreateQr)
				tradeParty.Post("/create-sdk", trade.CreateSdk)
				tradeParty.Post("/list-trade", trade.ListTrade)
				tradeParty.Post("/info", trade.TradeInfo)
				tradeParty.Post("/pay", trade.Pay)
				tradeParty.Post("/cancel", nil)
				tradeParty.Post("/refund", trade.ReFund)
				tradeParty.Post("/list-refund", trade.ReFundList)
			}
			//处理form
			{
				form := new(controllers.Form)
				tradeParty.Post("/form", form.FormAction)
			}

			// 创建咖啡订单
			{
				coffee := new(controllers.CoffeeTrade)
				tradeParty.Post("/coffee", coffee.CoffeeQr)
			}
			// 轮训咖啡订单
			{
				coffee := new(controllers.CoffeeTrade)
				tradeParty.Post("/order-status", coffee.OrderStatus)
			}

			// pos机
			{
				pos := new(controllers.Pos)
				tradeParty.Post("/pos", pos.CreatePos)
				tradeParty.Post("/pos_orders", pos.PosOrders)
				tradeParty.Post("/pos_list", pos.List)
			}
			//查询支持的币种信息
			{
				assets := new(controllers.Assets)
				tradeParty.Post("/avali-assets", assets.AvaliAssets)
			}

			//atm数字货币到法币的quote
			{
				assets := new(controllers.Quote)
				tradeParty.Post("/quote", assets.GetQuote)
			}

			//atm存钱转到相应的数字货币数量
			{
				amount := new(controllers.Amount)
				tradeParty.Post("/get-amount", amount.GetAmount)
			}

			fundParty := bas.Party("/fundout")
			{
				trade := new(controllers.FundOut)

				fundParty.Post("/create", trade.Create)
				fundParty.Post("/list", trade.List)
			}

			callbackParty := bas.Party("/callback")
			{
				trade := new(controllers.CallBacker)
				callbackParty.Post("/notify", trade.Notify)
			}

			// 转账接口
			atmTransfer := bas.Party("/atm")
			{
				transfer := new(controllers.Transfer)
				atmTransfer.Post("/transfer", transfer.TransferAtmOpenApi)  //待改
			}

			//有价格折扣手续费的 现在有atm
			feeQr := bas.Party("/fee")
			{
				atm := new(controllers.Atm)
				feeQr.Post("/atm_qr", atm.CreateAtmQr)
			}
		}

		//把之前按照设备分的接口 现在按照前后端分
		configPartyBk := bas.Party("/config-device")
		{
			config := new(controllers.BkConfig)

			configPartyBk.Post("/add", config.Add)
			configPartyBk.Post("/update", config.Update)
			configPartyBk.Post("/get", config.Get)
			configPartyBk.Post("/del", config.Delete)
			configPartyBk.Post("/list", config.List)
			configPartyBk.Post("/coinlist", config.GetCoinList)
		}
		{
			param := new(controllers.PrePareParam)
			configPartyBk.Post("/price", param.GetPrice)
		}

		//url白名单，钱包可以扫的二维码
		{
			param := new(controllers.WhiteList)
			configPartyBk.Post("/get-url", param.Get)
		}



//v2/bas-merchant-bk/config-device/list    参数 {"page":1}


		//后台配置
		basbk := v1.Party("/bas-merchant-bk")
		configParty := basbk.Party("/config-device")
		{
			config := new(controllers.BkConfig)

			configParty.Post("/add", config.Add)
			configParty.Post("/update", config.Update)
			configParty.Post("/get", config.Get)
			configParty.Post("/del", config.Delete)
			configParty.Post("/list", config.List)
			configParty.Post("/coinlist", config.GetCoinList)
		}
		//{
		//	param := new(controllers.PrePareParam)
		//	configParty.Post("/price", param.GetPrice)
		//}
		//
		////url白名单，钱包可以扫的二维码
		//{
		//	param := new(controllers.WhiteList)
		//	configParty.Post("/get-url", param.Get)
		//}


		cigMerchantParty := basbk.Party("/config-merchant")
		{
			mConfig := new(controllers.MerchantConfig)

			cigMerchantParty.Post("/add", mConfig.Add)
			cigMerchantParty.Post("/update", mConfig.Update)
			cigMerchantParty.Post("/get", mConfig.Get)
			cigMerchantParty.Post("/del", mConfig.Delete)
			cigMerchantParty.Post("/list", mConfig.List)
		}
	}



}
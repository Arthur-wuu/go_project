package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	"BastionPay/pay-user-merchant-api/controllers"
)

func (this *WebServer) bkroutes()  {
	app := this.mBkIris
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language", "Api-Key", "Req-Real-Ip"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(
			map[string]interface{}{
				"code": 0,
			})
	})

	v1InnerUser := app.Party("/v1/bk/usermerchant", crs)
	{
		userBkParty := v1InnerUser.Party("/login", crs)
		{
			userCtrl     := controllers.NewUserController()

			//accountBkParty.Get("/getuserinfo", infoCtrl.GetUserInfo)
			userBkParty.Post("/qr_callback", userCtrl.BkLoginQrCallBack)
		}
	}
}
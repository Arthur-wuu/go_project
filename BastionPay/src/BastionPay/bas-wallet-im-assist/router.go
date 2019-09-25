package main

import (
	"BastionPay/bas-wallet-im-assist/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
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

	v1 := app.Party("/v1", crs)
	{
		bas := v1.Party("/wallet/imassist")
		{

			callbackParty := bas.Party("/callback")
			{
				callBack := new(controllers.CallBacker)

				callbackParty.Post("/notify", callBack.Handle)
			}
		}
	}

}

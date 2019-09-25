package main

import (
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	"BastionPay/pay-user-merchant-api/controllers"
)

func (this *WebServer) routes() {
	app := this.mIris
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

	interceptor := new(controllers.Interceptor)

	v1 := app.Party("/v1/ft/usermerchant", crs, interceptor.VerifyAccess)
	{
		v1.Any("/", func(ctx iris.Context) {
			ctx.JSON(
				map[string]interface{}{
					"code": 0,
				})
		})
		{
			logCtrl := new(controllers.LogController)

			v1.Done(logCtrl.RecodeLog)
			v1.Get("/log/login", logCtrl.GetLoginLog)
			v1.Get("/log/safe", logCtrl.GetOperationLog)
		}
		{
			verificationCtrl := controllers.NewVerificationController()

			v1.Get("/verification/{type}", verificationCtrl.Send)
			v1.Post("/verification", verificationCtrl.Verification)
		}
		{
			userCtrl   := controllers.NewUserController()

			v1.Get("/refresh", userCtrl.RefreshToken)
			v1.Post("/login", userCtrl.Login)
			v1.Post("/login/ga", userCtrl.LoginWithGa)
			v1.Post("/login/qr_gen", userCtrl.LoginQr)
			v1.Post("/login/qr_check", userCtrl.LoginQrCheck)
			//v1.Post("/register", userCtrl.Register)
			//v1.Post("/exists", userCtrl.Exists)
		}
		{
			//gaCtrl  := controllers.NewGaController()
			//
			//v1.Get("/ga", gaCtrl.Generate)
			//v1.Post("/ga/bind", gaCtrl.Bind)
			//v1.Post("/ga/unbind", gaCtrl.UnBind)
		}
		{
			//passwordCtrl := controllers.NewPasswordController()

			//v1.Post("/password/modify", passwordCtrl.Modify)
			//v1.Post("/password/inquire", passwordCtrl.Inquire)
			//v1.Post("/password/reset", passwordCtrl.Reset)
		}
		{
			infoCtrl     := controllers.NewInfoController()

			v1.Get("/info", infoCtrl.GetInformation)
			v1.Get("/info/nohide", infoCtrl.GetInformationNoHide)
			v1.Post("/info", infoCtrl.SetInformation)
			//v1.Post("/info/email", infoCtrl.BindEmail)
			//v1.Post("/info/phone", infoCtrl.BindPhone)
			//v1.Post("/info/phone/rebind", infoCtrl.RebindPhone)
		}
		{
			//optionalCtrl  := controllers.NewOptionalController()
			//
			//v1.Get("/optional", optionalCtrl.Get)
			//v1.Post("/optional", optionalCtrl.Update)
		}
		merPty := v1.Party("/merchant")
		{
			merchantCtrl := new(controllers.Merchant)

			//merPty.Post("/create", merchantCtrl.Create)
			merPty.Post("/get", merchantCtrl.Get)
			merPty.Post("/add", merchantCtrl.Add)
			merPty.Post("/update", merchantCtrl.Update)
			//merPty.Post("/save", merchantCtrl.Update)
		}


		orderPty := v1.Party("/order")
		{
			orderCtrl := new(controllers.Trade)
			orderPty.Post("/trade-list", orderCtrl.ListTrade)
			orderPty.Post("/refund-list", orderCtrl.ReFundList)

		}
	}
}


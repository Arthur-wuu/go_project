package main

import (
	"BastionPay/baspay-recharge/controllers"
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

	//app.Post("/bk_all_url", func(ctx iris.Context) {
	//	ctx.JSON([]string{
	//		"/v1/user-help-bk/message/add",
	//		"/v1/user-help-bk/message/list",
	//		"/v1/user-help-bk/message/update",
	//		"/v1/app-version-bk/whitelabel/add",
	//		"/v1/app-version-bk/whitelabel/update",
	//		"/v1/app-version-bk/whitelabel/list",
	//		"/v1/app-version-bk/version/add",
	//		"/v1/app-version-bk/version/update",
	//		"/v1/app-version-bk/version/list",
	//		"/v1/app-version-bk/file/upload",
	//	})
	//})
	//用户充值
	v1 := app.Party("/v1", crs)
	{
		userHelpBk := v1.Party("/fissionshare")
		{
			userHelpBkParty := userHelpBk.Party("/callback")
			{
				notify := controllers.Notify{}
				userHelpBkParty.Post("/notify", notify.TransferCallBack)

				rubbishNotify := controllers.RubbishNotify{}
				userHelpBkParty.Post("/rubbish_notify", rubbishNotify.TransferCallBack)

				commonNotify := controllers.CommonNotify{}
				userHelpBkParty.Post("/common-notify", commonNotify.TransferCallBackCommon)
				//userHelpBkParty.Post("/list", user.List)
				//userHelpBkParty.Post("/update", user.Update)
			}
			{
				//
			}
		}
	}

}

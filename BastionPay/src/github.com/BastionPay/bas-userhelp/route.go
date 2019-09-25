package main

import (
	"BastionPay/bas-userhelp/controllers"
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

	app.Post("/bk_all_url", func(ctx iris.Context) {
		ctx.JSON([]string{
			"/v1/user-help-bk/message/add",
			"/v1/user-help-bk/message/list",
			"/v1/user-help-bk/message/update",
			"/v1/app-version-bk/whitelabel/add",
			"/v1/app-version-bk/whitelabel/update",
			"/v1/app-version-bk/whitelabel/list",
			"/v1/app-version-bk/version/add",
			"/v1/app-version-bk/version/update",
			"/v1/app-version-bk/version/list",
			"/v1/app-version-bk/file/upload",
		})
	})
//用户帮助信息
	v1 := app.Party("/v1", crs)
	{
		userHelpBk := v1.Party("/user-help-bk")
		{
			userHelpBkParty := userHelpBk.Party("/message")
			{
				user := controllers.UserHelp{}

				userHelpBkParty.Post("/add", user.Add)
				userHelpBkParty.Post("/list", user.List)
				userHelpBkParty.Post("/update", user.Update)
			}
		}
	}
//app版本信息
	{
		appVersionBk := v1.Party("/app-version-bk")
		{
			wBk := appVersionBk.Party("/whitelabel")
			{
				upload := controllers.AppWhiteLabel{}
				wBk.Post("/add",  upload.Add)
				wBk.Post("/update",  upload.Update)
				wBk.Post("/list",  upload.List)

			}
			fBk := appVersionBk.Party("/file")
			{
				upload := controllers.UploadFile{}
				fBk.Post("/upload", upload.HandlePicFiles)
			}
			verBk :=  appVersionBk.Party("/version")
			{
				appVersion := controllers.AppVersion{}
				verBk.Post("/add",  appVersion.Add)
				verBk.Post("/list", appVersion.List)
				verBk.Post("/update", appVersion.Update)
			}
		}
	}
	{
		appVersionBk := v1.Party("/app-version")
		{
			appVersion := controllers.AppVersion{}
			appVersionBk.Post("/get",  appVersion.GetForFront)
		}
	}

}

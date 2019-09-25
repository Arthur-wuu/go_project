package main

import (
	"BastionPay/bas-wallet-fission-share/controllers"
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

	v1 := app.Party("/v1/fissionshare", crs)
	{
		//活动，添加，list，更新，警用 全是后端的
		acParty := v1.Party("/activity")
		{
			ac := new(controllers.Activity)

			acParty.Post("/list", ac.ListForFront)
		}
		//大红包，创建，list（back）
		redParty := v1.Party("/red")
		{
			viplv := new(controllers.Red)

			redParty.Post("/add", viplv.Add)
			redParty.Post("/list", viplv.ListForFront)

		}
		//小红包，创建，list（back）
		robParty := v1.Party("/robber")
		{
			viplv := new(controllers.Robber)

			robParty.Post("/add", viplv.Rob)
			robParty.Post("/list", viplv.ListForFront)
		}

		sloganParty := v1.Party("/slogan")
		{
			slogan := new(controllers.Slogan)

			//sloganParty.Post("/add", slogan.Add)
			sloganParty.Post("/list", slogan.ListForFront)
			sloganParty.Post("/getall", slogan.GetAll)
		}

	}

	v1bk := app.Party("/v1/bk/fissionshare", crs)
	{
		//活动，添加，list，更新，警用 全是后端的
		acParty := v1bk.Party("/activity")
		{
			ac := new(controllers.Activity)

			acParty.Post("/add", ac.Add)
			acParty.Post("/list", ac.ListForBack)
			acParty.Post("/update", ac.Update)
		}
		//大红包，创建，list（back）
		redParty := v1bk.Party("/red")
		{
			viplv := new(controllers.Red)

			//redParty.Post("/add", viplv.Add)
			redParty.Post("/list", viplv.ListforBack)

		}
		//小红包，创建，list（back）
		robParty := v1bk.Party("/robber")
		{
			viplv := new(controllers.Robber)

			//robParty.Post("/add", viplv.Rob)
			robParty.Post("/list", viplv.ListForBack)
		}

		sloganParty := v1bk.Party("/slogan")
		{
			slogan := new(controllers.Slogan)

			sloganParty.Post("/add", slogan.Add)
			sloganParty.Post("/list", slogan.ListForBack)
			sloganParty.Post("/update", slogan.Update)
		}

	}

}

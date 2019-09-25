package main

import (
	"BastionPay/merchant-workattence-api/controllers"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
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

	v1bk := app.Party("/v1/bk", crs)
	{
		//账户关联
		acParty := v1bk.Party("/account")
		{
			ac := new(controllers.AccountMap)

			acParty.Post("/add", ac.AddForBack)
			acParty.Post("/list", ac.ListForBack)
			acParty.Post("/update", ac.UpdateForBack)
			acParty.Post("/get", ac.GetAccount)
		}

		//垃圾分类相关
		rcPart := v1bk.Party("/rubbishclassify")
		{
			rc := new(controllers.RubbishClassify)
			rcPart.Post("/send", rc.SendAward)
			rcPart.Post("/list", rc.ListForBack)
			rcPart.Post("/preadd", rc.DepartmentListForBack)
			rcPart.Post("/add", rc.AddForBack)
			rcPart.Post("/recordlist", rc.RecordListForBack)
			rcPart.Post("/edit", rc.GetRecordForBackEdit)
			rcPart.Post("/update", rc.RecordUpdateForBack)
			rcPart.Post("/detail", rc.SendDetailForBack)
		}

		//工龄奖励
		smPart := v1bk.Party("/staffmotivation")
		{
			sm := new(controllers.StaffMotivation)
			smPart.Post("/day", sm.DayListForBack)
			smPart.Post("/total", sm.TotalListForBack)
		}

		//考勤
		waPart := v1bk.Party("/workattence")
		{
			wa := new(controllers.WorkAttendance)
			waPart.Post("/list", wa.OvertimeAwardListForBack)
		}
	}

	v1ft := app.Party("/api/data", crs)
	{
		ac := new(controllers.Receive)

		v1ft.Post("/post", ac.CheckinRecord)
	}

}

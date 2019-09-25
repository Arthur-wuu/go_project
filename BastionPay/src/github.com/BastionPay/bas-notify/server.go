package main

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/config"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/pprof"
	"github.com/kataras/iris/middleware/recover"
	"go.uber.org/zap"
	"BastionPay/bas-notify/controllers"
	"BastionPay/bas-notify/sms"
	"BastionPay/bas-notify/email"
	"BastionPay/bas-notify/db"
	"BastionPay/bas-notify/models"
)
//
//const (
//	ErrCode_Success    = 0
//	ErrCode_Param      = 10001   //参数问题
//	ErrCode_InerServer = 10002   //内部错误
//	ErrCode_Exist      = 10003   //模板名已存在错误
//	ErrCode_NoAlive    = 10004   //模板未激活
//	ErrCode_NoFindAliveTemplate = 10005 //获取模板的函数未定义
//	ErrCode_SubRecipient = 10006 //多收件人中，有一个收件人报错
//	ErrCode_NoRecipient = 10007 //没有任何收件人
//	ErrCode_TemplateParse = 10008 //模板无法解析
//)

type WebServer struct {
	mIris        *iris.Application
	//mTemplateMgr TemplateMgr
	//mSmsMgr      sms.GSmsMgr
	//mMailMgr     MailMgr
}

func NewWebServer() *WebServer {
	web := new(WebServer)
	if err := web.Init(); err != nil {
		ZapLog().With(zap.Error(err)).Error("WebServer Init err")
		panic("WebServer Init err")
	}
	return web
}

func (this *WebServer) Init() error {
	//err := this.mTemplateMgr.Init()
	//if err != nil {
	//	ZapLog().With(zap.Error(err)).Error("TemplateMgr Init err")
	//	return err
	//}
	if err := sms.GSmsMgr.Init(); err != nil {
		ZapLog().With(zap.Error(err)).Error("SmsMgr Init err")
		return err
	}
	if err := email.GMailMgr.Init(); err != nil {
		ZapLog().With(zap.Error(err)).Error("MailMgr Init err")
		return err
	}
	dbOp := &db.DbOptions{
		Host:        config.GConfig.Db.Host,
		Port:        config.GConfig.Db.Port,
		User:        config.GConfig.Db.User,
		Pass:        config.GConfig.Db.Pwd,
		DbName:      config.GConfig.Db.Quote_db,
		MaxIdleConn: config.GConfig.Db.Max_idle_conn,
		MaxOpenConn: config.GConfig.Db.Max_open_conn,
	}
	ZapLog().Sugar().Infof("%v", *dbOp)
	err := db.GDbMgr.Init(dbOp)
	if err != nil {
		ZapLog().With(zap.Error(err)).Error("Db Init err")
		return err
	}
	models.InitDbTable()
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	if config.GConfig.Server.Debug {
		app.Any("/debug/pprof/{action:path}", pprof.New())
	}
	this.mIris = app
	this.controller()
	ZapLog().Info("WebServer Init ok")
	return nil
}

func (this *WebServer) Run() error {
	//if err := this.mTemplateMgr.Start(); err != nil {
	//	ZapLog().With(zap.Error(err)).Error("TemplateMgr Start err")
	//	return err
	//}
	err := this.mIris.Run(iris.Addr(":" + config.GConfig.Server.Port)) //阻塞模式
	if err != nil {
		if err == iris.ErrServerClosed {
			ZapLog().Sugar().Infof("Iris Run[%d] Stoped[%v]", config.GConfig.Server.Port, err)
		} else {
			ZapLog().Sugar().Errorf("Iris Run[%d] err[%v]", config.GConfig.Server.Port, err)
		}
	}
	return nil
}

func (this *WebServer) Stop() error {
	return nil
}

/********************内部接口************************/
func (a *WebServer) controller() {
	app := a.mIris
	//app.UseGlobal(interceptorCtrl.Interceptor)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "X-Requested-With", "X_Requested_With", "Content-Type", "Access-Token", "Accept-Language"},
		AllowCredentials: true,
	})

	app.Any("/", func(ctx iris.Context) {
		ctx.JSON(map[string]interface{}{
			"code":    0,
			"message": "ok",
			"data":    "",
		})
	})

	//app.Any("/sms/test" , a.smsSendTest)

	v1 := app.Party("/v1/notify", crs, func(ctx iris.Context) { ctx.Next() }).AllowMethods(iris.MethodOptions)
	{

		tmpGroupPy := v1.Party("/templategroup")
		{
			group := controllers.TemplateGroup{}
			tmpGroupPy.Post("/getall", group.List)
			tmpGroupPy.Post("/add", group.Add)
			tmpGroupPy.Post("/update", group.Update)
			tmpGroupPy.Post("/alive", group.Alive)
			tmpGroupPy.Post("/del", group.Del)
			tmpGroupPy.Post("/copy", group.Copys)
			tmpGroupPy.Post("/recipient/set", group.SetRecipient)
			tmpGroupPy.Post("/smsplatform/set", group.SetSmsPlatom)
			tmpGroupPy.Post("/smsplatform/setall", group.SetSmsPlatom)

		}
		tmpPy := v1.Party("/template")
		{
			temp := controllers.Template{}
			tmpPy.Get("/getall", temp.Gets)
			tmpPy.Post("/add", temp.Add)
			tmpPy.Post("/alive", temp.Update)
			tmpPy.Post("/update", temp.Update)
			tmpPy.Post("/madd", temp.Adds) //批量添加
			tmpPy.Post("/mupdate", temp.Updates) //批量更新
			tmpPy.Post("/saves", temp.Saves) //批量更新
		}

		hisPy := v1.Party("/templatehistory")
		{
			his := controllers.History{}
			hisPy.Post("/getall", his.List)
		}
		smsPy := v1.Party("/sms")
		{
			his := controllers.Sms{}
			smsPy.Post("/send", his.Send)
			smsPy.Post("/msend", his.Sends)
		}
		mailPy := v1.Party("/mail")
		{
			his := controllers.Email{}
			mailPy.Post("/send", his.Send)
			mailPy.Post("/msend", his.Sends)
		}

		dingPy := v1.Party("/ding")
		{
			his := controllers.DingDing{}
			dingPy.Post("/send", his.Send)
			dingPy.Post("/list-qun", his.GetQuns)
		}

		//v1.Get("/template/getall", a.getAllTemplatesByGroupId)  //按照groupid查temps
		//v1.Post("/templategroup/getall", a.List) //查groupids
		//v1.Post("/template/add", a.addTemplate)
		//v1.Post("/templategroup/recipient/set", a.SetTemplateGroupDefaultRecipient)
		//v1.Post("/templategroup/add", a.addTemplateGroup)
		//v1.Post("/templategroup/update", a.updateTemplateGroup)
		//v1.Post("/templategroup/alive", a.AliveTemplateGroup)
		//v1.Post("/templategroup/del", a.DelTemplateGroup)
		//v1.Post("/template/alive", a.AliveTemplate)
		//v1.Post("/sms/send", a.handleSmsSend)
		//v1.Post("/mail/send", a.handleMailSend)
		//v1.Post("/templatehistory/getall", a.getAllTemplateHistory)
		//v1.Post("/sms/msend", a.handleSmsMSend)
		//v1.Post("/mail/msend", a.handleMailMSend)
		//v1.Post("/template/update", a.updateTemplate)
		//v1.Post("/template/madd", a.mAddTemplate) //批量添加
		//v1.Post("/template/mupdate", a.mUpdateTemplate) //批量更新
		//v1.Post("/templategroup/smsplatform/set", a.SetTemplateGroupSmsPlatFrom)
		//v1.Post("/templategroup/smsplatform/setall", a.SetAllTemplateGroupSmsPlatFrom)
		//v1.Any("/", a.defaultRoot)
	}
}

//func (this *WebServer) getAllTemplatesByGroupId(ctx iris.Context) {
//	defer PanicPrint()
//	gidStr := ctx.URLParam("groupid")
//	if len(gidStr) == 0 {
//		ZapLog().With(zap.Error(errors.New("groupid is nil", 0))).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:nil groupid"))
//		return
//	}
//	gidStr = strings.TrimSpace(gidStr)
//	id, err := strconv.ParseUint(gidStr, 10, 32)
//	if err != nil {
//		ZapLog().With(zap.Error(err), zap.String("groupid", gidStr)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	tmps, err := this.mTemplateMgr.GetAllTempaltesByGroupIdFromDb(uint(id))
//	if err != nil {
//		ZapLog().With(zap.Error(err), zap.Any("groupid", gidStr)).Error("GetAllTempaltesByGroupIdFromDb err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//	for i := 0; i < len(tmps); i++ {
//		encodeString := base64.StdEncoding.EncodeToString([]byte(tmps[i].GetContent()))
//		tmps[i].SetContent(encodeString)
//	}
//	res := NewResNotifyMsg(ErrCode_Success, "")
//	res.SetTemplates(tmps)
//	ctx.JSON(res)
//}
//
////func (this *WebServer) getAllTemplateGroups(ctx iris.Context) {
////	defer PanicPrint()
////	param := new(TemplateGroupListParam)
////	if err := ctx.ReadJSON(param); err != nil {
////		ZapLog().With(zap.Error(err)).Error("param err")
////		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
////		return
////	}
////	totalLine:= param.Total_lines
////	pageIndex:= param.Page_index
////	pageNum:= param.Max_disp_lines
////	var err error
////	if totalLine == 0 {
////		totalLine, err = this.mTemplateMgr.CountTempaltesGroupFromDb(param.Type, param.GetName())
////		if err != nil {
////			ZapLog().With(zap.Error(err), zap.Any("type", param.Type), zap.String("likeName", param.GetName())).Error("CountTempaltesGroupFromDb err")
////			ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
////			return
////		}
////	}
////
////	if pageNum < 1 || pageNum > 100 {
////		pageNum = 50
////	}
////	beginIndex := pageNum * (pageIndex - 1)
////
////	gs, err := this.mTemplateMgr.GetAllTempaltesGroupFromDb(beginIndex, pageNum, param.Type, param.GetName())
////	if err != nil {
////		ZapLog().With(zap.Error(err), zap.Any("type", param.Type)).Error("GetAllTempaltesGroupFromDb err")
////		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
////		return
////	}
////	for i:=0; i < len(gs); i++ {
////		arr, err := this.mTemplateMgr.GetAllTempalteGroupLangsByGroupIdFromDb(gs[i].GetId())
////		if err != nil {
////			ZapLog().With(zap.Error(err), zap.Any("groupid", gs[i].GetId())).Error("GetAllTempalteGroupLangsByGroupIdFromDb err")
////			ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
////			return
////		}
////		gs[i].Langs = arr
////		hisArr, err := this.mTemplateMgr.GetAllTempalteHistoryFromDb(0, 3, gs[i].GetId())
////		if err != nil {
////			ZapLog().With(zap.Error(err), zap.Any("groupid", gs[i].GetId())).Error("GetAllTempalteHistoryFromDb err")
////			continue
////		}
////		gs[i].RateFail = GenRateFailArr(hisArr)
////	}
////	res := NewResNotifyMsg(ErrCode_Success, "")
////	ll := new(TemplateGroupList)
////	ll.SetTemplateGroups(gs)
////	ll.Max_disp_lines = pageNum
////	ll.Page_index = pageIndex
////	ll.Total_lines = totalLine
////	res.SetTemplateGroupList(ll)
////	ctx.JSON(res)
////
////}
//
//func (this *WebServer) addTemplate(ctx iris.Context) {
//	defer PanicPrint()
//	template := new(Template)
//	if err := ctx.ReadJSON(template); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	template.Id = nil
//	if template.GroupId == nil || *template.GroupId == 0 {
//		ZapLog().Error("param err:GroupId is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	if template.Content == nil || len(*template.Content) == 0 {
//		ZapLog().Error("param err:Content is nil")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	if template.Alive == nil {
//		template.SetAlive(0)
//	}
//	if template.SmsPlatform == nil {
//		template.SetSmsPlatform(SMSPlatform_AWS)
//	}
//	decodeBytes, err := base64.StdEncoding.DecodeString(template.GetContent())
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("base64.StdEncoding.DecodeString err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:content base64Decode, "+err.Error()))
//		return
//	}
//	template.SetContent(string(decodeBytes))
//	err = this.mTemplateMgr.AddTemplate(template)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("AddTemplate err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) mAddTemplate(ctx iris.Context) {
//	defer PanicPrint()
//	templateArr := make([]*Template, 0)
//	if err := ctx.ReadJSON(templateArr); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	for i:=0; i < len(templateArr); i++ {
//		templateArr[i].Id = nil
//		if templateArr[i].GroupId == nil || *templateArr[i].GroupId == 0 {
//			ZapLog().Error("param err:GroupId is nil or 0")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//			return
//		}
//		if templateArr[i].Content == nil || len(*templateArr[i].Content) == 0 {
//			ZapLog().Error("param err:Content is nil")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//			return
//		}
//		if templateArr[i].Alive == nil {
//			templateArr[i].SetAlive(0)
//		}
//		decodeBytes, err := base64.StdEncoding.DecodeString(templateArr[i].GetContent())
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("base64.StdEncoding.DecodeString err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:content base64Decode, "+err.Error()))
//			return
//		}
//		templateArr[i].SetContent(string(decodeBytes))
//		err = this.mTemplateMgr.AddTemplate(templateArr[i])
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("AddTemplate err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//			return
//		}
//	}
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) updateTemplate(ctx iris.Context) {
//	defer PanicPrint()
//	template := new(Template)
//	if err := ctx.ReadJSON(template); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if template.Id == nil || *template.Id == 0 {
//		ZapLog().Error("param err:Id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	if template.Content != nil  {
//		decodeBytes, err := base64.StdEncoding.DecodeString(template.GetContent())
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("base64.StdEncoding.DecodeString err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:content base64Decode, "+err.Error()))
//			return
//		}
//		template.SetContent(string(decodeBytes))
//	}
//	err := this.mTemplateMgr.UpdateTemplate(template)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("AddTemplate err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) mUpdateTemplate(ctx iris.Context) {
//	defer PanicPrint()
//	templateAddr := make([]*Template, 0 )
//	if err := ctx.ReadJSON(templateAddr); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//
//	for i:=0; i < len(templateAddr); i++ {
//		if templateAddr[i].Id == nil || *templateAddr[i].Id == 0 {
//			ZapLog().Error("param err:Id is nil or 0")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//			return
//		}
//		if templateAddr[i].Content != nil  {
//			decodeBytes, err := base64.StdEncoding.DecodeString(templateAddr[i].GetContent())
//			if err != nil {
//				ZapLog().With(zap.Error(err)).Error("base64.StdEncoding.DecodeString err")
//				ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:content base64Decode, "+err.Error()))
//				return
//			}
//			templateAddr[i].SetContent(string(decodeBytes))
//		}
//		err := this.mTemplateMgr.UpdateTemplate(templateAddr[i])
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("AddTemplate err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//			return
//		}
//	}
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) addTemplateGroup(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	templateGroup.Id = nil
//	if len(templateGroup.GetName()) == 0{
//		ZapLog().With(zap.String("error", "nil name")).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:nil name"))
//		return
//	}
//	if templateGroup.Alive == nil {
//		templateGroup.SetAlive(0)
//	}
//	if templateGroup.SmsPlatform == nil {
//		templateGroup.SetSmsPlatform(SMSPlatform_AWS)
//	}
//	templateGroup.SetName(strings.Replace(templateGroup.GetName(), " ", "", len(templateGroup.GetName())))
//	mm := make(map[string]interface{})
//	mm["name"] = templateGroup.GetName()
//	mm["type"] = templateGroup.GetType()
//	exist, err := this.mTemplateMgr.ExistTemplateGroup(mm)//name+type唯一
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("ExistTemplateGroup err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//	if exist {
//		ZapLog().With(zap.String("name", templateGroup.GetName())).Error("name already have")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Exist, "exist name"))
//		return
//	}
//	err = this.mTemplateMgr.AddTemplateGroup(templateGroup)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("AddTemplateGroup err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) updateTemplateGroup(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if templateGroup.Id == nil || *templateGroup.Id == 0 {
//		ZapLog().Error("param err:id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	if len(templateGroup.GetName()) != 0{
//		mm := make(map[string]interface{})
//		mm["name"] = templateGroup.GetName()
//		mm["type"] = templateGroup.GetType()//这里有bug，需要从数据库查询得到
//		exist, err := this.mTemplateMgr.ExistTemplateGroup(mm)
//		if err != nil {
//			ZapLog().With(zap.Error(err)).Error("ExistTemplateGroup err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//			return
//		}
//		if exist {
//			ZapLog().With(zap.String("name", templateGroup.GetName())).Error("name already have")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Exist, "exist name"))
//			return
//		}
//	}
//
//	err := this.mTemplateMgr.UpdateTemplateGroup(templateGroup)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("AddTemplateGroup err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
////tempId, title，params
//func (this *WebServer) handleSmsSend(ctx iris.Context) {
//	defer PanicPrint()
//	req := new(ReqNotifyMsg)
//	err := ctx.ReadJSON(req)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	ZapLog().With(zap.Any("param", req)).Info("handleSmsSend start")
//	errCode, errMsg := this.smsSend(req, true)
//	resMsg := NewResNotifyMsg(errCode, errMsg)
//	ctx.JSON(resMsg)
//	if errCode == ErrCode_Success {
//		ZapLog().With(zap.Any("param", req)).Info("handleSmsSend success")
//	}else{
//		ZapLog().With(zap.Any("param", req), zap.Int("errCode", errCode), zap.String("errmsg", errMsg)).Error("handleSmsSend fail")
//	}
//}
//
//func (this *WebServer) handleSmsMSend(ctx iris.Context) {
//	defer PanicPrint()
//	req := make([]*ReqNotifyMsg, 0)
//	err := ctx.ReadJSON(req)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	ZapLog().With(zap.Any("param", req)).Info("handleSmsMSend start")
//
//	resMsg := make([]ResNotifyMsg, len(req), len(req))
//	for m:=0; m < len(req); m++ {
//		errCode, errMsg := this.smsSend(req[m], true)
//		resMsg[m].SetErr(errCode)
//		resMsg[m].SetErrMsg(errMsg)
//	}
//	ctx.JSON(resMsg)
//	ZapLog().With(zap.Any("param", req)).Info("handleSmsMSend success")
//}
//
//func (this *WebServer) smsSendTest(ctx iris.Context) {
//	tos := ctx.URLParam("to")
//	toArr := strings.Split(tos, ",")
//	for i:=0; i<len(toArr); i++ {
//		if len(toArr[i]) == 0 {
//			continue
//		}
//		if err := this.mSmsMgr.DirectSendTwl(fmt.Sprintf("smsSendTest %d", toArr[i]), toArr[i]); err!= nil {
//			ZapLog().With(zap.Error(err), zap.Int("i", i)).Error("param err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_Param, "smsSendTest err:"+err.Error()+fmt.Sprintf("%d", i)))
//			return
//		}
//		ZapLog().With(zap.Int("i", i), zap.String("to", toArr[i])).Info("smsSendTest success")
//	}
//}
//
////返回值 errCode，errMsg， []sub_errCode, []sub_errMsg
//func (this *WebServer) smsSend(req *ReqNotifyMsg, recordFlag bool) (int, string){
//	ZapLog().With(zap.Any("param", req)).Info("SmsSend start")
//	if !this.ReqNotifyMsgIsValid(req) {
//		ZapLog().With(zap.Any("req", *req), zap.Error(errors.New("param noValid", 0))).Error("param err")
//		return ErrCode_Param, "param err"
//	}
//	var err error
//	var tmplate *Template
//	if req.TempId != nil {
//		tmplate, err = this.mTemplateMgr.GetTempalteById(req.GetTempId())
//	} else if req.TempAlias != nil {
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByAlias(req.GetTempAlias(), Notify_Type_Sms)
//	} else if req.GroupId != nil {
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByGroupId(req.GetGroupId(), req.GetLang())
//	} else if req.GroupName != nil {
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByGroupName(req.GetGroupName(), req.GetLang(), Notify_Type_Sms)
//	} else if req.GroupAlias != nil{
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByGroupAlias(req.GetGroupAlias(), req.GetLang(), Notify_Type_Sms)
//	} else {
//		ZapLog().With(zap.Any("req", *req)).Error("NoFindAliveTemplate err")
//		return ErrCode_NoFindAliveTemplate, "NoFindAliveTemplate"
//	}
//	if err != nil {
//		ZapLog().With(zap.Any("req", *req), zap.Error(err)).Error("GetTempalte err")
//		return ErrCode_InerServer, err.Error()
//	}
//	ZapLog().Sugar().Infof("template groupid[%d] id[%d] name[%d]", tmplate.GetGroupId(), tmplate.GetId(), tmplate.GetName())
//	if tmplate.GetALive() == Notify_AliveMode_Dead {
//		ZapLog().With(zap.Any("req", *req),zap.Uint("tempid", tmplate.GetId()), zap.Int("alive", tmplate.GetALive())).Error("GetTempalte err:template not alive")
//		go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Sms, 0,req.GetRecipientSize(), recordFlag)
//		return ErrCode_NoAlive, "template not alive"
//	}
//	if tmplate.GetType() != Notify_Type_Sms {
//		ZapLog().With(zap.Uint("tempid", tmplate.GetId()), zap.Error(errors.New("search wrong type", 0))).Error("type err")
//		return ErrCode_InerServer, "search wrong type"
//	}
//	if req.GetUseDefaultRecipient() {
//		AppendRecipient(req, tmplate.GetDefaultRecipient())
//	}
//	if req.Recipient == nil || len(req.Recipient) == 0 {
//		ZapLog().With(zap.Error(errors.New("no Recipient", 0))).Error("param err")
//		return ErrCode_NoRecipient, "param err:no Recipient"
//	}
//	ZapLog().Info("sms all recipient", zap.Any("recipient", req.Recipient))
//
//	smsBody, err := this.mTemplateMgr.ParseTextTemplate(tmplate.GetContent(), req.Params)
//	if err != nil {
//		ZapLog().With(zap.Any("tempParam", req.Params),zap.Uint("tempid", tmplate.GetId()),zap.Error(err)).Error("ParseTextTemplate err")
//		go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Mail, 0,req.GetRecipientSize(), recordFlag)
//		return ErrCode_TemplateParse, err.Error()
//	}
//
//	haveSubErr := false
//	subErrMsg := ""
//	failCount := 0
//	if tmplate.GetSmsPlatform() == SMSPlatform_TWL {
//		for i:=0; i < len(req.Recipient); i++ {
//			if len(req.Recipient[i]) == 0 {
//				continue
//			}
//			if err := this.mSmsMgr.DirectSendTwl(smsBody, req.Recipient[i]); err!= nil {
//				failCount++
//				haveSubErr = true
//				subErrMsg += " " + err.Error()
//				ZapLog().With(zap.Error(err), zap.String("recipient", req.Recipient[i])).Error("DirectSendTWL err")
//			}
//		}
//	}
//	if tmplate.GetSmsPlatform() == SMSPlatform_CHUANGLAN {
//		zhPhones,noZhPhones := ChuangLanSplitPhones(req.Recipient)
//		//newClSmsBody,newClParam := ParseSmsBodyChuanglan(tmplate.GetContent(), req.Params)
//		if num, err := this.mSmsMgr.DirectSendChuanglan(smsBody, zhPhones, nil);err != nil {
//			failCount += len(zhPhones)
//			haveSubErr = true
//			subErrMsg += " " + err.Error()
//			ZapLog().With(zap.Error(err), zap.Any("recipient",zhPhones)).Error("DirectSendChuanglan err")
//			//这里不能退出 还得继续发送
//		}else{
//			failCount += len(zhPhones) - num
//		}
//		for i:=0; i < len(noZhPhones); i++ {
//			if len(noZhPhones[i]) == 0 {
//				continue
//			}
//			err = this.mSmsMgr.DirectSendAws(smsBody, noZhPhones[i])
//			if err != nil {
//				failCount++
//				haveSubErr = true
//				subErrMsg += " " + err.Error()
//				ZapLog().With(zap.Error(err), zap.String("recipient", noZhPhones[i])).Error("DirectSendAws err")
//			}
//		}
//	}
//	if tmplate.GetSmsPlatform() == SMSPlatform_AWS {
//		for i:=0; i < len(req.Recipient); i++ {
//			if len(req.Recipient[i]) == 0 {
//				continue
//			}
//			err = this.mSmsMgr.DirectSendAws(smsBody, req.Recipient[i])
//			if err != nil {
//				failCount++
//				haveSubErr = true
//				subErrMsg += " " + err.Error()
//				ZapLog().With(zap.Error(err), zap.String("recipient", req.Recipient[i])).Error("DirectSendAWS err")
//			}
//		}
//	}
//
//	if haveSubErr {
//		go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Sms, req.GetRecipientSize() - failCount, failCount, recordFlag)
//		return ErrCode_SubRecipient, subErrMsg
//	}
//
//	go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Sms, req.GetRecipientSize(),0, recordFlag)
//	ZapLog().With(zap.Any("param", req)).Info("SmsSend success")
//	return ErrCode_Success,""
//}
//
//func (this *WebServer) handleMailSend(ctx iris.Context) {
//	defer PanicPrint()
//	req := new(ReqNotifyMsg)
//	err := ctx.ReadJSON(req)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	ZapLog().With(zap.Any("param", req)).Info("handleMailSend start")
//	errCode, errMsg := this.mailSend(req, true)
//	resMsg := NewResNotifyMsg(errCode, errMsg)
//	ctx.JSON(resMsg)
//	if errCode == ErrCode_Success {
//		ZapLog().With(zap.Any("param", req)).Info("handleMailSend success")
//	}else{
//		ZapLog().With(zap.Any("param", req), zap.Int("errCode", errCode), zap.String("errmsg", errMsg)).Error("handleMailSend fail")
//	}}
//
//func (this *WebServer) handleMailMSend(ctx iris.Context) {
//	defer PanicPrint()
//	req := make([]*ReqNotifyMsg, 0)
//	err := ctx.ReadJSON(req)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	ZapLog().With(zap.Any("param", req)).Info("handleMailMSend start")
//	resMsg := make([]ResNotifyMsg, len(req), len(req))
//	for m:=0; m < len(req); m++ {
//		errCode, errMsg := this.mailSend(req[m], true)
//		resMsg[m].SetErr(errCode)
//		resMsg[m].SetErrMsg(errMsg)
//	}
//	ctx.JSON(resMsg)
//	ZapLog().With(zap.Any("param", req)).Info("handleMailMSend success")
//}
//
//func (this *WebServer) mailSend(req* ReqNotifyMsg, recordFlag bool) (int, string){
//	ZapLog().With(zap.Any("param", req)).Info("MailSend start")
//	if !this.ReqNotifyMsgIsValid(req) {
//		ZapLog().With(zap.Any("req", *req), zap.Error(errors.New("param noValid", 0))).Error("param err")
//		return ErrCode_Param, "param err"
//	}
//	var err error
//	var tmplate *Template
//	if req.TempId != nil {
//		tmplate, err = this.mTemplateMgr.GetTempalteById(req.GetTempId())
//	} else if req.TempAlias != nil {
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByAlias(req.GetTempAlias(), Notify_Type_Mail)
//	} else if req.GroupId != nil {
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByGroupId(req.GetGroupId(), req.GetLang())
//	} else if req.GroupName != nil {
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByGroupName(req.GetGroupName(), req.GetLang(), Notify_Type_Mail)
//	} else if req.GroupAlias != nil{
//		tmplate, err = this.mTemplateMgr.GetAliveTempalteByGroupAlias(req.GetGroupAlias(), req.GetLang(), Notify_Type_Mail)
//	} else {
//		ZapLog().With(zap.Any("req", *req)).Error("NoFindAliveTemplate err")
//		return ErrCode_NoFindAliveTemplate, "NoFindAliveTemplate"
//	}
//	if err != nil {
//		ZapLog().With(zap.Any("req", *req), zap.Error(err)).Error("GetTempalte err:or not alive")
//		return ErrCode_InerServer, err.Error()+" or not alive"
//	}
//	ZapLog().Sugar().Infof("template groupid[%d] id[%d] name[%d]", tmplate.GetGroupId(), tmplate.GetId(), tmplate.GetName())
//	if tmplate.GetALive() == Notify_AliveMode_Dead {
//		ZapLog().With(zap.Any("req", *req), zap.Int("tempAlive", tmplate.GetALive()), zap.Uint("Id", tmplate.GetId())).Error("GetTempalte err:template not alive")
//		go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Mail, 0,1, recordFlag)
//		return ErrCode_NoAlive, "template not alive"
//	}
//	if tmplate.GetType() != Notify_Type_Mail {
//		ZapLog().With(zap.Any("req", *req), zap.Any("Id", tmplate.GetId()), zap.Error(errors.New("search wrong type", 0))).Error("type err")
//		return ErrCode_InerServer, "search wrong type"
//	}
//	if req.GetUseDefaultRecipient() {
//		AppendRecipient(req, tmplate.GetDefaultRecipient())
//	}
//	if req.Recipient == nil || len(req.Recipient) == 0 {
//		ZapLog().With(zap.Any("req", *req), zap.Error(errors.New("no Recipient", 0))).Error("param err")
//		return ErrCode_NoRecipient, "param err:no Recipient"
//	}
//
//	ZapLog().Info("mail all recipient", zap.Any("recipient", req.Recipient))
//
//	if req.Params == nil {
//		req.Params = make(map[string]interface{})
//	}
//	req.Params["title_key"] = tmplate.GetTitle()
//	_, body, err := this.mTemplateMgr.ParseHtmlTemplate(this.mTemplateMgr.BodyToHtml(tmplate.GetContent()), req.Params)
//	if err != nil {
//		ZapLog().With(zap.Any("tempParam", req.Params), zap.Error(err), zap.Uint("Id", tmplate.GetId())).Error("ParseHtmlTemplate err")
//		go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Mail, 0,1, recordFlag)
//		return ErrCode_TemplateParse, err.Error()
//	}
//
//	haveSubErr := false
//	subErrMsg := ""
//	failCount := 0
//	for i:=0; i < len(req.Recipient); i++ {
//		if len(req.Recipient[i]) == 0 {
//			continue
//		}
//		err = this.mMailMgr.DirectSend(tmplate.GetTitle(), body, req.Recipient[i])
//		if err != nil {
//			failCount++
//			haveSubErr = true
//			subErrMsg += " " + err.Error()
//			ZapLog().With(zap.Error(err), zap.String("recipient", req.Recipient[i])).Error("DirectSend err")
//		}
//	}
//
//	if haveSubErr {
//		go this.recordHistory(tmplate.GetGroupId(),Notify_Type_Mail, req.GetRecipientSize() - failCount, failCount, recordFlag)
//		return ErrCode_SubRecipient, subErrMsg
//	}
//
//	go this.recordHistory(tmplate.GetGroupId(), Notify_Type_Mail,req.GetRecipientSize(),0, recordFlag)
//	ZapLog().With(zap.Any("param", req)).Info("MailSend success")
//	return ErrCode_Success, ""
//}
//
//func (this *WebServer) DelTemplateGroup(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if templateGroup.Id == nil {
//		ZapLog().Error("param err:id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	err := this.mTemplateMgr.DelTemplateGroupAndTemplates(templateGroup.GetId())
//	if err != nil {
//		ZapLog().With(zap.Any("req", *templateGroup), zap.Error(err)).Error("DelTemplateGroupAndTemplates err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) AliveTemplateGroup(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if templateGroup.Id == nil || *templateGroup.Id == 0 {
//		ZapLog().Error("param err:id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	if templateGroup.Alive == nil {
//		ZapLog().Error("param err:alive is nil")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//
//	err := this.mTemplateMgr.AliveTemplateGroupAndTemplates(templateGroup.GetId(), templateGroup.GetAlive())
//	if err != nil {
//		ZapLog().With(zap.Any("req", *templateGroup), zap.Error(err)).Error("AliveTemplateGroupAndTemplates err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
////这个接口可以用update代替
//func (this *WebServer) AliveTemplate(ctx iris.Context) {
//	defer PanicPrint()
//	template := new(Template)
//	if err := ctx.ReadJSON(template); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if template.Id == nil || *template.Id == 0 {
//		ZapLog().Error("param err:id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//
//	err := this.mTemplateMgr.UpdateTemplate(template)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("UpdateTemplate err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, "UpdateTemplate err:"+err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) SetTemplateGroupSmsPlatFrom(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if templateGroup.Id == nil || *templateGroup.Id == 0 {
//		ZapLog().Error("param err:id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//	if templateGroup.SmsPlatform == nil {
//		ZapLog().Error("param err:SmsPlatform is nil")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//
//	err := this.mTemplateMgr.SetSmsPlatformTemplateGroupAndTemplates(templateGroup.GetId(), templateGroup.GetSmsPlatform())
//	if err != nil {
//		ZapLog().With(zap.Any("req", *templateGroup), zap.Error(err)).Error("SetSmsPlatformTemplateGroupAndTemplates err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) SetAllTemplateGroupSmsPlatFrom(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if templateGroup.SmsPlatform == nil {
//		ZapLog().Error("param err:SmsPlatform is nil")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err"))
//		return
//	}
//
//	err := this.mTemplateMgr.SetAllSmsPlatformTemplateGroupAndTemplates(templateGroup.GetSmsPlatform())
//	if err != nil {
//		ZapLog().With(zap.Any("req", *templateGroup), zap.Error(err)).Error("SetAllSmsPlatformTemplateGroupAndTemplates err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
//func (this *WebServer) SetTemplateGroupDefaultRecipient(ctx iris.Context) {
//	defer PanicPrint()
//	templateGroup := new(TemplateGroup)
//	if err := ctx.ReadJSON(templateGroup); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if templateGroup.Id == nil || *templateGroup.Id == 0 {
//		ZapLog().Error("param err:id is nil or 0")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:id is nil"))
//		return
//	}
//	if templateGroup.DefaultRecipient == nil {
//		ZapLog().Error("param err:DefaultRecipient is nil")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:default_recipient is nil"))
//		return
//	}
//	*templateGroup.DefaultRecipient = strings.Replace(*templateGroup.DefaultRecipient, " ", "", len(*templateGroup.DefaultRecipient))
//
//	err := this.mTemplateMgr.SetDefaultRecipientOfTemplateGroupAndTemplates(templateGroup.GetId(), templateGroup.GetDefaultRecipient())
//	if err != nil {
//		ZapLog().With(zap.Any("req", *templateGroup), zap.Error(err)).Error("AliveTemplateGroupAndTemplates err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//
//	resMsg := NewResNotifyMsg(ErrCode_Success, "")
//	ctx.JSON(resMsg)
//}
//
////flag 是需要的，避免变成死循环
//func (this *WebServer) recordHistory(gid uint,tp, succ, fail int, flag bool) {
//	defer PanicPrint()
//	informAdminFlag, err := this.mTemplateMgr.IncrTemplateHistoryCount(gid, tp, succ, fail)
//	if err != nil {
//		ZapLog().With(zap.Error(err)).Error("IncrTemplateHistoryCount err")
//		return
//	}
//	if !informAdminFlag {
//		return
//	}
//	if !flag {
//		return
//	}
//	mon := &config.GConfig.Monitor
//	ZapLog().Warn("Monitor start", zap.Uint("groupid", gid))
//	if len(mon.TmpGNameMail) != 0 {
//		req := new(ReqNotifyMsg)
//		req.Params = make(map[string]interface{})
//		req.SetGroupName(mon.TmpGNameMail)
//		req.SetLang(mon.TmpLangMail)
//		req.Params["key1"] = gid
//		errCode, errMsg := this.mailSend(req, false)
//		if errCode != ErrCode_Success {
//			ZapLog().With(zap.String("error", errMsg), zap.Int("errcode",errCode )).Error("monitor sendMail err")
//		}
//	}
//	if len(mon.TmpGNameSms) != 0 {
//		req := new(ReqNotifyMsg)
//		req.Params = make(map[string]interface{})
//		req.Params["key1"] = gid
//		req.SetGroupName(mon.TmpGNameSms)
//		req.SetLang(mon.TmpLangSms)
//		errCode, errMsg := this.smsSend(req, false)
//		if errCode != ErrCode_Success {
//			ZapLog().With(zap.String("error", errMsg), zap.Int("errcode",errCode )).Error("monitor sendSms err")
//		}
//	}
//
//}
//
//func (this *WebServer) getAllTemplateHistory(ctx iris.Context) {
//	defer PanicPrint()
//	param := new(TemplateHistoryListParam)
//	if err := ctx.ReadJSON(param); err != nil {
//		ZapLog().With(zap.Error(err)).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:"+err.Error()))
//		return
//	}
//	if param.GroupId == 0 {
//		ZapLog().With(zap.String("error", "groupid is 0")).Error("param err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_Param, "param err:groupid"))
//		return
//	}
//	totalLine:= param.Total_lines
//	pageIndex:= param.Page_index
//	pageNum:= param.Max_disp_lines
//	var err error
//	if totalLine == 0 {
//		totalLine, err = this.mTemplateMgr.CountTempalteHistoryFromDb(param.GroupId)
//		if err != nil {
//			ZapLog().With(zap.Error(err), zap.Any("GroupId", param.GroupId)).Error("CountTempalteHistoryFromDb err")
//			ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//			return
//		}
//	}
//
//	if pageNum < 1 || pageNum > 100 {
//		pageNum = 50
//	}
//	beginIndex := pageNum * (pageIndex - 1)
//
//	gs, err := this.mTemplateMgr.GetAllTempalteHistoryFromDb(beginIndex, pageNum, param.GroupId)
//	if err != nil {
//		ZapLog().With(zap.Error(err), zap.Any("GroupId", param.GroupId)).Error("GetAllTempalteHistoryFromDb err")
//		ctx.JSON(NewResNotifyMsg(ErrCode_InerServer, err.Error()))
//		return
//	}
//	res := NewResNotifyMsg(ErrCode_Success, "")
//	ll := new(TemplateHistoryList)
//	ll.TemplateHistorys = gs
//	ll.Max_disp_lines = pageNum
//	ll.Page_index = pageIndex
//	ll.Total_lines = totalLine
//	res.TemplateHistoryList = ll
//	ctx.JSON(res)
//}
//
//func (this *WebServer) ReqNotifyMsgIsValid(req *ReqNotifyMsg) bool {
//	//if req.Recipient == nil {
	//	return false
	//}
//	if req.TempId != nil {
//		return true
//	}
//	if req.TempAlias != nil {
//		return true
//	}
//	if (req.GroupId == nil) &&(req.GroupName == nil)&&(req.GroupAlias == nil) {
//		return false
//	}
//	if req.Lang == nil {
//		return false
//	}
//	return true
//}

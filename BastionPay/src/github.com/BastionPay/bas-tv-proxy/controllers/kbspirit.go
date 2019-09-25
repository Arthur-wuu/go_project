package controllers

import (
	"github.com/kataras/iris"
	"BastionPay/bas-tv-proxy/models/kbspirit"

	"strings"
	"strconv"
	"BastionPay/bas-tv-proxy/api"
	"BastionPay/bas-tv-proxy/models"
	"BastionPay/bas-tv-proxy/common"
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"BastionPay/bas-tv-proxy/config"
)

var GKBSpirit KBSpirit

type KBSpirit struct {
	Controllers
}

func (this * KBSpirit) Init() error {
	models.GKBSpirit.Init()
	models.GKBSpirit.Start()
	return nil
}

// input&type(暂时不用)&market(空或者*为全部，否则匹配指定市场)&viplevel
func (this * KBSpirit) HandleGetObjs(ctx iris.Context) {
	env,err := this.paraseGetObjsParam(ctx)
	if  err != nil {
		this.ExceptionSerive(ctx, api.ErrCode_Param)
		return
	}
	ZapLog().Info("HandleGetObjs",zap.Any("env", *env))

	data := models.GKBSpirit.GetObjs(env)

	//total := int64(env.Count)
	//env.Count = env.Count + 20 //
	//
	//kbspirit.GKBSpirit.GetObjs("", env)

	//data = this.removeDuplicate(data)
	//if len(data) != 0 {
	//	shuchu := &dzhyun.JPBShuChu{
	//		LeiXing: dzhyun.JPBLeiXing(t).Enum(),
	//		ShuJu:   data,
	//	}
	//	result = append(result, shuchu)
	//	tpName := dzhyun.JPBLeiXing(t)
	//	if dzhyun.JPBLeiXing_TYPE_OBJATFUND == tpName || dzhyun.JPBLeiXing_TYPE_OBJATBLOCK == tpName {
	//		continue
	//	}
	//	// 截取结果
	//	if total-int64(len(data)) < 0 {
	//		shuchu.ShuJu = shuchu.ShuJu[:total]
	//		break
	//	}
	//	total -= int64(len(data))
	//}
	bigMsg := api.NewMSG(int32(api.EnumID_IDJianPanBaoShuChu)).AddJianPanBaoShuChu(
		[]*api.JianPanBaoShuChu{
			&api.JianPanBaoShuChu{
				GuanJianZi: &env.Input,
				JieGuo:     data,
			},
		}...
	)

	this.Response(ctx, bigMsg)
}

func (this *KBSpirit) paraseGetObjsParam(ctx iris.Context) (*kbspirit.Env,error) {
	env := kbspirit.NewSearchEnv()

	input := strings.ToUpper(ctx.URLParam("input"))
	if strings.HasPrefix(input, "_") || strings.Contains(input, " _"){
		env.SuffixFlag = true
	}
	env.Input = strings.Replace(input, " ", "", -1)
	env.Input = strings.Replace(env.Input, "_", "", -1)
	strMarkets := strings.Replace(strings.ToUpper(ctx.URLParam("market")), " ", "", -1)
	if len(strMarkets) == 0 || strMarkets[0] == '*'{
		env.Markets = config.GPreConfig.Markets
	}else{
		env.Markets = strings.Split(strMarkets, ",")
	}

	// 个数
	strCount := ctx.URLParam("count")
	env.Count, _ = strconv.ParseInt(strCount, 10, 64)
	if env.Count == 0 || env.Count > 100 {
		env.Count = 100
	}

	vipStr := ctx.URLParam("viplevel")
	env.VipLevel,_ = strconv.Atoi(vipStr)

	//// 是否退市
	//delist := true
	//strDelist := param.Values.Get("delist")
	//if "" != strDelist {
	//	delist, err = strconv.ParseBool(strDelist)
	//}
	//env.Delist = delist
	if len(env.Input) >= 6 {
		env.ChineseFlag =  common.IsChineseChar(env.Input[:6])
	}

	return env,nil
}

// intSlice 将字符串分割为整型数组
func (self *KBSpirit) intSlice(s string, sep string) ([]int64, error) {
	result := make([]int64, 0)
	if "" == s {
		return result, nil
	}
	strs := strings.Split(s, sep)
	for _, str := range strs {
		i, err := strconv.ParseInt(str, 10, 64)
		if nil != err {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}
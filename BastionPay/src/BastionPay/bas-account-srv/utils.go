package main

import (
	"BastionPay/bas-base/config"
	notify "BastionPay/bas-tools/sdk.notify.mail"
	l4g "github.com/alecthomas/log4go"
)

func InitBasNotify(cfgPath string) error {
	url := ""
	err := config.LoadJsonNode(cfgPath, "bas_notify", &url)
	if err != nil {
		l4g.Error("LoadJsonNode bas_notify err[%s]", err.Error())
		return err
	}
	err = notify.GNotifySdk.Init(url, "bas-account-srv")
	if err != nil {
		l4g.Error("GNotifySdk Init[%s] err[%s]", url, err.Error())
		return err
	}
	l4g.Info("bas_notify url[%s]", url)
	return nil
}

func GetAuditeTemplateName(cfgPath string) string {
	sms := ""
	err := config.LoadJsonNode(cfgPath, "audite_doing_temp", &sms)
	if err != nil {
		l4g.Error("LoadJsonNode bas_notify temp[audite_doing_temp] err[%s]", err.Error())
	}
	//	l4g.Info("tmp_name[%s]", sms)
	return sms
}

func GetUserFrozen1ForAdminTemp(cfgPath string) string {
	sms := ""
	err := config.LoadJsonNode(cfgPath, "userfrozen_1_foradmin_temp", &sms)
	if err != nil {
		l4g.Error("LoadJsonNode bas_notify tmp[userfrozen_1_foradmin] err[%s]", err.Error())
	}
	//	l4g.Info("tmp_name[%s]", sms)
	return sms
}

func GetUserFrozen1ForUserTemp(cfgPath string) string {
	sms := ""
	err := config.LoadJsonNode(cfgPath, "userfrozen_1_foruser_temp", &sms)
	if err != nil {
		l4g.Error("LoadJsonNode bas_notify tmp[userfrozen_1_foruser] err[%s]", err.Error())
	}
	//	l4g.Info("tmp_name[%s]", sms)
	return sms
}

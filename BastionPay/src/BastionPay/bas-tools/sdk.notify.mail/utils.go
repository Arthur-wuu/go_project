package sdk_notify_mail

func GenReqNotifyMsgWithGroupName(name, lang string, recipient []string, params map[string]interface{}) *ReqNotifyMsg {
	if len(name) == 0 || len(lang) == 0 {
		return nil
	}
	req := new(ReqNotifyMsg)
	req.SetGroupName(name)
	req.SetLang(lang)
	req.Recipient = recipient
	req.Params = params
	return req
}

func ReqNotifyMsgArrSetAppName(reqs []*ReqNotifyMsg, appName string) {
	for i := 0; i < len(reqs); i++ {
		reqs[i].SetAppName(appName)
	}
}

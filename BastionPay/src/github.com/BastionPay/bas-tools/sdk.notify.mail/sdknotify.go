package sdk_notify_mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var GNotifySdk NotifySdk

type NotifySdk struct {
	mUrl     string
	mAppName string
	mFlag    bool
	sync.Mutex
}

/*
*	addr: url
*   appName: 服务调用方名称
 */
func (this *NotifySdk) Init(addr, appName string) error {
	if this.mFlag {
		return nil
	}
	addr = strings.TrimRight(addr, "/")
	addr = strings.TrimSpace(addr)
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	if len(addr) == 0 || len(appName) == 0 {
		return errors.New("wrong param")
	}
	this.mUrl = addr
	this.mAppName = appName
	this.mFlag = true
	fmt.Println("init ", this.mFlag)
	return nil
}

/*
* name: 组名称，也就是 模板名称
* lang: 语言
* recipient: 多个接收人，只要有一个出错最终返回err，其余接收人任然能收到信息。若此需求不满足，可采用MSend接口
* params: 模板参数，参考模板内容
 */
func (this *NotifySdk) SendMailByGroupName(name, lang string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	if len(name) == 0 || len(lang) == 0 {
		return errors.New("nil name or lang")
	}
	req := new(ReqNotifyMsg)
	req.SetGroupName(name)
	req.SetLang(lang)
	req.SetAppName(this.mAppName)
	req.Recipient = recipient
	req.Params = params
	res, err := this.send(req, "/v1/notify/mail/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

//参数同上
func (this *NotifySdk) SendSmsByGroupName(name, lang string, recipient []string, params map[string]interface{}, useDefaultRecipient ...bool) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	if len(name) == 0 || len(lang) == 0 {
		return errors.New("nil name or lang")
	}
	req := new(ReqNotifyMsg)
	req.SetGroupName(name)
	req.SetAppName(this.mAppName)
	req.SetLang(lang)
	req.Recipient = recipient
	req.Params = params
	if len(useDefaultRecipient) > 0 {
		req.SetUseDefaultRecipient(useDefaultRecipient[0])
	}
	res, err := this.send(req, "/v1/notify/sms/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

/*多模板批量发送
* 参数：level 优先级
* 返回值：通用错误信息，各个模板错误信息
* 注意点：参数生成可采用 GenReqNotifyMsgWithGroupName
 */
func (this *NotifySdk) MSendMail(reqs []*ReqNotifyMsg, level int) (error, []error) {
	if !this.mFlag {
		return errors.New("not init"), nil
	}
	if len(reqs) == 0 {
		return errors.New("nil reqs"), nil
	}
	ReqNotifyMsgArrSetAppName(reqs, this.mAppName)
	resArr, err := this.msend(reqs, "/v1/notify/mail/msend", "POST")
	if err != nil {
		return err, nil
	}
	errArr := make([]error, len(reqs), len(reqs))
	for i := 0; i < len(resArr); i++ {
		if !resArr[i].IsOk() {
			errArr[i] = fmt.Errorf("%d %s", resArr[i].GetErr(), resArr[i].GetErrMsg())
		}
	}
	return nil, errArr
}

//参数同上
func (this *NotifySdk) MSendSms(reqs []*ReqNotifyMsg, level int) (error, []error) {
	if !this.mFlag {
		return errors.New("not init"), nil
	}
	if len(reqs) == 0 {
		return errors.New("nil reqs"), nil
	}
	ReqNotifyMsgArrSetAppName(reqs, this.mAppName)
	resArr, err := this.msend(reqs, "/v1/notify/sms/msend", "POST")
	if err != nil {
		return err, nil
	}
	errArr := make([]error, len(reqs), len(reqs))
	for i := 0; i < len(resArr); i++ {
		if !resArr[i].IsOk() {
			errArr[i] = fmt.Errorf("%d %s", resArr[i].GetErr(), resArr[i].GetErrMsg())
		}
	}
	return nil, errArr
}

/****************************************扩展接口*****************************************/
//根据组id和lang 发送
func (this *NotifySdk) SendMailById(id uint, lang string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	req := new(ReqNotifyMsg)
	req.SetGroupId(id)
	req.SetAppName(this.mAppName)
	req.SetLang(lang)
	req.Recipient = recipient
	req.Params = params
	res, err := this.send(req, "/v1/notify/mail/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

//根据组id和lang
func (this *NotifySdk) SendMailLvHighById(id uint, lang string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	req := new(ReqNotifyMsg)
	req.SetGroupId(id)
	req.SetAppName(this.mAppName)
	req.SetLang(lang)
	req.Recipient = recipient
	req.Params = params
	req.SetLevel(Notify_Level_High)
	res, err := this.send(req, "/v1/notify/mail/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

func (this *NotifySdk) SendSmsById(id uint, lang string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	req := new(ReqNotifyMsg)
	req.SetGroupId(id)
	req.SetAppName(this.mAppName)
	req.SetLang(lang)
	req.Recipient = recipient
	req.Params = params
	res, err := this.send(req, "/v1/notify/sms/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

func (this *NotifySdk) SendSmsLvHighById(id uint, lang string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	req := new(ReqNotifyMsg)
	req.SetGroupId(id)
	req.SetAppName(this.mAppName)
	req.SetLang(lang)
	req.Recipient = recipient
	req.Params = params
	req.SetLevel(Notify_Level_High)
	res, err := this.send(req, "/v1/notify/sms/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

func (this *NotifySdk) SendMailByTempAlias(alias string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	req := new(ReqNotifyMsg)
	req.SetTempAlias(alias)
	req.SetAppName(this.mAppName)
	req.Recipient = recipient
	req.Params = params
	res, err := this.send(req, "/v1/notify/mail/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

func (this *NotifySdk) SendMail(req *ReqNotifyMsg) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	res, err := this.send(req, "/v1/notify/mail/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

func (this *NotifySdk) SendSmsByTempAlias(alias string, recipient []string, params map[string]interface{}) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	req := new(ReqNotifyMsg)
	req.SetTempAlias(alias)
	req.SetAppName(this.mAppName)
	req.Recipient = recipient
	req.Params = params
	res, err := this.send(req, "/v1/notify/sms/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

func (this *NotifySdk) SendSms(req *ReqNotifyMsg) error {
	if !this.mFlag {
		return errors.New("not init")
	}
	res, err := this.send(req, "/v1/notify/sms/send", "POST")
	if err != nil {
		return err
	}
	if !res.IsOk() {
		return fmt.Errorf("%d %s", res.GetErr(), res.GetErrMsg())
	}
	return nil
}

/*******************************************************************************/

func (this *NotifySdk) send(req *ReqNotifyMsg, path, method string) (*ResNotifyMsg, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body := ioutil.NopCloser(bytes.NewReader(content))
	content, err = this.curl(path, method, body)
	if err != nil {
		return nil, err
	}
	res := new(ResNotifyMsg)
	if err = json.Unmarshal(content, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (this *NotifySdk) msend(req []*ReqNotifyMsg, path, method string) ([]ResNotifyMsg, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body := ioutil.NopCloser(bytes.NewReader(content))
	content, err = this.curl(path, method, body)
	if err != nil {
		return nil, err
	}
	res := make([]ResNotifyMsg, 0)
	if err = json.Unmarshal(content, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (this *NotifySdk) curl(path, method string, body io.Reader) ([]byte, error) {
	this.Lock()
	defer this.Unlock()
	newUrl := this.mUrl + path
	req, err := http.NewRequest(method, newUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return nil, errors.New("res no content")
	}
	if string(content) == "Not Found" {
		return nil, errors.New("res Not Found")
	}
	return content, nil
}

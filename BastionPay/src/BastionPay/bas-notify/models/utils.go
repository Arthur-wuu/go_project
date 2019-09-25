package models

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-notify/config"
	"BastionPay/bas-notify/models/table"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	html "html/template"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	text "text/template"
	"time"
)

func ParseHtmlTemplate(tpl string, params map[string]interface{}) (string, string, error) {
	tmp := html.New("mailTmp")
	t, err := tmp.Parse(tpl)
	if err != nil {
		return "", "", err
	}
	var tplBuf bytes.Buffer
	if err := t.Execute(&tplBuf, params); err != nil {
		return "", "", err
	}
	result := tplBuf.String()
	return getTemplateTitle(result), result, nil
}

func getTemplateTitle(result string) string {
	r, err := regexp.Compile(`<title>(.*)<\/title>`)
	if err != nil || r == nil {
		return ""
	}
	title := r.FindString(result)
	rr, err := regexp.Compile(`<[^>]*>`)
	if err != nil || rr == nil {
		return ""
	}
	return rr.ReplaceAllString(title, "")
}

func ParseTextTemplate(tmpBody string, params map[string]interface{}) (string, error) {
	tmp := text.New("smsTemp")
	_, err := tmp.Parse(tmpBody)
	if err != nil {
		return "", err
	}
	var tplBuf bytes.Buffer
	if err := tmp.Execute(&tplBuf, params); err != nil {
		return "", err
	}
	return tplBuf.String(), nil
}

//国内，国外，蓝创独用
func ChuangLanSplitPhones(phones []string) ([]string, []string) {
	zhPhones := make([]string, 0)
	noZhPhones := make([]string, 0)
	for i := 0; i < len(phones); i++ {
		if strings.HasPrefix(phones[i], "0086") || strings.HasPrefix(phones[i], "+86") {
			newPhone := strings.TrimLeft(phones[i], "0086")
			newPhone = strings.TrimLeft(newPhone, "+86")
			zhPhones = append(zhPhones, newPhone)
		} else {
			noZhPhones = append(noZhPhones, phones[i])
		}
	}
	return zhPhones, noZhPhones
}

//国内，国外，蓝创独用
func SplitChinaAndUnChinaPhones(phones []string) ([]string, []string) {
	zhPhones := make([]string, 0)
	noZhPhones := make([]string, 0)
	for i := 0; i < len(phones); i++ {
		if strings.HasPrefix(phones[i], "0086") || strings.HasPrefix(phones[i], "+86") {
			zhPhones = append(zhPhones, phones[i])
		} else {
			noZhPhones = append(noZhPhones, phones[i])
		}
	}
	return zhPhones, noZhPhones
}

func BodyToHtml(body string) string {
	head := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <!-- <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0" /> -->
    <title>{{.title_key}}</title>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="renderer" content="webkit">
</head>
<body style="margin: 0;padding:0">
`

	tail := `
            </body>
            </html>
            `

	return head + body + tail
}

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

func IncrHistoryCount(groupId int, tp, succ, fail int) (bool, error) {
	if succ == 0 && fail == 0 {
		return false, nil
	}
	day := GenDay()
	his := &TemplateHistory{
		GroupId: &groupId,
		Day:     &day,
		Type:    &tp,
		DaySucc: &succ,
		DayFail: &fail,
	}
	firstOverFlag, err := his.IncrTemplateHistoryCountAndRate(GenRateFail2, FirstOverRateFailThd)
	if err != nil {
		return false, err
	}
	return firstOverFlag, err
}

func GenDay() int64 {
	t1 := time.Now().Year()  //年
	t2 := time.Now().Month() //月
	t3 := time.Now().Day()   //日
	currentTimeData := time.Date(t1, t2, t3, 0, 0, 0, 0, time.Local)
	return currentTimeData.Unix()
}

func GenRateFail2(succ, fail int) float32 {
	sum := succ + fail
	if sum == 0 {
		return 0
	}
	return float32(fail) / float32(sum)
}

func FirstOverRateFailThd(oldhis, newhis *table.History, tp int) bool {
	if (oldhis != nil) && (oldhis.Inform != nil) && (*oldhis.Inform != 0) {
		return false
	}
	if newhis == nil {
		return false
	}
	failThd := 0
	rateFailThd := float32(0)
	if tp == Notify_Type_Sms {
		failThd = config.GConfig.Sms.FailThreshold
		rateFailThd = config.GConfig.Sms.FailRateThreshold
	} else {
		failThd = config.GConfig.Email.FailThreshold
		rateFailThd = config.GConfig.Email.FailRateThreshold
	}
	if failThd == 0 {
		return false
	}
	newSum := 0
	if newhis.DaySucc != nil {
		newSum += *newhis.DaySucc
	}
	if newhis.DayFail != nil {
		newSum += *newhis.DayFail
	}

	if newSum < failThd {
		return false
	}
	if newhis.RateFail == nil {
		return false
	}
	if *newhis.RateFail > rateFailThd {
		return true
	}
	return false
}

func SortMap(m map[string]interface{}) ([]string, []string) {
	length := len(m)
	arr := make([]string, 0, length)
	for k, _ := range m {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	values := make([]string, 0, length)
	for i := 0; i < len(arr); i++ {
		values = append(values, fmt.Sprintf("%v", m[arr[i]]))
	}
	return arr, values
}

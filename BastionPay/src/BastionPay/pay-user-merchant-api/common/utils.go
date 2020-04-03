package common

import (
	. "BastionPay/bas-base/log/zap"
	"bytes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"go.uber.org/zap"
	html "html/template"
	"os"
	"regexp"
	"strings"
	text "text/template"
	"time"
)

const (
	CompanyName = "BASTIONPAY"
)

type AwsConfig struct {
	Region      string
	AccessKeyId string
	SecretKey   string
}

func ParseHtmlTemplate(tpl string, data interface{}) (string, string, error) {
	var (
		err error
		t   *html.Template
	)

	t, err = html.ParseFiles(tpl)
	if err != nil {
		return "", "", err
	}

	var tplBuf bytes.Buffer
	if err := t.Execute(&tplBuf, data); err != nil {
		return "", "", err
	}

	result := tplBuf.String()

	return getTemplateTitle(result), result, nil
}

func getTemplateTitle(result string) string {
	r, _ := regexp.Compile(`<title>(.*)<\/title>`)

	title := r.FindString(result)
	rr, _ := regexp.Compile(`<[^>]*>`)
	return rr.ReplaceAllString(title, "")
}

func ParseTextTemplate(tpl string, data interface{}) (string, error) {
	var (
		err error
		t   *text.Template
	)

	t, err = text.ParseFiles(tpl)
	if err != nil {
		return "", err
	}

	var tplBuf bytes.Buffer
	if err := t.Execute(&tplBuf, data); err != nil {
		return "", err
	}

	result := tplBuf.String()

	return result, nil
}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		return true, nil
	}
}

func NowTimestamp() int64 {
	return time.Now().UnixNano() / (1000 * 1000)
}

func GetAppClaims(ctx iris.Context) (*AppClaims, error) {
	appClaimsInterface := ctx.Values().Get("app_claims")
	if appClaimsInterface == nil {
		return nil, errors.New("ERR_APP_CLAIMS")
	}

	return appClaimsInterface.(*AppClaims), nil
}

func GetUserIdFromCtx(ctx iris.Context) uint {
	appClaims, err := GetAppClaims(ctx)
	if err != nil {
		return 0
	}
	if appClaims.Safe {
		return appClaims.UserId
	} else {
		return 0
	}
}

func GetUserIdFromCtxUnsafe(ctx iris.Context) (*AppClaims, error) {
	appClaims, err := GetAppClaims(ctx)
	if err != nil {
		return new(AppClaims), err
	}

	return appClaims, err
}

func InArray(arr interface{}, value interface{}) bool {
	switch arr.(type) {
	case []string:
		arr = arr.([]string)
		if _, ok := value.(string); !ok {
			return false
		}
		for _, v := range arr.([]string) {
			if v == value.(string) {
				return true
			}
		}
		break
	case []uint:
		arr = arr.([]uint)
		if _, ok := value.(uint); !ok {
			return false
		}
		for _, v := range arr.([]uint) {
			if v == value.(uint) {
				return true
			}
		}
		break
	case []int:
		arr = arr.([]int)
		if _, ok := value.(int); !ok {
			return false
		}
		for _, v := range arr.([]int) {
			if v == value.(int) {
				return true
			}
		}
		break
	case []int64:
		arr = arr.([]int64)
		if _, ok := value.(int64); !ok {
			return false
		}
		for _, v := range arr.([]int64) {
			if v == value.(int64) {
				return true
			}
		}
		break
	case []float64:
		arr = arr.([]float64)
		if _, ok := value.(float64); !ok {
			return false
		}
		for _, v := range arr.([]float64) {
			if v == value.(float64) {
				return true
			}
		}
		break
	case []float32:
		arr = arr.([]float32)
		if _, ok := value.(float32); !ok {
			return false
		}
		for _, v := range arr.([]float32) {
			if v == value.(float32) {
				return true
			}
		}
		break
	}

	return false
}

func Substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return ""
	}

	if end < 0 || end > length {
		return ""
	}

	return string(rs[start:end])
}

func SecretEmail(email string) string {
	if email == "" {
		return ""
	}

	atAt := strings.IndexAny(email, "@")
	ename := Substr(email, 0, atAt)

	return Substr(ename, 0, 1) + "***" + Substr(ename, len(ename)-1, len(ename)) + Substr(email, atAt, len(email))
}

func SecretPhone(phone string) string {
	if phone == "" {
		return ""
	}

	return Substr(phone, 0, 3) + "***" + Substr(phone, len(phone)-4, len(phone))
}

func GetRealIp(ctx iris.Context) string {
	if ctx == nil {
		return ""
	}
	var (
		ips string
	)
	ZapLog().With(zap.String("X-Real-IP", ctx.GetHeader("X-Real-IP")),
		zap.String("X-Real-IP", ctx.GetHeader("X-Real-IP")),
		zap.String("X-Real-IP", ctx.GetHeader("X-Real-IP"))).
		Info("Get IP from header")
		//	glog.Info("Get IP from header X-Real-IP", ctx.GetHeader("X-Real-IP"))
		//	glog.Info("Get IP from header X-Forwarded-For", ctx.GetHeader("X-Forwarded-For"))
		//	glog.Info("Get IP from ctx.RemoteAddr()", ctx.RemoteAddr())

	ips = ctx.GetHeader("X-Forwarded-For")
	if ips == "" {
		ips = ctx.GetHeader("X-Real-IP")
	}
	if ips == "" {
		ips = ctx.RemoteAddr()
	}

	return strings.Split(ips, ",")[0]
}

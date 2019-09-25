package common

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Tools struct {
}

var (
	t    *Tools
	once sync.Once
)

/**
 * 返回单例实例
 * @method New
 */
func New() *Tools {
	once.Do(func() { //只执行一次
		t = &Tools{}
	})

	return t
}

/**
 * md5 加密
 * @method MD5
 * @param  {[type]} data string [description]
 */
func (t *Tools) MD5(data string) string {
	m := md5.New()
	io.WriteString(m, data)

	return fmt.Sprintf("%x", m.Sum(nil))
}

/**
 * string转换int
 * @method parseInt
 * @param  {[type]} b string        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) ParseInt(b string, defInt int64) int64 {
	id, err := strconv.ParseInt(b, 10, 64)
	if err != nil {
		return defInt
	} else {
		return id
	}
}

/**
 * 结构体转换成map对象
 * @method func
 * @param  {[type]} t *Tools        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) GetDateNowString() string {

	return time.Now().Format("2006-01-02 15:04:05")
}

/**
 * 结构体转换成map对象
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) StructToMap(obj interface{}) map[string]interface{} {
	k := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < k.NumField(); i++ {
		data[strings.ToLower(k.Field(i).Name)] = v.Field(i).Interface()
	}

	return data
}

/**
 * 判断手机号码
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) IsMobile(mobile string) bool {

	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`

	rgx := regexp.MustCompile(reg)

	return rgx.MatchString(mobile)
}

/**
 * 验证密码
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) CheckPassword(password, metaPassword string) bool {

	return strings.EqualFold(password, metaPassword)
}

/**
 * 生成随机字符串
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) GetRandomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()+[]{}/<>;:=.,?"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

/**
 * 生成用户Redis key
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) UserRedisKey(userId int64) string {
	userKey := fmt.Sprintf("user_login_%d", userId)

	return userKey
}

/**
 * 生成用户Token
 * @method func
 * @param  {[type]} u *Utils        [description]
 * @return {[type]}   [description]
 */
func (t *Tools) GenerateUserLoginToken(userId int64) string {
	key := t.UserRedisKey(userId)
	sum := sha256.Sum256([]byte(key))
	token := fmt.Sprintf("%x", sum)

	return token
}



func (t *Tools) SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)

	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}

		data = append(data, d)
	}

	return strings.ToLower(string(data[:]))
}

func (t *Tools) GetNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (t *Tools) GetAcceptLanguage(acceptLanguage string) string {
	language := "zh-CN"

	lang := strings.Split(acceptLanguage, ";")
	if len(lang) >= 1 {
		langs := strings.Split(lang[0], ",")
		language = langs[0]
	}

	return language
}



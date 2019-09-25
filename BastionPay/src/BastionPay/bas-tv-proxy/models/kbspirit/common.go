package kbspirit

import (
	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-tv-proxy/api"
	"encoding/binary"
	"github.com/mozillazg/go-pinyin"
	"go.uber.org/zap"
	"runtime/debug"
	"strings"
	"time"
	"unicode/utf16"
)

// FirstLetter 获取拼音首字母组成的字符串
func FirstLetter(s string) []string {
	s = strings.Replace(s, " ", "", -1)
	args := pinyin.NewArgs()
	args.Style = pinyin.FirstLetter // 首字母模式
	args.Heteronym = true           // 开启多音字
	pinyins := pinyin.Pinyin(s, args)
	result := make([]string, 0)
	for i, firstLetters := range pinyins {
		var null struct{}
		// 去重
		m := make(map[string]struct{})
		for _, firstLetter := range firstLetters {
			m[firstLetter] = null
		}
		// 非汉字处理
		if len(firstLetters) == 0 {
			m[string([]rune(s)[i:i+1])] = null
		}

		resultLen := len(result)
		if resultLen == 0 {
			for firstLetter, _ := range m {
				result = append(result, firstLetter)
			}
		} else {
			for firstLetter, _ := range m {
				for i := 0; i < resultLen; i++ {
					result = append(result, result[i]+firstLetter)
				}
			}
			// 结果截取
			result = result[resultLen:]
		}
	}
	return result
}

// Reverse 字符串反转
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Substrs 获取字符串s的所有不重复连续子串,包含本身
func Substrs(s string) []string {
	sLen := len(s)
	var null struct{}
	m := make(map[string]struct{})
	for step := 1; step <= sLen; step++ {
		for i := 0; i+step <= sLen; i++ {
			m[s[i:i+step]] = null
		}
	}

	result := make([]string, 0)
	for substr, _ := range m {
		result = append(result, substr)
	}
	return result
}

func Subrunes(s []rune) []string {
	var null struct{}
	m := make(map[string]struct{})

	for len(s) > 0 {
		bhave := 0
		for mapstr, _ := range m {
			if strings.HasPrefix(mapstr, string(s)) {
				bhave = 1
				break
			}
		}
		if 0 == bhave {
			m[string(s)] = null
		}

		s = s[1:]
	}

	result := make([]string, 0)
	for substr, _ := range m {
		result = append(result, substr)
	}
	return result
}

func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func DeleteNotLetter(name string) (string, int) {
	// name = strings.ToUpper(name)
	// num := 0
	// by := []byte(name)
	// if len(name) > 0{
	// 	for index,data := range by{
	// 		if data>='A' && data<='Z'{
	// 			num = index
	// 			break
	// 		}
	// 	}
	// }
	// return string(by[num:]),num

	name = strings.ToUpper(name)
	num := 0
	by := []byte(name)
	tembys := make([]byte, 0)
	if len(name) > 0 {
		for _, data := range by {
			if data >= 'A' && data <= 'Z' {
				tembys = append(tembys, data)
			} else {
				num = 1
			}
		}
	}
	if num == 0 {
		return name, 0
	}
	return string(tembys[:]), num
}

//键盘宝排
type JPBShuJuList []*api.JPBShuJu

func (list JPBShuJuList) Len() int {
	return len(list)
}

func (list JPBShuJuList) Less(i, j int) bool {
	return list[i].GetDaiMa() < list[j].GetDaiMa()
}

func (list JPBShuJuList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func Half(input string) string {
	bFull := false
	for _, inputu := range utf16.Encode([]rune(input)) {
		if inputu >= 0xFF21 && inputu <= 0xFF3A {
			bFull = true
			break
		}
	}
	if bFull {
		temprune := make([]rune, 0)
		for _, indexdata := range []rune(input) {
			temu := utf16.Encode([]rune(string(indexdata)))
			if len(temu) == 1 {
				if temu[0] >= 0xFF21 && temu[0] <= 0xFF3A {
					uu := uint8((temu[0] & 0x00FF) + 0x0020)
					temprune = append(temprune, rune(uu))
					// logger.Info("common Half input:%v",input)
				} else {
					temprune = append(temprune, indexdata)
				}
			} else {
				ZapLog().Info(" common Half else input:" + input)
			}
		}
		return string(temprune)
	} else {
		return input
	}
}

func handlePanic(exit bool) {
	if e := recover(); e != nil {
		ZapLog().Error("panic", zap.Any("msg", e))
		ZapLog().Error("call stack:\n" + string(debug.Stack()))
		if exit {
			time.Sleep(time.Second)
			panic("exit after recover")
		}
	}
}

func RemoveDuplicate(data []*api.JPBShuJu) []*api.JPBShuJu {
	found := make(map[string]bool)
	j := 0
	for i, d := range data {
		if !found[d.GetDaiMa()] {
			found[d.GetDaiMa()] = true
			data[j] = data[i]
			j++
		}
	}
	return data[:j]
}

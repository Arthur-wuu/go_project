package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
)

// Scan input to string with '\n' with end line
func ScanLine() string {
	var c byte
	var err error
	var b []byte
	for err == nil {
		_, err = fmt.Scanf("%c", &c)
		if c != '\n' && c != '\r' {
			b = append(b, c)
		} else {
			break
		}
	}
	return string(b)
}

func GetMd5Text(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	sum := h.Sum(nil)

	return hex.EncodeToString(sum)
}

func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

const stringSpace = " "

func space(level int) string {
	var out string
	for i := 0; i < level; i++ {
		out += stringSpace
	}
	return out
}

func FieldTag(v interface{}, level int) string {
	if v == nil {
		return ""
	}

	return fieldFormat(reflect.ValueOf(v), level)
}

func fieldFormat(vv reflect.Value, level int) string {
	t := vv.Type()
	v := vv
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = vv.Elem()
	}

	var out string
	if t.Kind() == reflect.Struct {
		out += space(level) + "{\n"
		n := t.NumField()
		for i := 0; i < n; i++ {
			tagJson := space(level+1) + t.Field(i).Tag.Get("json")
			if tagJson == "-" {
				continue
			}
			out += space(level+1) + tagJson + "--(" + t.Field(i).Tag.Get("doc") + ") "
			out += fieldFormat(v.Field(i), level+1)
			out += "\n"
		}
		out += space(level) + "}"
	} else if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		n := v.Len()
		out += space(level) + "[\n"
		for i := 0; i < n; i++ {
			rs := v.Index(i)
			out += space(level) + fieldFormat(rs, level+1)
			out += "\n"
		}
		out += space(level) + "]"
	} else if t.Kind() == reflect.Map {
		ks := v.MapKeys()
		out += "{\n"
		for i := 0; i < len(ks); i++ {
			out += fieldFormat(ks[i], level+1)
			out += ":"
			key := v.MapIndex(ks[i])
			out += fieldFormat(key, level+1)
		}
		out += "\n}"
	} else {
		out = v.Type().String()
	}

	return out
}

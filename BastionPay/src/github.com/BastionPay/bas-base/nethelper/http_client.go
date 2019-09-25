package nethelper

import (
	"bytes"
	"io/ioutil"
	"net/http"
	l4g "github.com/alecthomas/log4go"
)

// Do a post to Http server
func CallToHttpServer(addr string, path string, body string) (int, string, error) {
	url := addr + path
	contentType := "application/json;charset=utf-8"

	b := []byte(body)
	b2 := bytes.NewBuffer(b)

	resp, err := http.Post(url, contentType, b2)
	if err != nil {
		l4g.Error("Post failed: %s", err.Error())
		return -1, "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l4g.Error("Read failed: %s", err.Error())
		return -1, "", err
	}

	return resp.StatusCode, string(content), nil
}
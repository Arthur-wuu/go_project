package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type (
	Source struct {
		ApiKey    string `json:"user_key"`
		Message   string `json:"message"`
		Signature string `json:"signature"`
	}
)

func Post(url string, body interface{}) ([]byte, error) {

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	client := http.Client{Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ResponseHeaderTimeout: time.Second * 30,
	}}

	resp, err := client.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

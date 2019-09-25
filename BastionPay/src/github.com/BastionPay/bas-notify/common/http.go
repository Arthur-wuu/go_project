package common

import (
	"io"
	"net/http"
	"io/ioutil"
	"errors"
)

func HttpSend(url string, body io.Reader, method string, headers map[string]string) ([]byte, error) {
	if len(method) == 0 {
		method = "GET"
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Content-Type", "application/json")
	for k,v := range headers {
		//fmt.Println(k, v)
		req.Header.Set(k, v)
	}
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
		return nil,err
	}

	if len(content) == 0 {
		return nil, errors.New("nil resp")
	}
	return content, nil
}

func PostForm(url string, form map[string][]string) ([]byte, error) {
	resp, err := http.PostForm(url, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	if len(content) == 0 {
		return nil, errors.New("nil resp")
	}
	return content, nil
}
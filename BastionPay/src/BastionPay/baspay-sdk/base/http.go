package base

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

func HttpSend(url string, body io.Reader, method string, headers map[string]string) ([]byte, error) {
	if len(method) == 0 {
		method = "GET"
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		//fmt.Println(k, v)
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return content, errors.New(resp.Status)
	}

	if len(content) == 0 {
		return nil, err
	}
	return content, nil
}

//
func HttpSendSer(url string, body io.Reader, method string, headers map[string]string) ([]byte, error) {
	if len(method) == 0 {
		method = "GET"
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for k, v := range headers {
		//fmt.Println(k, v)

		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return content, errors.New(resp.Status)
	}

	if len(content) == 0 {
		return nil, err
	}
	return content, nil
}

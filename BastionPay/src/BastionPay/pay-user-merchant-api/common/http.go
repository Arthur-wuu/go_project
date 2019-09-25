package common

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"errors"
	"io"
)

const (
	HttpTimeout = 10
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
	for k,v := range headers {

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

type Http struct {
	client *http.Client
}

func NewHttp() *Http {
	return &Http{
		client: &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Second*HttpTimeout)
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(time.Second * HttpTimeout))
					return conn, nil
				},
				ResponseHeaderTimeout: time.Second * HttpTimeout,
			},
		},
	}
}

func (h *Http) Get(url string) ([]byte, error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return result, nil
}

package common

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	HttpTimeout = 10
)

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

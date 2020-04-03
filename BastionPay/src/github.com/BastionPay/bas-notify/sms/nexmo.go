package sms

import (
	"BastionPay/bas-notify/common"
	"BastionPay/bas-notify/config"
	"encoding/json"
	"fmt"
)

type ReqNexmo struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	From      string `json:"from"`
	To        string `json:"to"`
	Text      string `json:"text"`
	Type      string `json:"type"`
}

type ResNexmo struct {
	MessageCount string `json:"message-count"`
	Messages     []ResNexmoMessages
}

type ResNexmoMessages struct {
	To               string `json:"to"`
	MessageId        string `json:"message-id"`
	Status           string `json:"status"`
	RemainingBalance string `json:"remaining-balance"`
	MessagePrice     string `json:"message-price"`
	Network          string `json:"network"`
	AccountRef       string `json:"account-ref"`
	ErrorText        string `json:"error-text"`
}

func (this *SmsMgr) DirectSendNexmo(body, phone string, senderId *string) error {
	//if len(phone) > 4 {
	//	phone = phone[2:]
	//}

	form := make(map[string][]string)
	form["api_key"] = []string{config.GConfig.Nexmo.ApiKey}
	form["api_secret"] = []string{config.GConfig.Nexmo.ApiSecret}
	form["to"] = []string{phone}
	form["text"] = []string{body}
	form["type"] = []string{"unicode"}

	if senderId != nil {
		form["from"] = []string{*senderId}
	} else {
		form["from"] = []string{config.GConfig.Nexmo.DefaultSender}
	}

	resData, err := common.PostForm(config.GConfig.Nexmo.Url, form)
	if err != nil {
		return err
	}

	res := new(ResNexmo)
	if err := json.Unmarshal(resData, res); err != nil {
		return err
	}
	if res.Messages == nil || len(res.Messages) <= 0 {
		return fmt.Errorf("no response Messages")
	}

	if res.Messages[0].Status != "0" {
		return fmt.Errorf("%s-%s-%s", res.Messages[0].Status, res.Messages[0].MessageId, res.Messages[0].ErrorText)
	}
	return nil
}

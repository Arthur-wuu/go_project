package api

import "github.com/shopspring/decimal"

type NotifyRequest struct {
	Id             *int             `json:"id,omitempty"  `
	ActivityUuid   *string          `json:"activity_uuid,omitempty"  `
	RedId          *string          `json:"red_uuid,omitempty"  `
	AppId          *string          `json:"app_id,omitempty"`
	UserId         interface{}      `json:"user_id,omitempty"`
	CountryCode    *string          `json:"country_code,omitempty"`
	Phone          *string          `json:"phone,omitempty" `
	Symbol         *string          `json:"symbol,omitempty" `
	Coin           *decimal.Decimal `json:"coin,omitempty"`
	SponsorAccount *string          `json:"sponsor_account,omitempty" `
	ApiKey         *string          `json:"api_key,omitempty" `
	OffAt          *int64           `json:"off_at,omitempty" `
	Lang           *string          `json:"language,omitempty" `
	//TransferFlag 	*int       `json:"transfer_flag,omitempty"`
}

type NotifyResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

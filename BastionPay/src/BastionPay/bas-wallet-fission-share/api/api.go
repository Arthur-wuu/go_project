package api

import "github.com/shopspring/decimal"

type (
	//创建订单, 查询
	ActivityAdd struct {
		Id           *int             `valid:"optional" json:"id,omitempty"`
		Name         *string          `valid:"required" json:"name,omitempty"`
		Sponsor      *string          `valid:"required" json:"sponsor,omitempty"`
		Merchant     *string          `valid:"optional" json:"merchant,omitempty"`
		Country      *string          `valid:"required" json:"country,omitempty"`
		Language     *string          `valid:"required" json:"language,omitempty"`
		Mode         *int             `valid:"required" json:"mode,omitempty"`
		Symbol       *string          `valid:"required" json:"symbol,omitempty"`
		TotalCoin    *decimal.Decimal `valid:"required" json:"total_coin,omitempty"`
		TotalRed     *int64           `valid:"required" json:"total_red,omitempty"`
		RedExpire    *int64           `valid:"required" json:"red_expire,omitempty"`
		RobRedExpire *int64           `valid:"required" json:"rob_red_expire,omitempty"`
		TotalRob     *int64           `valid:"required" json:"total_rob,omitempty"`
		OnAt         *int64           `valid:"required" json:"on_at,omitempty"`
		OffAt        *int64           `valid:"required" json:"off_at,omitempty"`
		Valid        *int             `valid:"optional"  json:"valid,omitempty"`
		Precision    *int32           `valid:"required" json:"precision,omitempty"`
	}

	Activity struct {
		Id           *int             `valid:"required" json:"id,omitempty"`
		Name         *string          `valid:"optional" json:"name,omitempty"`
		Sponsor      *string          `valid:"optional" json:"sponsor,omitempty"`
		Merchant     *string          `valid:"optional" json:"merchant,omitempty"`
		Country      *string          `valid:"optional" json:"country,omitempty"`
		Language     *string          `valid:"optional" json:"language,omitempty"`
		Mode         *int             `valid:"optional" json:"mode,omitempty"`
		Symbol       *string          `valid:"optional" json:"symbol,omitempty"`
		TotalCoin    *decimal.Decimal `valid:"optional" json:"total_coin,omitempty"`
		TotalRed     *int64           `valid:"optional" json:"total_red,omitempty"`
		RedExpire    *int64           `valid:"optional" json:"red_expire,omitempty"`
		RobRedExpire *int64           `valid:"optional" json:"rob_red_expire,omitempty"`
		TotalRob     *int64           `valid:"optional" json:"total_rob,omitempty"`
		OnAt         *int64           `valid:"optional" json:"on_at,omitempty"`
		OffAt        *int64           `valid:"optional" json:"off_at,omitempty"`
		Valid        *int             `valid:"optional"  json:"valid,omitempty"`
		Precision    *int32           `valid:"optional" json:"precision,omitempty"`
	}

	ActivityList struct {
		Id           *int             `valid:"optional" json:"id,omitempty"`
		Name         *string          `valid:"optional" json:"name,omitempty"`
		Sponsor      *string          `valid:"optional" json:"sponsor,omitempty"`
		Merchant     *string          `valid:"optional" json:"merchant,omitempty"`
		Country      *string          `valid:"optional" json:"country,omitempty"`
		Language     *string          `valid:"optional" json:"language,omitempty"`
		Mode         *int             `valid:"optional" json:"mode,omitempty"`
		Symbol       *string          `valid:"optional" json:"symbol,omitempty"`
		TotalCoin    *decimal.Decimal `valid:"optional" json:"total_coin,omitempty"`
		TotalRed     *int64           `valid:"optional" json:"total_red,omitempty"`
		RedExpire    *int64           `valid:"optional" json:"red_expire,omitempty"`
		RobRedExpire *int64           `valid:"optional" json:"rob_red_expire,omitempty"`
		OnAt         *int64           `valid:"optional" json:"on_at,omitempty"`
		OffAt        *int64           `valid:"optional" json:"off_at,omitempty"`
		Valid        *int             `valid:"optional"  json:"valid,omitempty"`
		Page         int64            `valid:"optional" json:"page,omitempty"`
		Size         int64            `valid:"optional" json:"size,omitempty"`
	}

	RedAdd struct {
		ActivityId *int             `valid:"required" json:"activity_id,omitempty"`
		Symbol     *string          `valid:"optional" json:"-"`
		TotalCoin  *decimal.Decimal `valid:"optional" json:"-"`
		TotalSize  *int64           `valid:"optional" json:"-"`
		Mode       *int             `valid:"optional" json:"-"`
		SrcUid     *int64           `valid:"required" json:"src_uid,omitempty"`
		ExpireAt   *int64           `valid:"optional" json:"-"`
		RobExpire  *int64           `valid:"optional" json:"-"`
		OrderId    *string          `valid:"optional" json:"order_id,omitempty"`
		Valid      *int             `valid:"optional"  json:"valid,omitempty"`
	}

	RedList struct {
		Id         *int    `valid:"optional" json:"id,omitempty"`
		ActivityId *int    `valid:"optional" json:"activity_id,omitempty"`
		Symbol     *string `valid:"optional" json:"symbol,omitempty"`
		Mode       *int    `valid:"optional" json:"mode,omitempty"`
		SrcUid     *int64  `valid:"optional" json:"src_uid,omitempty"`
		ExpireAt   *int64  `valid:"optional" json:"expire_at,omitempty"`
		OrderId    *string `valid:"optional" json:"order_id,omitempty"`
		Valid      *int    `valid:"optional"  json:"valid,omitempty"`
		Page       int64   `valid:"optional" json:"page,omitempty"`
		Size       int64   `valid:"optional" json:"size,omitempty"`
	}

	RobberList struct {
		Id          *int             `valid:"optional" json:"id,omitempty"`
		RedId       *int             `valid:"required" json:"red_id,omitempty"`
		CountryCode *string          `valid:"required" json:"country_code,omitempty"`
		Phone       *string          `valid:"required" json:"phone,omitempty"`
		Symbol      *string          `valid:"optional" json:"symbol,omitempty"`
		Coin        *decimal.Decimal `valid:"optional" json:"coin,omitempty"`
		SrcUrl      *string          `valid:"required" json:"src_url,omitempty"`
		SrcUid      *int64           `valid:"optional" json:"src_uid,omitempty"`
		ExpireAt    *int64           `valid:"optional" json:"expire_at,omitempty"`
		Page        int64            `valid:"optional" json:"page,omitempty"`
		Size        int64            `valid:"optional" json:"size,omitempty"`
	}

	RobberAdd struct {
		RedId       *int             `valid:"required" json:"red_id,omitempty"`
		CountryCode *string          `valid:"required" json:"country_code,omitempty"`
		Phone       *string          `valid:"required" json:"phone,omitempty"`
		Symbol      *string          `valid:"optional" json:"-"`
		Coin        *decimal.Decimal `valid:"optional" json:"-"`
		SrcUrl      *string          `valid:"required" json:"src_url,omitempty"`
		SrcUid      *int64           `valid:"optional" json:"src_uid,omitempty"`
		ExpireAt    *int64           `json:"optional" json:"-"`
	}

	SloganAdd struct {
		ActivityId *int    `valid:"required" json:"activity_id,omitempty"`
		TitleId    *int    `valid:"required" json:"title_id,omitempty"`
		Text       *string `valid:"required" json:"text,omitempty"`
		Valid      *int    `valid:"optional"  json:"valid,omitempty"`
	}

	Slogan struct {
		Id         *int    `valid:"required" json:"id,omitempty"`
		ActivityId *int    `valid:"optional" json:"activity_id,omitempty"`
		TitleId    *int    `valid:"optional" json:"title_id,omitempty"`
		Text       *string `valid:"optional" json:"text,omitempty"`
		Valid      *int    `valid:"optional"  json:"valid,omitempty"`
	}

	SloganList struct {
		Id         *int    `valid:"optional" json:"id,omitempty"`
		ActivityId *int    `valid:"optional" json:"activity_id,omitempty"`
		TitleId    *int    `valid:"optional" json:"title_id,omitempty"`
		Text       *string `valid:"optional" json:"text,omitempty"`
		Valid      *int    `valid:"optional"  json:"valid,omitempty"`
		Page       int64   `valid:"optional" json:"page,omitempty"`
		Size       int64   `valid:"optional" json:"size,omitempty"`
	}

	SloganGets struct {
		ActivityId *int `valid:"required" json:"activity_id,omitempty"`
	}
)

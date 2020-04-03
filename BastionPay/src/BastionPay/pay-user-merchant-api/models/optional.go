package models

import (
	"encoding/json"

	"BastionPay/pay-user-merchant-api/db"
)

var GOptionalModel OptionalModel

type OptionalModel struct {
}

func (o *OptionalModel) Get(userId uint) ([]string, error) {
	var (
		markets  []string
		optional = Optional{}
		err      error
	)

	if err = db.GDbMgr.Get().Where(&Optional{UserId: userId}).First(&optional).Error; err != nil {
		return nil, err
	}

	if optional.Optional == "" {
		return markets, nil
	}

	if err = json.Unmarshal([]byte(optional.Optional), &markets); err != nil {
		return nil, err
	}

	return markets, nil
}

func (o *OptionalModel) Create(userId uint, markets []string) (uint, error) {
	var (
		mkStr []byte
		err   error
	)

	mkStr, err = json.Marshal(markets)
	if err != nil {
		return 0, err
	}

	optional := Optional{
		UserId:   userId,
		Optional: string(mkStr),
	}

	if err := db.GDbMgr.Get().Save(&optional).Error; err != nil {
		return 0, err
	}

	return optional.ID, nil
}

func (o *OptionalModel) Delete(userId uint) error {
	err := db.GDbMgr.Get().Where(Optional{UserId: userId}).Delete(&Optional{}).Error
	if err != nil {
		return err
	}
	return nil
}

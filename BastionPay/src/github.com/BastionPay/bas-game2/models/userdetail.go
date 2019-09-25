package models

import (
	//"BastionPay/bas-game/common"
	"BastionPay/bas-game2/db"
	"BastionPay/bas-game2/models/table"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type (
	UserDetailMsg struct {
        Id                  *int64         `valid:"required"     json:"ID"`
        AccId               *int           `valid:"optional"     json:"ACC_ID"`
        UserId              *int64         `valid:"required"     json:"USER_ID"`
        Assets              *string        `valid:"optional"     json:"ASSETS"`
        Amount              *float64       `valid:"optional"     json:"AMOUNT"`
        AfterAmount         *float64       `valid:"optional"     json:"AFTER_AMOUNT"`
        Direction           *int32         `valid:"optional"     json:"DIRECTION"`
        Remark              *string        `valid:"optional"     json:"REMARK"`
        CreateTmie          *time.Duration `valid:"optional"     json:"CREATE_TIME"`
        LastUpdateTmie      *time.Duration `valid:"optional"     json:"LAST_UPDATE_TIME"`
        OrderId             *int64         `valid:"optional"     json:"ORDER_ID"`
        FrozenAmount        *float64       `valid:"optional"     json:"FROZEN_AMOUNT"`
        AfterFrozenAmount   *float64       `valid:"optional"     json:"AFTER_FROZEN_AMOUNT"`
}
)


func (this *UserDetailMsg) Gets() (*table.UserDetail, error) {
	//model := &table.UserDetail{
	//	Id:                this.Id,
	//	AccId:             this.AccId,
	//	UserId:            this.UserId,
	//	Assets:            this.Assets,
	//	Amount:            this.Amount,
	//	AfterAmount:       this.AfterAmount,
	//	Direction:         this.Direction,
	//	Remark:            this.Remark,
	//	CreateTmie:        this.CreateTmie,
	//	LastUpdateTmie:    this.LastUpdateTmie,
	//	OrderId:           this.OrderId,
	//	FrozenAmount:      this.FrozenAmount,
	//	AfterFrozenAmount: this.AfterFrozenAmount,
	//}

	userDetail := new(table.UserDetail)
	 err := db.GDbMgr.Get().Where("USER_ID=?",30).Last(&userDetail).Error

	if err == gorm.ErrRecordNotFound {
		fmt.Println("err gorm***")
		return nil,nil
	}
	return userDetail, err
}


package access

import (
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
)

type (
	UserLoginLog struct {
		UserId int64
		Ip     string
	}
)

var (
	Tools = common.New()
)

func (this *UserLoginLog) SaveLog() error {
	model := &models.UserLoginLog{
		UserId:    this.UserId,
		Ip:        this.Ip,
		CreatedAt: Tools.GetDateNowString(),
	}

	return models.DB.Save(model).Error
}

func (this *UserLoginLog) GetLoginLog() (*models.UserLoginLog, error) {
	var info models.UserLoginLog

	err := models.DB.Where("user_id = ?", this.UserId).
		Order("created_at DESC").First(&info).Error

	if err != nil {
		return nil, err
	}

	return &info, nil
}

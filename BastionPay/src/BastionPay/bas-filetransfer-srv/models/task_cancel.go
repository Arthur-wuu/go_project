package models

import (
	"BastionPay/bas-filetransfer-srv/db"
)

type TaskCancelInfo struct {
	Order_id string `form:"order_id" valid:"required"`
}

func (this *TaskCancelInfo) Cancel() error {
	_, err := db.GRedis.Do("del", EXPORT_Status_KeyPrefix+this.Order_id)
	return err
}

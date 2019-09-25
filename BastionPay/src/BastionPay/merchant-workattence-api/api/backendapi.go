package api

import "github.com/shopspring/decimal"

type (
	BkAccountMapAdd struct {
		AccId      *int    `valid:"required" json:"account_id,omitempty"` //bastionpay account id
		StaffId    *string `valid:"required" json:"staff_id,omitempty"`   //考勤机员工id
		CorpId     *string `valid:"required" json:"corp_id,omitempty"`    //企业id
		Phone      *string `valid:"optional" json:"phone,omitempty"`      //员工手机号
		Name       *string `valid:"optional" json:"name,omitempty"`       //员工姓名
		DeleteFlag *int    `valid:"optional" json:"valid,omitempty"`      //删除标志0（删除）| 1（正常）
	}

	BkAccountMapUpdate struct {
		Id         *int    `valid:"required"	json:"id,omitempty"`         //bastionpay account id
		AccId      *int    `valid:"optional" json:"account_id,omitempty"` //bastionpay account id
		StaffId    *string `valid:"optional" json:"staff_id,omitempty"`   //考勤机员工id
		CorpId     *string `valid:"optional" json:"corp_id,omitempty"`    //企业id
		Phone      *string `valid:"optional" json:"phone,omitempty"`      //员工手机号
		Name       *string `valid:"optional" json:"name,omitempty"`       //员工姓名
		DeleteFlag *int    `valid:"optional" json:"valid,omitempty"`      //删除标志0（删除）| 1（正常）
	}

	BkAccountMapList struct {
		Id         *int    `valid:"optional"	json:"id,omitempty"`         //bastionpay account id
		AccId      *int    `valid:"optional" json:"account_id,omitempty"` //bastionpay account id
		StaffId    *string `valid:"optional" json:"staff_id,omitempty"`   //考勤机员工id
		CorpId     *string `valid:"optional" json:"corp_id,omitempty"`    //企业id
		Phone      *string `valid:"optional" json:"phone,omitempty"`      //员工手机号
		Name       *string `valid:"optional" json:"name,omitempty"`       //员工姓名
		DeleteFlag *int    `valid:"optional" json:"valid,omitempty"`      //删除标志0（删除）| 1（正常）
		CreatedAt  *int64  `valid:"optional" json:"created_at,omitempty"`
		Page       int64   `valid:"optional" json:"page,omitempty"`
		Size       int64   `valid:"optional" json:"size,omitempty"`
	}

	BkAccountMap struct {
		Id         *int    `valid:"required"	json:"id,omitempty"`         //bastionpay account id
		AccId      *int    `valid:"optional" json:"account_id,omitempty"` //bastionpay account id
		StaffId    *string `valid:"optional" json:"staff_id,omitempty"`   //考勤机员工id
		CorpId     *string `valid:"optional" json:"corp_id,omitempty"`    //企业id
		Phone      *string `valid:"optional" json:"phone,omitempty"`      //员工手机号
		Name       *string `valid:"optional" json:"name,omitempty"`       //员工姓名
		DeleteFlag *int    `valid:"optional" json:"valid,omitempty"`      //删除标志0（删除）| 1（正常）
		CreatedAt  *int64  `valid:"optional" json:"created_at,omitempty"`
	}

	RubbishClassifyAwardSend struct {
		Id *int `valid:"required"	json:"id,omitempty"`
	}

	BkRubbishClassifyAwardList struct {
		UserId       *int    `valid:"optional" json:"user_id,omitempty"`
		CompanyId    *int    `valid:"optional" json:"company_id,omitempty"`
		DepartmentId *int    `valid:"optional" json:"department_id,omitempty"`
		Name         *string `valid:"optional" json:"name,omitempty"`
		BpUid        *int    `valid:"optional" json:"bp_uid,omitempty"`
		Score        *int    `valid:"optional" json:"score,omitempty"` // 0 差 1 良 2 优
		TransferFlag *int    `valid:"optional" json:"transfer_flag,omitempty"`
		Page         int64   `valid:"optional" json:"page,omitempty"`
		Size         int64   `valid:"optional" json:"size,omitempty"`
	}

	ResponseDepartmenInfo struct {
		Id           int    `valid:"optional" json:"id,omitempty"`
		Date         string `valid:"required" json:"date,omitempty"`
		TotalNumbers int    `valid:"required" json:"total_numbers,omitempty"`
		Department   []struct {
			DepartmentId   int    `valid:"required" json:"department_id,omitempty"`
			DepartmentName string `valid:"required" json:"department_name,omitempty"`
			Numbers        int64  `valid:"required" json:"numbers,omitempty"`
			Score          int    `valid:"optional" json:"score"`
		} `valid:"required" json:"department,omitempty"`
	}

	BkRubbishClassifyList struct {
		Id           *int             `valid:"optional"	json:"id,omitempty"`
		TotalNumbers *int             `valid:"optional" json:"total_numbers,omitempty"`
		TotalCoin    *decimal.Decimal `valid:"optional" json:"total_coin,omitempty"`
		ScoreDate    *string          `valid:"optional" json:"score_date,omitempty"`
		TransferFlag *int             `valid:"optional" json:"transfer_flag,omitempty"` //发送标志0（未发送）| 1（已发送）
		CreatedAt    *int64           `valid:"optional" json:"created_at,omitempty"`
		Page         int64            `valid:"optional" json:"page,omitempty"`
		Size         int64            `valid:"optional" json:"size,omitempty"`
	}

	BkEditRubbishClassify struct {
		Id           *int             `valid:"required"	json:"id,omitempty"`
		TotalNumbers *int             `valid:"optional" json:"total_numbers,omitempty"`
		TotalCoin    *decimal.Decimal `valid:"optional" json:"total_coin,omitempty"`
		ScoreDate    *string          `valid:"optional" json:"score_date,omitempty"`
		TransferFlag *int             `valid:"optional" json:"transfer_flag,omitempty"` //发送标志0（未发送）| 1（已发送）
		CreatedAt    *int64           `valid:"optional" json:"created_at,omitempty"`
	}

	RubbishClassifyAwardSendList struct {
		Id   int   `valid:"required"	json:"id,omitempty"`
		Page int64 `valid:"optional" json:"page,omitempty"`
		Size int64 `valid:"optional" json:"size,omitempty"`
	}

	StaffMotivationList struct {
		Datetime *string `valid:"optional"	json:"datetime,omitempty"`
		Page     int64   `valid:"optional" json:"page,omitempty"`
		Size     int64   `valid:"optional" json:"size,omitempty"`
	}

	OvertimeAwardList struct {
		Datetime *string `valid:"optional"	json:"datetime,omitempty"`
		Page     int64   `valid:"optional" json:"page,omitempty"`
		Size     int64   `valid:"optional" json:"size,omitempty"`
	}
)

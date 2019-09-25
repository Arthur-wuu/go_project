package api

type (
	PushCheckin struct {
		Id     *string `valid:"required" json:"id,omitempty"`
		Data   *string `valid:"optional" json:"data,omitempty"`
		Ccid   *string `valid:"required" json:"ccid,omitempty"`
		Time   *string `valid:"required" json:"time,omitempty"`
		Verify *int    `valid:"required" json:"verify,omitempty"`
		Pic    *string `valid:"optional" json:"pic,omitempty"`
	}

	PullCheckin struct {
		AttenId     *string `valid:"optional" json:"atten_id,omitempty"`     //考勤ID
		AttenDevice *string `valid:"optional" json:"atten_device,omitempty"` //考勤设备ID
		StaffId     *string `valid:"optional" json:"staff_id,omitempty"`     //考勤机员工id
		Verify      *int    `valid:"optional" json:"verify,omitempty"`       //打卡方式 0(密码)，1(指纹)，2(刷卡)，15(人脸)
		CheckinAt   *int64  `valid:"optional" json:"created_at,omitempty"`   //考勤打卡时间
		Realname    *string `valid:"optional" json:"realname,omitempty"`     //员工姓名
		Departname  *string `valid:"optional" json:"departname,omitempty"`   //部门名称
	}

	PullResult struct {
		Status *int    `valid:"required" json:"status,omitempty"`
		Error  *string `valid:"required" json:"error,omitempty"`
		Data   *struct {
			Total     *string `valid:"required" json:"total,omitempty"`
			Totalpage *int    `valid:"required" json:"totalpage,omitempty"`
			Page      *int    `valid:"required" json:"page,omitempty"`
			Attendata *[]PushCheckin
		}
	}
)

type RecordListResponse struct {
	Errcode      int            `json:"errcode,omitempty"`
	Errmsg       string         `json:"errmsg,omitempty"`
	Recordresult []Recordresult `json:"recordresult,omitempty"`
}

type Recordresult struct {
	WorkDate       int64  `json:"workDate,omitempty"`
	CorpId         string `json:"corpId,omitempty"`
	CheckType      string `json:"checkType,omitempty"`
	SourceType     string `json:"sourceType,omitempty"`
	TimeResult     string `json:"timeResult,omitempty"`
	UserAddress    string `json:"userAddress,omitempty"`
	UserCheckTime  int64  `json:"userCheckTime,omitempty"`
	LocationMethod string `json:"locationMethod,omitempty"`
	DeviceId       string `json:"deviceId,omitempty"`
	IsLegal        string `json:"isLegal,omitempty"`
	LocationResult string `json:"locationResult,omitempty"`
	UserId         string `json:"userId,omitempty"`
}

//Employee Motivation
type StaffList struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    []StaffListData `json:"data,omitempty"`
}

type StaffListData struct {
	Id             *int    `json:"id,omitempty"`
	CompanyId      *int    `json:"company_id,omitempty"`
	DepartmentId   *int    `json:"department_id,omitempty"`
	CompanyName    *string `json:"company_name,omitempty"`
	DepartmentName *string `json:"department_name,omitempty"`
	Name           *string `json:"name,omitempty"`
	Sex            *int    `json:"sex,omitempty"`
	HiredAt        *int64  `json:"hired_at,omitempty"`
	BirthAt        *int64  `json:"birth_at,omitempty"`
	EeNo           *string `json:"mech_no,omitempty"`
	BpUid          *int    `json:"bp_uid,omitempty"`
	Valid          *int    `json:"valid,omitempty"`
}

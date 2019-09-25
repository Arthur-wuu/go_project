package api

type (
	FtDepartmentList struct {
		CompanyId *int64  `valid:"required" json:"company_id,omitempty" `
		Name      *string `valid:"optional" json:"name,omitempty" `
		EpeNum    *int    `valid:"optional" json:"employee_num,omitempty"`
		Vaild     *int    `valid:"optional" json:"valid,omitempty" `

		Page int64 `valid:"required" json:"page,omitempty"`
		Size int64 `valid:"optional" json:"size,omitempty"`
	}
	FtDepartmentGet struct {
		Id        *int64  `valid:"required"  json:"id,omitempty" `
		CompanyId *int64  `valid:"optional" json:"company_id,omitempty" `
		Name      *string `valid:"optional" json:"name,omitempty" `
		EpeNum    *int    `valid:"optional" json:"employee_num,omitempty"`
		Vaild     *int    `valid:"optional" json:"valid,omitempty" `
	}
	FtDepartmentGets struct {
		CompanyId *int64 `valid:"required" json:"company_id,omitempty" `
	}
	FtEmployeeList struct {
		CompanyId    *int64  `valid:"required" json:"company_id,omitempty"`
		DepartmentId *int64  `valid:"optional" json:"department_id,omitempty"`
		Name         *string `valid:"optional" json:"name,omitempty" `
		Sex          *int    `valid:"optional" json:"sex,omitempty"`
		HiredAt      *int64  `valid:"optional" json:"hired_at,omitempty"`
		BirthAt      *int64  `valid:"optional" json:"birth_at,omitempty"`
		EeNo         *string `valid:"optional" json:"ee_no,omitempty"`
		BpUid        *int64  `valid:"optional" json:"bp_uid,omitempty"`
		MechNo       *string `valid:"optional" json:"mech_no,omitempty"`
		Vaild        *int    `valid:"optional" json:"valid,omitempty" `

		Page int64 `valid:"required" json:"page,omitempty"`
		Size int64 `valid:"optional" json:"size,omitempty"`
	}
	FtEmployeeGet struct {
		Id           *int64  `valid:"optional"  json:"id,omitempty" `
		CompanyId    *int64  `valid:"optional" json:"company_id,omitempty"`
		DepartmentId *int64  `valid:"optional" json:"department_id,omitempty"`
		Name         *string `valid:"optional" json:"name,omitempty" `
		Sex          *int    `valid:"optional" json:"sex,omitempty"`
		EeNo         *string `valid:"optional" json:"ee_no,omitempty"`
		BpUid        *int64  `valid:"optional" json:"bp_uid,omitempty"`
		MechNo       *string `valid:"optional" json:"mech_no,omitempty"`
		Vaild        *int    `valid:"optional" json:"valid,omitempty" `
	}

	FtEmployeeGets struct {
		CompanyId    *int64 `valid:"required" json:"company_id,omitempty"`
		DepartmentId *int64 `valid:"optional" json:"department_id,omitempty"`
		Vaild        *int   `valid:"optional" json:"valid,omitempty" `
	}
)

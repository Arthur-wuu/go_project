package api

type (
	CompanyAdd struct {
		Name    *string `valid:"required" json:"name,omitempty" `
		Address *string `valid:"optional" json:"address,omitempty" `
		Vaild   *int    `valid:"optional" json:"valid,omitempty" `
	}

	Company struct {
		Id      *int64  `valid:"required" json:"id,omitempty" `
		Name    *string `valid:"optional" json:"name,omitempty" `
		Address *string `valid:"optional" json:"address,omitempty" `
		Vaild   *int    `valid:"optional" json:"valid,omitempty" `
	}

	CompanyList struct {
		Name    *string `valid:"optional" json:"name,omitempty" `
		Address *string `valid:"optional" json:"address,omitempty" `
		Vaild   *int    `valid:"optional" json:"valid,omitempty" `

		Page int64 `valid:"required" json:"page,omitempty"`
		Size int64 `valid:"optional" json:"size,omitempty"`
	}

	DepartmentAdd struct {
		CompanyId *int64  `valid:"required" json:"company_id,omitempty" `
		Name      *string `valid:"required" json:"name,omitempty" `
		ParentId  *int64  ` valid:"optional" json:"parent_id,omitempty"`
		EpeNum    *int    `valid:"optional" json:"employee_num,omitempty"`
		Vaild     *int    `valid:"optional" json:"valid,omitempty" `
	}

	Department struct {
		Id        *int64  `valid:"required" json:"id,omitempty" ` //加上type:int(11)后AUTO_INCREMENT无效
		CompanyId *int64  `valid:"optional" json:"company_id,omitempty" `
		Name      *string `valid:"optional" json:"name,omitempty" `
		ParentId  *int64  ` valid:"optional" json:"parent_id,omitempty"`
		EpeNum    *int    `valid:"optional" json:"employee_num,omitempty"`
		Vaild     *int    `valid:"optional" json:"valid,omitempty" `
	}

	DepartmentList struct {
		CompanyId *int64  `valid:"optional" json:"company_id,omitempty" `
		Name      *string `valid:"optional" json:"name,omitempty" `
		EpeNum    *int    `valid:"optional" json:"employee_num,omitempty"`
		Vaild     *int    `valid:"optional" json:"valid,omitempty" `

		Page int64 `valid:"required" json:"page,omitempty"`
		Size int64 `valid:"optional" json:"size,omitempty"`
	}

	EmployeeAdd struct {
		CompanyId    *int64  `valid:"required"  json:"company_id,omitempty"`
		DepartmentId *int64  `valid:"required"  json:"department_id,omitempty"`
		Name         *string `valid:"optional" json:"name,omitempty" `
		Sex          *int    `valid:"optional" json:"sex,omitempty"`
		HiredAt      *int64  `valid:"optional" json:"hired_at,omitempty"`
		BirthAt      *int64  `valid:"optional" json:"birth_at,omitempty"`
		RegularAt    *int64  `valid:"optional" json:"regular_at,omitempty"`
		EeNo         *string `valid:"optional" json:"ee_no,omitempty"`
		MechNo       *string `valid:"optional" json:"mech_no,omitempty"`
		BpUid        *int64  `valid:"optional" json:"bp_uid,omitempty"`
		Phone        *string `valid:"optional" json:"phone,omitempty" `
		Wchat        *string `valid:"optional" json:"wchat,omitempty"`
		Email        *string `valid:"optional" json:"email,omitempty"`
		Vaild        *int    `valid:"optional" json:"valid,omitempty" `
	}

	Employee struct {
		Id           *int64  `valid:"required"  json:"id,omitempty" ` //加上type:int(11)后AUTO_INCREMENT无效
		CompanyId    *int64  `valid:"optional" json:"company_id,omitempty"`
		DepartmentId *int64  `valid:"optional" json:"department_id,omitempty"`
		Name         *string `valid:"optional" json:"name,omitempty" `
		Sex          *int    `valid:"optional" json:"sex,omitempty"`
		HiredAt      *int64  `valid:"optional" json:"hired_at,omitempty"`
		BirthAt      *int64  `valid:"optional" json:"birth_at,omitempty"`
		RegularAt    *int64  `valid:"optional" json:"regular_at,omitempty"`
		EeNo         *string `valid:"optional" json:"ee_no,omitempty"`
		MechNo       *string `valid:"optional" json:"mech_no,omitempty"`
		BpUid        *int64  `valid:"optional" json:"bp_uid,omitempty"`
		Phone        *string `valid:"optional" json:"phone,omitempty" `
		Wchat        *string `valid:"optional" json:"wchat,omitempty"`
		Email        *string `valid:"optional" json:"email,omitempty"`
		Vaild        *int    `valid:"optional" json:"valid,omitempty" `
	}

	EmployeeList struct {
		CompanyId    *int64  `valid:"optional" json:"company_id,omitempty"`
		DepartmentId *int64  `valid:"optional" json:"department_id,omitempty"`
		Name         *string `valid:"optional" json:"name,omitempty" `
		Sex          *int    `valid:"optional" json:"sex,omitempty"`
		HiredAt      *int64  `valid:"optional" json:"hired_at,omitempty"`
		BirthAt      *int64  `valid:"optional" json:"birth_at,omitempty"`
		EeNo         *string `valid:"optional" json:"ee_no,omitempty"`
		BpUid        *int64  `valid:"optional" json:"bp_uid,omitempty"`
		MechNo       *string `valid:"optional" json:"mech_no,omitempty"`
		Phone        *string `valid:"optional" json:"phone,omitempty" `
		Wchat        *string `valid:"optional" json:"wchat,omitempty"`
		Email        *string `valid:"optional" json:"email,omitempty"`
		Vaild        *int    `valid:"optional" json:"valid,omitempty" `

		Page int64 `valid:"required" json:"page,omitempty"`
		Size int64 `valid:"optional" json:"size,omitempty"`
	}
)

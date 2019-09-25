package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/models"
	"github.com/BastionPay/bas-bkadmin-api/models/redis"
	"github.com/BastionPay/bas-bkadmin-api/services/access"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
)

const (
	CHANGEPASSWORDEXPIRE = 30
)

var (
	Tools = common.New()
)

type (
	Account struct {
		Id           int64  `valid:"optional" json:"id"`
		Name         string `valid:"required" json:"name"`
		Password     string `valid:"required,length(6|50)" json:"password"`
		Email        string `valid:"email,required" json:"email"`
		Mobile       string `valid:"length(7|13),required" json:"mobile"`
		Department   string `valid:"optional" json:"department"`
		GoogleSecret string `valid:"optional" json:"id"`
		ActualName   string `valid:"required" json:"actual_name"`
		RoleId       int64  `valid:"required" json:"role_id"`
	}

	AccountList struct {
		Id        string `valid:"optional" json:"name" params:"name"`
		Name      string `valid:"optional" json:"name" params:"name"`
		Email     string `valid:"email,optional" json:"email" params:"email"`
		Mobile    string `valid:"length(7|13),optional" json:"mobile" params:"mobile"`
		StartDate string `valid:"optional" json:"start_date" params:"start_date" format:"2018-01-01"`
		EndDate   string `valid:"optional" json:"end_date" params:"end_date" format:"2018-01-01"`
		RoleId    int64  `valid:"optional" json:"role_id" params:"role_id"`
		Status    string `valid:"optional" json:"status" params:"status"`
		Page      int64  `valid:"optional" json:"page" params:"page"`
		Size      int64  `valid:"optional" json:"size" params:"size"`
	}

	AccountUpdate struct {
		Id           int64  `valid:"required" json:"id"`
		Name         string `valid:"optional" json:"name"`
		Password     string `valid:"optional" json:"password"`
		Department   string `valid:"optional" json:"department"`
		GoogleSecret string `valid:"optional" json:"google_secret"`
		Email        string `valid:"optional, email" json:"email"`
		Mobile       string `valid:"optional, length(7|13)" json:"mobile"`
		ActualName   string `valid:"optional" json:"actual_name"`
		RoleId       int64  `valid:"optional" json:"role_id"`
		IsAdmin      int    `valid:"optional,in(1|2)" json:"is_admin"`
		Status       int    `valid:"optional" json:"status"`
		Valid        int    `valid:"optional" json:"valid"`
	}

	AccountDel struct {
		Id    int64 `valid:"required" json:"id"`
		Valid int   `valid:"optional" json:"valid"`
	}

	ChangeBeforePassword struct {
		Id              int64  `valid:"-" json:"id"`
		OldPassword     string `valid:"required" json:"old_password"`
		Password        string `valid:"required" json:"password"`
		ConfirmPassword string `valid:"required" json:"confirm_password"`
		Uuid            string `json:"uuid"`
	}

	ChangeAfterPassword struct {
		Id   int64  `valid:"-" json:"id"`
		Uuid string `json:"uuid" valid:"required"`
		Code string `json:"code" valid:"required"`
	}

	ChangeUserPassword struct {
		Id              int64  `valid:"required" json:"id"`
		Password        string `valid:"required" json:"password"`
		ConfirmPassword string `valid:"required" json:"confirm_password"`
	}

	BatchUserByIds struct {
		Id   string `valid:"required" json:"id" params:"id"`
		Page int64  `valid:"optional" json:"page" params:"page"`
		Size int64  `valid:"optional" json:"size" params:"size"`
	}

	UserInfo struct {
		Id int64 `valid:"required" json:"id" params:"id"`
	}

	SearchList struct {
		models.Account
		LoginTime string `json:"login_time"`
		RoleName  string `json:"role_name"`
	}

	AccountResult struct {
		Id           int64  `json:"id"`
		Name         string `json:"name"`
		ActualName   string `json:"actual_name"`
		Department   string `json:"department"`
		GoogleSecret string `json:"google_secret"`
		Email        string `json:"email"`
		Mobile       string `json:"mobile"`
		RoleId       int64  `json:"role_id"`
		IsAdmin      int    `json:"is_admin"`
		Status       int    `json:"status"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}
)

func (u *UserInfo) GetUserInfo() (*AccountResult, error) {
	var info AccountResult

	err := models.DB.Table(new(models.Account).TableName()).Where("id = ? AND valid = ?", u.Id, "0").First(&info).Error
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (this *Account) Register(params *Account) (*models.Account, error) {
	params.Name = strings.ToLower(params.Name)
	password := Tools.MD5(params.Password)

	account := &models.Account{
		Name:       params.Name,
		Password:   password,
		ActualName: params.ActualName,
		RoleId:     params.RoleId,
		Email:      params.Email,
		Mobile:     params.Mobile,
		Status:     1,
		Department: params.Department,
	}

	tx := models.DB.Begin()

	err := tx.Table(account.TableName()).Create(account).Error
	if err == nil {
		if account.RoleId > 0 {
			userRole := &models.UserRole{
				UserId: account.Id,
				RoleId: this.RoleId,
				Model: models.Model{
					CreatedAt: Tools.GetDateNowString(),
					UpdatedAt: Tools.GetDateNowString(),
				},
			}

			err = tx.Table(new(models.UserRole).TableName()).Create(userRole).Error
			if err == nil {
				tx.Commit()
				return account, nil
			}
		} else {
			tx.Commit()
			return account, nil
		}
	}

	tx.Rollback()
	return nil, err
}

func (this *Account) GetUserById(userId int64) (*models.Account, error) {
	user := &models.Account{
		Id: userId,
	}

	err := models.DB.First(&user).Error

	return user, err
}

func (this *Account) GetUserInfoByToken(token string) (*models.Account, error) {
	var user *models.Account

	body, err := redis.RedisClient.Get(token).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(body), &user)
	return user, err
}

func (this *AccountUpdate) UpdateUserInfo() (*models.Account, error) {
	var user models.Account

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&user).Error
	if err != nil {
		return nil, err
	}

	if this.Name != "" {
		user.Name = this.Name
	}

	if this.Email != "" {
		user.Email = this.Email
	}

	if this.Mobile != "" {
		user.Mobile = this.Mobile
	}

	if this.Department != "" {
		user.Department = this.Department
	}

	if this.GoogleSecret != "" {
		user.GoogleSecret = this.GoogleSecret
	}

	if this.Password != "" {
		user.Password = Tools.MD5(this.Password)
	}

	if this.RoleId > 0 && this.RoleId != user.RoleId {
		user.RoleId = this.RoleId

		var userRole models.UserRole
		models.DB.Where("role_id = ? AND user_id = ? AND valid = ?", user.RoleId, user.Id, "0").First(&userRole)

		if userRole.Id > 0 {
			userRole.RoleId = this.RoleId
			userRole.UpdatedAt = Tools.GetDateNowString()
			models.DB.Save(userRole)
		} else {
			data := &models.UserRole{
				RoleId: this.RoleId,
				UserId: user.Id,
			}

			models.DB.Create(data)
		}
	}

	if this.ActualName != "" {
		user.ActualName = this.ActualName
	}

	err = models.DB.Model(&models.Account{}).Where("id = ? AND valid = ?", this.Id, "0").Update(&user).Error
	return &user, err
}

func (this *AccountUpdate) SetUserAdmin(token string) (*models.Account, error) {
	var account models.Account

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&account).Error
	if err != nil {
		return nil, err
	}

	tx := models.DB.Begin()

	account.IsAdmin = this.IsAdmin
	err = tx.Save(&account).Error
	if err == nil {
		body, err := json.Marshal(account)
		if err == nil {
			_, err = redis.RedisClient.Set(Tools.GenerateUserLoginToken(account.Id), body, 24*time.Second).Result()
			if err == nil {
				tx.Commit()
				return &account, err
			}
		}
	}

	tx.Rollback()
	return &account, err
}

func (this *AccountUpdate) DisabledUser(token string) (*models.Account, error) {
	var account models.Account

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&account).Error
	if err != nil {
		return nil, err
	}

	tx := models.DB.Begin()

	account.Status = this.Status
	err = tx.Save(&account).Error
	if err == nil {
		RemoveToken(Tools.GenerateUserLoginToken(account.Id))
		tx.Commit()
		return &account, err
	}

	tx.Rollback()
	return nil, err
}

func (this *AccountDel) Delete() (*models.Account, error) {
	var account models.Account

	err := models.DB.Where("id = ? AND valid = ?", this.Id, "0").First(&account).Error
	if err != nil {
		return nil, err
	}

	account.Valid = this.Valid
	err = models.DB.Save(&account).Error
	if err != nil {
		return nil, err
	}

	RemoveToken(Tools.GenerateUserLoginToken(this.Id))
	return &account, nil
}

func RemoveToken(token string) bool {
	num, err := redis.RedisClient.Del(token).Result()
	if err == nil && num > 0 {
		return true
	}

	return false
}

func (this *AccountList) Search() (*common.Result, error) {
	var list []*models.Account

	query := models.DB.Where("valid = ?", "0")

	if this.Name != "" {
		query = query.Where("name = ?", this.Name)
	}

	if this.Email != "" {
		query = query.Where("email = ?", this.Email)
	}

	if this.Mobile != "" {
		query = query.Where("mobile = ?", this.Mobile)
	}

	if this.Id != "" {
		ids := strings.Split(this.Id, ",")
		query = query.Where("id in (?)", ids)
	}

	if this.RoleId > 0 {
		query = query.Where("role_id = ?", this.RoleId)
	}

	if this.StartDate != "" && this.EndDate != "" {
		start := this.StartDate + " 00:00:00"
		end := this.EndDate + " 23:59:59"

		query = query.Where("created_at >= ? AND created_at <= ?", start, end)
	}

	if this.Status != "" {
		query = query.Where("status = ?", this.Status)
	}

	return new(common.Result).PageQuery(
		query,
		&models.Account{},
		&list,
		this.Page,
		this.Size,
		&AccountList{},
		"GetLogTime")
}

func (this *BatchUserByIds) BatchUserByIds() (*common.Result, error) {
	var list []*models.Account
	var total int64

	ids := strings.Split(this.Id, ",")

	offset := (this.Page - 1) * this.Size
	err := models.DB.Model(&models.Account{}).
		Where("id IN (?)", ids).Count(&total).
		Offset(offset).Limit(this.Size).Find(&list).Error

	if err != nil {
		return nil, err
	}

	return new(common.Result).PageResult(list, total, this.Page, this.Size)
}

func (this *AccountList) GetLogTime(data interface{}) []*SearchList {
	var searchList []*SearchList

	body, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(body, &searchList)
	if err != nil {
		return nil
	}

	var list []*SearchList

	roleMapList := make(map[int64]string)
	roleIds := this.ColumnRoleId(searchList)
	if len(roleIds) > 0 {
		roleIds = Tools.Unique(roleIds)
		var result []*models.Role
		models.DB.Where("valid = ? AND id IN (?)", "0", roleIds).Find(&result)
		if err == nil {
			for _, val := range result {
				roleMapList[val.Id] = val.Name
			}
		}
	}

	for _, v := range searchList {
		log := access.UserLoginLog{
			UserId: v.Id,
		}

		logInfo, err := log.GetLoginLog()
		if err == nil {
			v.LoginTime = logInfo.CreatedAt
		}

		if v.RoleId > 0 {
			roleName, ok := roleMapList[v.RoleId]
			if ok {
				v.RoleName = roleName
			}
		}

		list = append(list, v)
	}

	return list
}

func (this *AccountList) ColumnRoleId(list []*SearchList) []string {
	var roleIdList []string

	for _, v := range list {
		roleIdList = append(roleIdList, strconv.Itoa(int(v.RoleId)))
	}

	return roleIdList
}

func (this *AccountList) GetRoleInfoById(roleId int64) (*models.Role, error) {
	var roles models.Role

	err := models.DB.Where("valid = ? AND id = ?", "0", roleId).Find(&roles).Error

	return &roles, err
}

func (this *ChangeUserPassword) ChangeUserPasswords() error {

	if !strings.EqualFold(this.Password, this.ConfirmPassword) {
		return errors.New("two passwords are inconsistent")
	}

	account := &models.Account{}

	err := models.DB.Where("valid = ? AND status = ? AND id = ?", "0", "1", this.Id).First(&account).Error
	if err != nil {
		return err
	}

	if account.Id <= 0 {
		return errors.New("user does not exist")
	}

	account.Password = Tools.MD5(this.Password)

	err = models.DB.Save(&account).Error
	if err != nil {
		return err
	}

	RemoveToken(Tools.GenerateUserLoginToken(account.Id))

	return nil
}

func (this *ChangeBeforePassword) ChangeBeforePassword() (string, error) {
	if !strings.EqualFold(this.Password, this.ConfirmPassword) {
		return "", errors.New("two passwords are inconsistent")
	}

	account := &models.Account{}

	err := models.DB.Where("valid = ? AND status = ? AND id = ?", "0", "1", this.Id).First(&account).Error
	if err != nil {
		return "", err
	}

	if account.Id <= 0 {
		return "", errors.New("user does not exist")
	}

	if strings.EqualFold(account.Password, Tools.MD5(this.Password)) {
		return "", errors.New("the new password and the old password cannot be consistent")
	}

	if !strings.EqualFold(account.Password, Tools.MD5(this.OldPassword)) {
		return "", errors.New("the original password is incorrect")
	}

	u := fmt.Sprintf("%s", uuid.Must(uuid.NewV4()))

	data := map[string]string{
		"password": this.Password,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	_, err = redis.RedisClient.Set(u, body, CHANGEPASSWORDEXPIRE*time.Minute).Result()
	if err != nil {
		return "", err
	}

	return u, nil
}

func (c *ChangeAfterPassword) ChangeAfterPassword() bool {
	body, err := redis.RedisClient.Get(c.Uuid).Result()
	if err != nil {
		return false
	}

	data := struct {
		Password string `json:"password"`
	}{}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return false
	}

	rows := models.DB.Table(new(models.Account).TableName()).
		Where("id = ? AND valid = ?", c.Id, "0").
		Update(map[string]interface{}{
			"password": Tools.MD5(data.Password),
		}).RowsAffected

	if rows <= 0 {
		return false
	}

	RemoveToken(Tools.GenerateUserLoginToken(c.Id))
	RemoveToken(c.Uuid)

	return true
}

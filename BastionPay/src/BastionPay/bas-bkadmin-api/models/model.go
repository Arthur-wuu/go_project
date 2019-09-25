package models

import (
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	"github.com/BastionPay/bas-bkadmin-api/tools"
	l4g "github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strings"
	"time"
)

type (
	Model struct {
		Valid     int    `json:"valid" gorm:"column:valid"`
		CreatedAt string `json:"created_at" gorm:"column:created_at"`
		UpdatedAt string `json:"updated_at" gorm:"column:updated_at"`
	}

	MysqlConfig struct {
		Config *tools.Mysql
	}
)

var (
	DB    *gorm.DB
	err   error
	Tools *common.Tools
)

func init() {
	Tools = common.New()
}

func New(conf *tools.Mysql) *MysqlConfig {
	return &MysqlConfig{
		Config: conf,
	}
}

func (this *MysqlConfig) Connection() *gorm.DB {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		this.Config.User,
		this.Config.Password,
		this.Config.Host,
		this.Config.Port,
		this.Config.DbName,
		this.Config.Charset,
		this.Config.ParseTime)

	DB, err = gorm.Open(this.Config.Dialect, conn)
	if err != nil {
		l4g.Crash(err, "mysql connection errors")
	}

	DB.DB().SetMaxIdleConns(this.Config.MaxIdle)
	DB.DB().SetMaxOpenConns(this.Config.MaxOpen)
	DB.Callback().Create().Replace("gorm:update_time_stamp", this.updateTimeStampForCreateCallback)
	DB.Callback().Update().Replace("gorm:update_time_stamp", this.updateTimeStampForUpdateCallback)
	//DB.Callback().Delete().Replace("gorm:delete", db.deleteCallback)

	DB.SingularTable(true)
	DB.LogMode(this.Config.Debug)
	DB.Set("gorm:table_options", "ENGINE=InnoDB")

	return DB
}

func (this *Model) BeforeSave() {
	this.CreatedAt = Tools.GetDateNowString()
	this.UpdatedAt = Tools.GetDateNowString()
}

func (this *Model) BeforeCreate() {
	this.CreatedAt = Tools.GetDateNowString()
	this.UpdatedAt = Tools.GetDateNowString()
}

func (this *Model) BeforeUpdate() {
	this.UpdatedAt = Tools.GetDateNowString()
}

func (this *Model) BatchInsert(table string, fields, values []string) *gorm.DB {
	sql := fmt.Sprintf("INSERT INTO %s ( %s ) VALUES ", table, strings.Join(fields, ","))

	var i = 0
	var strSql = ""
	for _, v := range values {
		if i > 0 {
			strSql += " , "
		}

		strSql += fmt.Sprintf("( %s )", v)
		i++
	}

	sql = sql + strSql
	return DB.Exec(sql)
}

func (d *MysqlConfig) updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := NowTimestamp()

		if createdAtField, ok := scope.FieldByName("CreatedAt"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

func (d *MysqlConfig) updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if !scope.HasError() {

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(NowTimestamp())
			}
		}
	}
}

func (d *MysqlConfig) deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedAtField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(NowTimestamp()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func NowTimestamp() int64 {
	return time.Now().UnixNano() / (1000 * 1000)
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

package common

import (
	"fmt"

	l4g "github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Model struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"default:null"`
	DeletedAt int64 `json:"deleted_at" gorm:"default:null"`
}

type Db struct {
	conn *gorm.DB
}

type DbOptions struct {
	Host        string
	Port        string
	User        string
	Pass        string
	DbName      string
	MaxIdleConn int
	MaxOpenConn int
	Debug       bool
}

func NewDb(options *DbOptions) (*Db, error) {
	var (
		db   *Db
		conn *gorm.DB
		err  error
	)
	db = &Db{}

	conn, err = gorm.Open("mysql",
		options.User+":"+options.Pass+"@tcp("+options.Host+":"+options.Port+")/"+options.DbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		l4g.Error(err.Error())
		return nil, err
	}
	db.conn = conn

	conn.SingularTable(true)
	conn.LogMode(options.Debug)

	conn.DB().SetMaxIdleConns(options.MaxIdleConn)
	conn.DB().SetMaxOpenConns(options.MaxOpenConn)

	conn.Callback().Create().Replace("gorm:update_time_stamp", db.updateTimeStampForCreateCallback)
	conn.Callback().Update().Replace("gorm:update_time_stamp", db.updateTimeStampForUpdateCallback)
	conn.Callback().Delete().Replace("gorm:delete", db.deleteCallback)

	return db, nil
}

func (d *Db) GetConn() *gorm.DB {
	return d.conn
}

func (d *Db) updateTimeStampForCreateCallback(scope *gorm.Scope) {
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

func (d *Db) updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if !scope.HasError() {

		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(NowTimestamp())
			}
		}
	}
}

func (d *Db) deleteCallback(scope *gorm.Scope) {
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

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

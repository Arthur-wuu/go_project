package db

import (
"fmt"
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/mysql"
"time"
)

var GDbMgr DbMgr

type DbMgr struct {
	mConn *gorm.DB
	mOption   *DbOptions
	mFlag bool
}


func (this *DbMgr) Init(options *DbOptions) (err error) {
	if this.mFlag {
		return nil
	}
	this.mOption = options
	fmt.Println("db=", *options)
	this.mConn, err = gorm.Open("mysql",
		options.User+":"+options.Pass+"@tcp("+options.Host+":"+options.Port+")/"+options.DbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return err
	}

	this.mConn.SingularTable(true)
	this.mConn.LogMode(false)

	this.mConn.DB().SetMaxIdleConns(options.MaxIdleConn)
	this.mConn.DB().SetMaxOpenConns(options.MaxOpenConn)

	//this.mConn.Callback().Create().Replace("gorm:update_time_stamp", this.updateTimeStampForCreateCallback)
	//this.mConn.Callback().Update().Replace("gorm:update_time_stamp", this.updateTimeStampForUpdateCallback)
	//this.mConn.Callback().Delete().Replace("gorm:delete", this.deleteCallback)

	//	this.mConn.AutoMigrate(&CodeTable{})
	this.mFlag = true
	return nil
}

func (this *DbMgr) Get() *gorm.DB {
	if !this.mFlag {
		err := this.Init(this.mOption)
		if err != nil {
			return nil
		}
	}
	return this.mConn
}

func (this *DbMgr) Close() {
	if !this.mFlag {
		return
	}
	this.mConn.Close()
}

func (d *DbMgr) updateTimeStampForCreateCallback(scope *gorm.Scope) {
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

func (d *DbMgr) updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		if updatedAtField, ok := scope.FieldByName("UpdatedAt"); ok {
			//	if updatedAtField.IsBlank {
			updatedAtField.Set(NowTimestamp())
			//	}
		}
	}
}

func (d *DbMgr) deleteCallback(scope *gorm.Scope) {
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

func NowTimestamp() int64 {
	return time.Now().Unix()
}


package db

import (
	"BastionPay/bas-base/log/zap"
	"BastionPay/bas-filetransfer-srv/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"strings"
	"time"
)

var GDbMgrs map[string]*DbMgr

func Init() error {
	GDbMgrs = make(map[string]*DbMgr)
	for i := 0; i < len(config.GConfig.Dbs); i++ {
		c := config.GConfig.Dbs[i]
		c.Dbname = strings.Replace(c.Dbname, " ", "", len(c.Dbname))
		dbnameArr := strings.Split(c.Dbname, ",")
		if len(c.Dbname) == 0 || len(dbnameArr) == 0 {
			continue
		}
		dbMgr := new(DbMgr)

		err := dbMgr.Init(&DbOptions{
			Host:        config.GConfig.Dbs[i].Host,
			Port:        config.GConfig.Dbs[i].Port,
			User:        config.GConfig.Dbs[i].User,
			Pass:        config.GConfig.Dbs[i].Pwd,
			DbName:      config.GConfig.Dbs[i].Dbname,
			MaxIdleConn: config.GConfig.Dbs[i].Max_idle_conn,
			MaxOpenConn: config.GConfig.Dbs[i].Max_open_conn,
		})
		if err != nil {
			log.ZapLog().Error("mysql init err", zap.Error(err), zap.String("dbname", config.GConfig.Dbs[i].Dbname))
		}
		for j := 0; j < len(dbnameArr); j++ {
			GDbMgrs[dbnameArr[j]] = dbMgr
		}
	}
	return nil
}

type DbMgr struct {
	mConn   *gorm.DB
	mOption *DbOptions
	mFlag   bool
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

	this.mConn.Callback().Create().Replace("gorm:update_time_stamp", this.updateTimeStampForCreateCallback)
	this.mConn.Callback().Update().Replace("gorm:update_time_stamp", this.updateTimeStampForUpdateCallback)
	this.mConn.Callback().Delete().Replace("gorm:delete", this.deleteCallback)

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

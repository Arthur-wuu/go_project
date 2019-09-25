package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"strconv"
	"time"
)

type DbMgr struct {
	mConn *gorm.DB
}

func (this *DbMgr) Init(options *DbOptions) (err error) {
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
	return nil
}

func (this *DbMgr) Close() {
	this.mConn.Close()
}

func (this *DbMgr) AddCode(info *CodeTable) error {
	if info.UpdatedAt == nil {
		info.UpdatedAt = new(int64)
		*info.UpdatedAt = NowTimestamp()
	}
	//	fmt.Printf("=====%v\n", *info)
	newDB := this.mConn.Model(&CodeTable{}).Where("symbol = ?", info.Symbol).Update(info)
	if newDB.Error != nil {
		return newDB.Error
	}
	if newDB.RowsAffected != 0 {
		return nil
	}
	return this.mConn.Create(info).Error
}

//修改码表的 vaild 值，将无效的码表的vaild值由 1 ==> 0
func (this *DbMgr) ModelCode(info *CodeTable) error {
	if info.UpdatedAt == nil {
		info.UpdatedAt = new(int64)
		*info.UpdatedAt = NowTimestamp()
	}
	//	fmt.Printf("=====%v\n", *info)
	newDB := this.mConn.Model(&CodeTable{}).Where("code = ?", info.Code).Update("vaild", 0)
	if newDB.Error != nil {
		return newDB.Error
	}
	if newDB.RowsAffected != 0 {
		return nil
	}
	return this.mConn.Create(info).Error
}

func (this *DbMgr) GetAllCode() ([]CodeTable, error) {
	arr := make([]CodeTable, 0)
	err := this.mConn.Find(&arr).Error
	return arr, err
}

//通过code找到码表的信息
func (this *DbMgr) GetCodeByCode(code string) *CodeTable {
	arr := new(CodeTable)
	codeInt, err := strconv.Atoi(code)
	if err != nil {
		return nil
	}
	this.mConn.First(arr, "code = ?", codeInt)
	fmt.Println("GetCodeByCode********", arr)
	return arr
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

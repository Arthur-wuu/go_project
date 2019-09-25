package models

import (
	"BastionPay/bas-filetransfer-srv/common"
	"BastionPay/bas-filetransfer-srv/db"
	"encoding/csv"
	"reflect"

	. "BastionPay/bas-base/log/zap"
	"BastionPay/bas-filetransfer-srv/config"
	"BastionPay/bas-tools/sdk.aws.sts"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func (this *TaskExportInfo) MakeFile() error {
	//判断任务是否过期或者被取消
	exist, err := db.GRedis.Exist(EXPORT_Status_KeyPrefix + this.OrderId)
	if err != nil {
		ZapLog().Error("Redis.Exist err", zap.Error(err), zap.String("orderid", this.OrderId))
		return err
	}
	if !exist {
		ZapLog().Info("task cancel or expire", zap.String("orderid", this.OrderId))
		return nil
	}

	if this.WaitExpireAt <= time.Now().Unix() {
		this.updateTaskStatus("", TASKSTATUS_EXPIRE, "wait expire")
		ZapLog().Info("task cancel or expire", zap.String("orderid", this.OrderId))
		return nil
	}

	//判断文件是否已存在于s3上
	if config.GConfig.Task.FileUseExist {
		fileStatus, err := db.GFileS3Sdk.Status(config.GConfig.Aws.FileBucket, this.FileName+EXPORT_File_Exd)
		if err != nil {
			ZapLog().Error("GFileS3Sdk.Status err", zap.Error(err), zap.String("orderid", this.OrderId), zap.String("filename", this.FileName+EXPORT_File_Exd))
			this.updateTaskStatus("", TASKSTATUS_FAIL, "s3, "+err.Error())
			return err
		}
		if fileStatus.Exist() && !fileStatus.Timeout() {
			newAddr, err := this.genSignUrl(fileStatus.Addr(), this.FileName+EXPORT_File_Exd)
			if err != nil {
				ZapLog().Error("genSignUrl err", zap.Error(err), zap.String("orderid", this.OrderId))
				this.updateTaskStatus("", TASKSTATUS_FAIL, "signurl, "+err.Error())
				return err
			}
			return this.updateTaskStatus(newAddr, TASKSTATUS_Finish, "use exist one")
		}
	}

	//更新状态
	if err := this.updateTaskStatus("", TASKSTATUS_Doing, ""); err != nil {
		ZapLog().Error("updateTaskStatus err", zap.Error(err), zap.String("orderid", this.OrderId))
		return err
	}

	//生成文件
	filePath := config.GConfig.Server.TmpPath + "/" + this.FileName + "." + this.FileFormat

	if this.FileFormat == "xlsx" {
		if err = this.dbToFile(filePath); err != nil {
			ZapLog().Error("dbToFile err", zap.Error(err), zap.Any("Task", *this))
			this.updateTaskStatus("", TASKSTATUS_FAIL, "dbTofile, "+err.Error())
			return err
		}
	} else if this.FileFormat == "csv" {
		if err = this.dbToFileCsv(filePath); err != nil {
			ZapLog().Error("dbToFile err", zap.Error(err), zap.Any("Task", *this))
			this.updateTaskStatus("", TASKSTATUS_FAIL, "dbTofile, "+err.Error())
			return err
		}
	}

	ZapLog().Info("dbToFile ok", zap.String("orderid", this.OrderId), zap.String("filepath", filePath))

	//压缩
	zipPath := config.GConfig.Server.TmpPath + "/" + this.FileName + EXPORT_File_Exd
	if err := common.NewZip().Compress(filePath, zipPath); err != nil {
		ZapLog().Error("Compress err", zap.Error(err), zap.Any("Task", *this))
		this.updateTaskStatus("", TASKSTATUS_FAIL, "compress, "+err.Error())
		return err
	}
	ZapLog().Info("Compress ok", zap.String("orderid", this.OrderId), zap.String("path", filePath+" => "+zipPath))

	//上传文件
	addr, err := this.UploadFile(zipPath)
	if err != nil {
		ZapLog().Error("GFileS3Sdk.UpLoad err", zap.Error(err), zap.String("zipPath", zipPath))
		this.updateTaskStatus("", TASKSTATUS_FAIL, "UploadFile, "+err.Error())
		return err
	} //上传出错了还没处理

	ZapLog().Info("upload zipfile ok", zap.String("orderid", this.OrderId))
	//清理下
	os.Remove(filePath)
	os.Remove(zipPath)

	//签名url
	newAddr, err := this.genSignUrl(addr, this.FileName+EXPORT_File_Exd)
	if err != nil {
		ZapLog().Error("genSignUrl err", zap.Error(err), zap.String("orderid", this.OrderId))
		this.updateTaskStatus("", TASKSTATUS_FAIL, "signUrl, "+err.Error())
		return err
	}

	if err = this.updateTaskStatus(newAddr, TASKSTATUS_Finish, "upload new one"); err != nil {
		ZapLog().Error("updateTaskStatus err", zap.Error(err), zap.String("orderid", this.OrderId))
		return err
	}
	return nil
}

func (this *TaskExportInfo) dbToFile(filePath string) error {
	dbMgr, ok := db.GDbMgrs[this.Dbname]
	if !ok {
		return errors.New("db not in config")
	}
	dbCon := dbMgr.Get()
	if dbCon == nil {
		return errors.New("db cannot connected")
	}
	allCount := uint64(0)

	ZapLog().Info("dbSql=" + string(this.Sql))
	rowArr := make([]interface{}, 0)
	var colNames []string

	for i := 0; true; i++ {
		count := uint64(0)
		newSql := this.toNewSql(string(this.Sql), this.PageSize*i, this.PageSize)

		stmt, err := dbCon.DB().Prepare(newSql)
		if err != nil {
			ZapLog().Error("Db Query err", zap.Error(err))
			break
		}
		rows, err := stmt.Query()
		if err != nil {
			ZapLog().Error("Db Query err", zap.Error(err))
			break
		}
		for rows.Next() {
			if err = rows.Err(); err != nil {
				ZapLog().Error("Db Query rows.Err err", zap.Error(err))
				break
			}
			var calArr []interface{}
			calArr, colNames, err = this.scan(rows)
			if err != nil {
				break
			}
			ZapLog().Debug("db Next  start scan", zap.Int("num", len(calArr)), zap.Any("data", calArr))

			rowArr = append(rowArr, calArr)
			count++
		}
		rows.Close()
		if count == 0 {
			ZapLog().Debug("Db Query count=0")
			break
		}
		allCount += count
		if this.Max_lines != 0 && allCount >= this.Max_lines {
			break
		}
		if err != nil {
			ZapLog().Error("Db Query rows.Err err", zap.Error(err))
			break
		}
		if this.PageHalt != 0 {
			time.Sleep(time.Second * time.Duration(this.PageHalt))
		}
	}
	ZapLog().Info("read from db ", zap.Int("num", len(rowArr)))

	return this.createFile(filePath, rowArr, colNames)
}

//生成csv文件
func (this *TaskExportInfo) dbToFileCsv(filePath string) error {
	dbMgr, ok := db.GDbMgrs[this.Dbname]
	if !ok {
		return errors.New("db not in config")
	}
	dbCon := dbMgr.Get()
	if dbCon == nil {
		return errors.New("db cannot connected")
	}
	allCount := uint64(0)

	ZapLog().Info("dbSql=" + string(this.Sql))
	totalValues := make([][]string, 0)
	var columns []string

	for i := 0; true; i++ {
		count := uint64(0)
		newSql := this.toNewSql(string(this.Sql), this.PageSize*i, this.PageSize)

		stmt, err := dbCon.DB().Prepare(newSql)
		if err != nil {
			ZapLog().Error("Db Query err", zap.Error(err))
			break
		}
		rows, err := stmt.Query()
		if err != nil {
			ZapLog().Error("Db Query err", zap.Error(err))
			break
		}

		columns, err = rows.Columns()
		if err != nil {
			ZapLog().Error("Db Query columns err", zap.Error(err))
			break
		}

		values := make([]sql.RawBytes, len(columns))

		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			var s []string
			err = rows.Scan(scanArgs...)
			if err != nil {
				ZapLog().Error("Db scan args err", zap.Error(err))
				break
			}

			for _, v := range values {
				s = append(s, string(v))
			}
			totalValues = append(totalValues, s)
		}

		rows.Close()
		if count == 0 {
			ZapLog().Debug("Db Query count=0")
			break
		}
		allCount += count
		if this.Max_lines != 0 && allCount >= this.Max_lines {
			break
		}
		if err != nil {
			ZapLog().Error("Db Query rows.Err err", zap.Error(err))
			break
		}
		if this.PageHalt != 0 {
			time.Sleep(time.Second * time.Duration(this.PageHalt))
		}
	}
	ZapLog().Info("read from db ", zap.Int("num", len(totalValues)))

	return this.writeTocsv(filePath, columns, totalValues)
}

//返回行数据，行名称
func (this *TaskExportInfo) scan(rows *sql.Rows) ([]interface{}, []string, error) {
	colNames, err := rows.Columns()
	if err != nil {
		ZapLog().Error("Db rows.Columns err", zap.Error(err))
		return nil, nil, err
	}
	calArr := make([]interface{}, len(colNames))
	for m := 0; m < len(colNames); m++ {
		calArr[m] = new(interface{})
	}
	if err := rows.Scan(calArr...); err != nil {
		ZapLog().Error("rows.Scan err", zap.Error(err))
		return nil, nil, err
	}

	for m := 0; m < len(colNames); m++ {
		dd, ok := calArr[m].(*interface{})
		if !ok {
			ZapLog().Error("type err")
			continue
		}
		tp := reflect.TypeOf(*dd)

		if tp.String() == "[]uint8" {
			v := reflect.ValueOf(*dd)
			*dd = string(v.Bytes())
		}
		calArr[m] = *dd
	}
	return calArr, colNames, nil
}

func (this *TaskExportInfo) toNewSql(oldSql string, offset, limit int) string {
	return oldSql + fmt.Sprintf("  limit %d offset %d ;", limit, offset)
}

func (this *TaskExportInfo) UploadFile(zipPath string) (string, error) {
	zipReader, err := os.Open(zipPath)
	if err != nil {
		ZapLog().Error("Open err", zap.Error(err), zap.String("zipPath", zipPath))
		//this.updateTaskStatus("", db.TASKSTATUS_FAIL)
		return "", err
	}
	addr, err := db.GFileS3Sdk.UpLoadEx(config.GConfig.Aws.FileBucket, this.FileName+EXPORT_File_Exd, zipReader, 60*5, config.GConfig.Task.FileKeepTime)
	if err != nil {
		ZapLog().Error("GFileS3Sdk.UpLoad err", zap.Error(err), zap.String("zipPath", zipPath))
		//this.updateTaskStatus("", db.TASKSTATUS_FAIL)
		return "", err
	}
	return addr, nil
}

func (this *TaskExportInfo) createFile(filePath string, rowArr []interface{}, colNames []string) error {
	switch this.FileFormat {
	case "xlsx":
		return this.createXlsFile(filePath, rowArr, colNames)
		break
	//case "csv":
	//	return this.createCsvFile(filePath , rowArr , colNames)
	//	break
	default:
		return errors.New("unknown fileformat")
	}
	return errors.New("unknown fileformat")
}

func (this *TaskExportInfo) createXlsFile(filePath string, rowArr []interface{}, colNames []string) error {
	xlsObj, err := common.NewXlsx(rowArr, colNames, nil)
	if err != nil {
		ZapLog().Error("NewXlsx err", zap.Error(err))
		return err
	}
	if err = xlsObj.Generate(); err != nil {
		ZapLog().Error("xlsObj.Generate err", zap.Error(err))
		return err
	}
	if err = xlsObj.File(filePath); err != nil {
		ZapLog().Error("xlsObj.File err", zap.Error(err))
		return err
	}
	return nil
}

func (this *TaskExportInfo) writeTocsv(file string, columns []string, totalValues [][]string) error {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(f)
	for i, row := range totalValues {
		if i == 0 {
			w.Write(columns)
			w.Write(row)
		} else {
			w.Write(row)
		}
	}
	w.Flush()
	return nil
}

func (this *TaskExportInfo) genSignUrl(url, filename string) (string, error) {
	stat := sdk_aws_sts.Statement{
		Action:   "s3:GetObject",
		Effect:   "Allow",
		Resource: "arn:aws:s3:::" + config.GConfig.Aws.FileBucket + "/" + filename,
	}
	cred, err := db.GStsSdk.GetFederationToken(3600, "pingzilao", []sdk_aws_sts.Statement{stat})
	if err != nil {
		return "", err
	}
	newCred := credentials.NewStaticCredentials(*cred.AccessKeyId, *cred.SecretAccessKey, *cred.SessionToken)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	sign := v4.NewSigner(newCred)
	_, err = sign.Presign(httpReq, nil, "s3", config.GConfig.Aws.FileRegion, time.Duration(time.Second*3600), time.Now())
	if err != nil {
		return "", err
	}
	fmt.Println(httpReq.URL.String())
	return httpReq.URL.String(), nil
}

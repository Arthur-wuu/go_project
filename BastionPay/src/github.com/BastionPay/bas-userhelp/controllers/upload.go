package controllers

import (
	s3 "BastionPay/bas-tools/sdk.aws.s3"
	"github.com/kataras/iris"

	"strings"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-userhelp/config"
	"go.uber.org/zap"
	. "BastionPay/bas-base/log/zap"
)

type UploadFile struct {
	Controllers
}

type LogoFileAddr struct {
	File string `json:"file"`
	Addr string `json:"addr"`
}

func (this *UploadFile) HandlePicFiles(ctx iris.Context) {
	this.handleFilesToAwsS3(ctx, config.GConfig.Aws.PicRegion, config.GConfig.Aws.PicBucket, "", "", config.GConfig.Aws.PicTimeout)
}

func (this *UploadFile) handleFilesToAwsS3(ctx iris.Context, region, bucket, bucketPath, filePrefix string, timeout int) {
	//设置内存大小
	err := ctx.Request().ParseMultipartForm(32 << 23)
	if err != nil {
		ZapLog().Error( "ParseMultipartForm err", zap.Error(err))
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "ParseMultipartForm_ERRORS")
		return
	}
	files := ctx.Request().MultipartForm.File["file"]
	if len(files) == 0 {
		ZapLog().Error( "NoFile err")
		this.ExceptionSerive(ctx, apibackend.BASERR_INVALID_PARAMETER.Code(), "NoFind file")
		return
	}
	//暂时s3Sdk作为一个临时的，因为压根很少用到，何必占内存呢
	s3Sdk := s3.NewS3Sdk(region, config.GConfig.Aws.AccessKeyId, config.GConfig.Aws.AccessKey, config.GConfig.Aws.AccessToken)
	if s3Sdk == nil {
		ZapLog().Error( "NewS3Sdk  err[return nil]")
		this.ExceptionSerive(ctx, apibackend.BASERR_UNKNOWN_BUG.Code(), "NewS3Sdk_ERRORS")
		return
	}
	defer s3Sdk.Close()
	if timeout < 10 {
		timeout = 10
	}
	results := make([]LogoFileAddr, 0)
	for i := 0; i < len(files); i++ {
		filename := files[i].Filename
		//l4g.Debug("fileName[%s] ", filename)
		file, err := files[i].Open()
		if err != nil {
			ZapLog().Error( "FileOpen  err", zap.Error(err))
			this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "FileOpen_ERRORS")
			return
		}
		defer file.Close()
		//content, err := ioutil.ReadAll(file)
		//if err != nil {
		//	ZapLog().Error( "FileReadAll  err", zap.Error(err))
		//	this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "File_ReadAll_ERRORS")
		//	return
		//}
		_, err = file.Seek(0, 0)
		if err != nil {
			ZapLog().Error( "FileSeek  err", zap.Error(err))
			this.ExceptionSerive(ctx, apibackend.BASERR_SYSTEM_INTERNAL_ERROR.Code(), "File_Seek_ERRORS")
			return
		}

		//l4g.Info("S3 start UpLoad[%s]md5[%s] ", filename, filenameMd5)
		addr, err := s3Sdk.UpLoad(bucket, filename, file, timeout)
		if err != nil {
			ZapLog().Error( "S3 UpLoad  err", zap.Error(err))
			this.ExceptionSerive(ctx, apibackend.BASERR_INTERNAL_SERVICE_ACCESS_ERROR.Code(), "S3_UpLoad_ERRORS")
			return
		}
		addr = strings.TrimRight(addr, "/")
		results = append(results, LogoFileAddr{filename, addr})
	}
	//content, err := json.Marshal(results)
	//if err != nil {
	//	ZapLog().Error( "Marshal  err", zap.Error(err))
	//	this.ExceptionSerive(ctx, apibackend.BASERR_DATA_PACK_ERROR.Code(), "Marshal_LogFileAddr_ERROR")
	//	return
	//}
	this.Response(ctx, results)
	//l4g.Debug("deal handleFilesToAwsS3 ok, result[%s]",  string(content))
}


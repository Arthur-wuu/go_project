package sdk_aws_s3

import (
	"context"
	"errors"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	pkgSess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"strings"
	"time"
)

// 参数说明：区域，密钥Id，密钥，token（一般为空）
func NewS3Sdk(region, accessKeyID, accessKey, accessToken string) *S3Sdk {
	config := &aws.Config{
		Region: aws.String(region),
	}
	if (len(accessKeyID) != 0) && (len(accessKey) != 0) {
		config.Credentials = credentials.NewStaticCredentials(accessKeyID, accessKey, accessToken)
	}
	client := &S3Sdk{
		mSess: pkgSess.Must(pkgSess.NewSession(config)),
	}
	client.mS3Client = s3.New(client.mSess)
	return client
}

type S3Sdk struct {
	mSess     *pkgSess.Session
	mS3Client *s3.S3
}

func (this *S3Sdk) Close() {

}

//阻塞模式，用并发
//参数说明：key是文件名；value是文件内容；timeout文件上传超时值，超时将中断传输。
//返回值：location，url地址
func (this *S3Sdk) UpLoad(bucket, key string, value io.Reader, timeout int, header map[string][]string) (string, error) {
	if timeout <= 0 {
		timeout = 120 //default 60s
	}
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancelFn()

	uploader := s3manager.NewUploaderWithClient(this.mS3Client, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024 // 64MB per part
		u.LeavePartsOnError = true    //出错的时候不删除已经上传的部分文件
	})
	if uploader == nil {
		return "", errors.New("NewUploaderWithClient failed")
	}

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   value,
	}

	if len(header) > 0 {
		if ctypes, ok := header["Content-Type"]; ok {
			ctypesStr := strings.Join(ctypes, "; ")
			if len(ctypesStr) > 3 {
				upParams.ContentType = &ctypesStr
			}
		}
	}

	result, err := uploader.UploadWithContext(ctx, upParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			l4g.Error("upload canceled due to timeout, %v\n", err)
		} else {
			l4g.Error("failed to upload object, %v\n", err)
		}
		return "", err
	}

	l4g.Info("successfully uploaded file to %s/%s\n", bucket, key)
	return result.Location, nil
}

//expire 与存储桶的超时策略 ，取哪个？？
func (this *S3Sdk) UpLoadEx(bucket, key string, value io.Reader, timeout int, expire int64) (string, error) {
	if timeout <= 0 {
		timeout = 120 //default 60s
	}
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancelFn()

	uploader := s3manager.NewUploaderWithClient(this.mS3Client, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024 // 64MB per part
		u.LeavePartsOnError = true    //出错的时候不删除已经上传的部分文件
	})
	if uploader == nil {
		return "", errors.New("NewUploaderWithClient failed")
	}

	upParams := &s3manager.UploadInput{
		Bucket:  aws.String(bucket),
		Key:     aws.String(key),
		Body:    value,
		Expires: aws.Time(time.Now().Add(time.Second * time.Duration(expire))),
	}

	result, err := uploader.UploadWithContext(ctx, upParams)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			l4g.Error("upload canceled due to timeout, %v\n", err)
		} else {
			l4g.Error("failed to upload object, %v\n", err)
		}
		return "", err
	}

	l4g.Info("successfully uploaded file to %s/%s\n", bucket, key)
	return result.Location, nil
}

//阻塞模式，用并发
//参数说明：key是文件名；value是写入流；timeout文件上传超时值，超时将中断传输。
//返回值：n，文件大小
func (this *S3Sdk) Download(bucket, key string, value io.WriterAt, timeout int) (int64, error) {
	if timeout <= 0 {
		timeout = 120 //default 60s
	}
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancelFn()

	downloader := s3manager.NewDownloaderWithClient(this.mS3Client, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
	})

	if downloader == nil {
		return -1, errors.New("NewDownloaderWithClient failed")
	}

	n, err := downloader.DownloadWithContext(ctx, value, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return n, err
	}

	l4g.Info("successfully download file %s/%s\n", bucket, key)
	return n, nil
}

type S3ObjectStatus struct {
	headObject *s3.HeadObjectOutput
	isExist    bool
}

func (this *S3ObjectStatus) Exist() bool {
	return this.isExist
}

func (this *S3ObjectStatus) Timeout() bool {
	//this.headObject.Expires
	if !this.isExist {
		return true
	}
	if this.headObject == nil {
		return true
	}
	if this.headObject.Expires == nil {
		return false
	}
	//fmt.Println(*this.headObject.Expires)
	loc, _ := time.LoadLocation("GMT")
	t, err := time.ParseInLocation("Thu, 02 Jan 2006 15:04:05 GMT", *this.headObject.Expires, loc)
	if err != nil {
		fmt.Println("sdks3 S3ObjectStatus", err)
	}
	//fmt.Println(t.Unix(), t.Before(time.Now()))
	return t.Before(time.Now())
}

func (this *S3ObjectStatus) Addr() string {
	if this.headObject == nil || this.headObject.WebsiteRedirectLocation == nil {
		return ""
	}
	return *this.headObject.WebsiteRedirectLocation
}

func (this *S3Sdk) Status(bucket, key string) (*S3ObjectStatus, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	out, err := this.mS3Client.HeadObject(input)
	if err != nil { //err是不是 包含不存在的情况
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			return &S3ObjectStatus{isExist: false}, nil
		}
		return nil, err
	}
	return &S3ObjectStatus{headObject: out, isExist: true}, nil
}

package aws

import (
	"bytes"
	"fmt"
	"github.com/BastionPay/bas-bkadmin-api/api-common"
	l4g "github.com/alecthomas/log4go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"net/http"
	"time"
)

var (
	Tools *common.Tools
)

const (
	MODULES = "kyc"
)

type (
	Aws struct {
		AccessKeyid     string
		SecretAccessKey string
		Token           string
		Region          string
		Bucket          string
	}
)

func init() {
	Tools = common.New()
}

func New(AccessKeyid, SecretAccessKey, Token, Region, Bucket string) *Aws {
	return &Aws{
		AccessKeyid:     AccessKeyid,
		SecretAccessKey: SecretAccessKey,
		Token:           Token,
		Region:          Region,
		Bucket:          Bucket,
	}
}

func (this Aws) UploadFile(file *multipart.FileHeader, model string) (string, error) {
	creds := credentials.NewStaticCredentials(
		this.AccessKeyid,
		this.SecretAccessKey,
		this.Token)

	_, err := creds.Get()
	if err != nil {
		return "", err
	}

	cfg := aws.NewConfig().WithRegion(this.Region).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	fileOpen, err := file.Open()
	if err != nil {
		return "", err
	}

	defer fileOpen.Close()

	size := file.Size

	buffer := make([]byte, size)

	fileOpen.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := this.generateFileName(model)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(this.Bucket),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}

	_, err = svc.PutObject(params)
	if err != nil {
		l4g.Warn(err, "upload file")
		return "", err
	}

	return path, nil
}

func (this Aws) generateFileName(model string) string {
	now := time.Now().UnixNano()
	random := Tools.GetRandomString(12)

	return fmt.Sprintf("/%s/%s/%d_%s", MODULES, model, now, Tools.MD5(random))
}

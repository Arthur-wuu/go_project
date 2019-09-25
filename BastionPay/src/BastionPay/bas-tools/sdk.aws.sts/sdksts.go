package sdk_aws_sts

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"
	"net/http"
	"time"
)

// 参数说明：区域，密钥Id，密钥，token（一般为空）
func NewStsSdk(region, accessKeyID, accessKey, accessToken string) *StsSdk {
	config := &aws.Config{
		Region: aws.String(region),
	}
	if (len(accessKeyID) != 0) && (len(accessKey) != 0) {
		config.Credentials = credentials.NewStaticCredentials(accessKeyID, accessKey, accessToken)
	}
	client := &StsSdk{
		mSess: session.Must(session.NewSession(config)),
	}
	client.mStsClient = sts.New(client.mSess)
	return client
}

type Policy struct {
	Version   string      `json:"Version,omitempty"`
	Statement []Statement `json:"Statement,omitempty"`
}

type Statement struct {
	Sid       string `json:"Sid,omitempty"`
	Effect    string `json:"Effect,omitempty"`
	Action    string `json:"Action,omitempty"`
	Resource  string `json:"Resource,omitempty"`
	Principal string `json:"Principal,omitempty"`
}

type StsSdk struct {
	mStsClient *sts.STS
	mSess      *session.Session
}

func (this *StsSdk) GetFederationToken(tokenExpire int64, name string, stats []Statement) (*sts.Credentials, error) {
	policy := &Policy{
		Version:   "2012-10-17",
		Statement: stats,
	}
	policyContent, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	input := &sts.GetFederationTokenInput{
		DurationSeconds: aws.Int64(tokenExpire),
		Name:            aws.String(name),
		Policy:          aws.String(string(policyContent)),
	}

	result, err := this.mStsClient.GetFederationToken(input)
	if err != nil {
		return nil, err
	}
	return result.Credentials, nil
}

func (this *StsSdk) SignGetUrl(url, region string, stat *Statement, expire int64) (string, error) {
	cred, err := this.GetFederationToken(expire, "pingzilao", []Statement{*stat})
	if err != nil {
		return "", err
	}
	newCred := credentials.NewStaticCredentials(*cred.AccessKeyId, *cred.SecretAccessKey, *cred.SessionToken)

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	sign := v4.NewSigner(newCred)
	_, err = sign.Presign(httpReq, nil, "s3", region, time.Second*time.Duration(expire), time.Now())
	if err != nil {
		return "", err
	}
	return httpReq.URL.String(), nil
}

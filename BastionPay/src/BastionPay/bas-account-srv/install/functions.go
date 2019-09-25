package install

import (
	"BastionPay/bas-account-srv/handler"
	"BastionPay/bas-api/apibackend"
	"BastionPay/bas-api/apibackend/v1/backend"
	"BastionPay/bas-api/utils"
	"BastionPay/bas-base/config"
	"BastionPay/bas-base/data"
	"encoding/json"
	"errors"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"io/ioutil"
	"os"
)

func BuildWebAdmin() (*backend.ReqUserRegister, *backend.ReqUserUpdateProfile, error) {
	uc := &backend.ReqUserRegister{}
	up := &backend.ReqUserUpdateProfile{}

	//0:普通用户 1:热钱包; 2:管理员
	uc.UserClass = 2
	// 100：普通管理员 200：创世管理员
	uc.Level = 200

	var input string

	// 1
	fmt.Println("请输入后台公钥路径: ")
	input = ""
	input = utils.ScanLine()
	pubKey, err := ioutil.ReadFile(input)
	if err != nil {
		return nil, nil, err
	}
	up.PublicKey = string(pubKey)

	// 2
	fmt.Println("请输入后台限制IP: ")
	input = ""
	input = utils.ScanLine()
	up.SourceIP = input

	// 3
	fmt.Println("请输入后台回调地址（后台可不填）: ")
	input = ""
	input = utils.ScanLine()
	up.CallbackUrl = input

	return uc, up, nil
}

func InstallBastionPay(dir string) error {
	var err error

	fi, err := os.Open(dir + "/" + config.BastionPayInstall)
	if err == nil {
		defer fi.Close()
		l4g.Info("BastionPay已经安装！！！")
		return nil
	}

	fmt.Println("***正在安装超级管理员，请谨慎操作***")

	l4g.Info("1. 创建BastionPay RSA密钥对...")
	bCreateNewRsa := false
	priKeyPath := dir + "/" + config.BastionPayPrivateKey
	pubKeyPath := dir + "/" + config.BastionPayPublicKey
	_, err = os.Open(priKeyPath)
	if err != nil {
		bCreateNewRsa = true
	}
	_, err = os.Open(pubKeyPath)
	if err != nil {
		bCreateNewRsa = true
	}

	if bCreateNewRsa {
		err = utils.RsaGen(2048, priKeyPath, pubKeyPath)
		if err != nil {
			return err
		}
		l4g.Info("~~~成功创建BastionPay RSA密钥对: %s", dir)
	} else {
		l4g.Info("!!!使用存在BastionPay RSA密钥对: %s", dir)
	}

	l4g.Info("2. 创建Web超级管理员账号")
	ackUc := backend.AckUserRegister{}
	ackUp := backend.AckUserUpdateProfile{}

	uc, up, err := BuildWebAdmin()
	if err != nil {
		return err
	}

	err = func() error {
		// register
		err = func() error {
			b, err := json.Marshal(*uc)
			if err != nil {
				return err
			}

			var req data.SrvRequest
			var res data.SrvResponse
			req.Argv.Message = string(b)
			handler.AccountInstance().Register(&req, &res)
			if res.Err != apibackend.NoErr {
				return errors.New(res.ErrMsg)
			}

			err = json.Unmarshal([]byte(res.Value.Message), &ackUc)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}

		// update profile
		err = func() error {
			b, err := json.Marshal(*up)
			if err != nil {
				return err
			}

			var req data.SrvRequest
			var res data.SrvResponse
			req.Argv.SubUserKey = ackUc.UserKey
			req.Argv.Message = string(b)
			handler.AccountInstance().UpdateProfile(&req, &res)
			if res.Err != apibackend.NoErr {
				return errors.New(res.ErrMsg)
			}

			err = json.Unmarshal([]byte(res.Value.Message), &ackUp)
			if err != nil {
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		return err
	}

	// write a tag file
	fo, err := os.Create(dir + "/" + config.BastionPayInstall)
	if err != nil {
		return err
	}
	defer fo.Close()

	fmt.Println("～～～Web超级管理员安装成功～～～")
	fmt.Println("请进行以下操作：")
	fmt.Printf("1. 记录Web超级管理员user_key(%s)\n", ackUc.UserKey)
	fmt.Printf("2. 将BastionPay公钥(%s)保存到Web后台\n", dir+"/public.pem")

	return nil
}

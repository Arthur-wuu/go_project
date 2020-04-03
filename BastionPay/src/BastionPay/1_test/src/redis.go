package main

import (
	//"encoding/json"
	//"fmt"
	//"reflect"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"strings"

	//"crypto/md5"
	"fmt"
	//"reflect"
	//"strconv"

	//"github.com/go-redis/redis"
)
type ParamsForRedis struct {
	Body           map[string]interface{}    `json:"body,omitempty"`
	Url            string                    `json:"url,omitempty"`
	AccountId      string                    `json:"account_id,omitempty"`
}
func main() {
	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})

	//paramsForRedis := ParamsForRedis{
	//	Body: map[string]interface{}{"user_id":"sdd"},
	//	Url: "http://test.url.com",
	//	AccountId: "account_id",
	//}
	//
	//paramsBytes, err := json.Marshal(paramsForRedis)
	//if err != nil {
	//	fmt.Println("err:",err)
	//}
	//
	//client.LPush("test11", "xx")
	//
	//
	//str, err := client.RPop("test11").Result()
	//str1, err1 := client.RPop("test11").Result()
	//if str1 == "" {
	//	fmt.Println("mei zhi")
	//
	//}
	//fmt.Printf("[%s],[%v]",str,err)
	//fmt.Printf("[%s],[%s]",str1,err1)
	//
	//paramsWx := new(ParamsForRedis)
	//
	//err = json.Unmarshal([]byte(str),paramsWx)

	//fmt.Println(paramsWx.Url)
	//body, _ := json.Marshal(paramsWx.Body)
	//fmt.Println(string(body))
	//fmt.Println(paramsWx.AccountId)

	//client.HMSet("test21", map[string]interface{}{"key1":"value1","key2":"value2"})
	//
	//v, err := client.HMGet("test21","key11111","key22222").Result()
	////fmt.Println("len:", len(v))
	//if err == nil && len(v) > 0 && v[0] != nil {
	//	 fmt.Println( v[0].(string))
	//}
	var a uint32
	a = 23
	str := fmt.Sprintf("%d",a)
	fmt.Println("str",str)

	//vs := v[0].(string)
	//fmt.Println("ok")
//	ssss := hashContent("11")
//	a:= hashCode("sas")
//fmt.Println(a)


	out, err := AesEncryptss( "AIK?H<>m4t+vSP)@", "%}U=@{HBwAGY[zS1", "C3A6772DA561D05BDCFFF81F1AA3DF32592DCE975A272D30F5CD84F1B2C8BA27EF54BDA19DAD531D2E87ADC578321F73F14FA6F56E247CCBF75650409E912E6415E0FB41D58219B42061DEA3A544D4F3A9D1ACDF5B5038DAE76049D5002BECFC3A42180A24CF0E2562EECFCF04687D4C57E2B3CCFB84B4A3399E117A469AC57104DA78104534572A36146F7F766F4FAC514C4D12ADC9CD71643677216E8527003D45173FE3A64E33BB2F1FBCC04BAF456C4617FBB9F6221C32EE4EC3B034CDB4D02FAFE406289985A77F18769C7C7A22962B56F208B8E9C47C12449A04DEF6E99D5FFFEEAD0443BC56F40BD7E77539EF47D6A323DDBF83B2A06B132E063EACFC003DE6D2522210DD72157E67923DDBA41CAE6823BA2B67A5157F16C4D28DBA56B280306024DC254E8285B22C4EC02AD28ADDFD9D394F9D13A4F4A5A42F8B3DBAD9B8B37F67B92A13A3D293CA0EDEE7D506D438DE4059F98F01CADB8AE7AC68184946889BACC3EFDE99ACE9CE56F723ABF68662F96968FADC9572C4BD93663D622C0786A7575AC682C927793F2F94E350374920DAF94C3B50310D44F65B14A7C2424BE4AD734224855660A1D934207E93B8ADC3456344848810A524AFDBC1AF6E1801B416D326BB43F580BA3163256B74DD05EBC94139846856F6B1738AEB9DB9AEF0141A276A34388F1F1DCF556B4E40BA2BECB66327CFF694F903703100F90D74BCC19DB6347D655E138E35E95252C8FDF219704A655F9C6AE90C2E70FF422FB0B33CC3D93A967A18ECFA40E573A0374EB52F54E632605ECDED23E72F27C73BC6D6E294F2EA256D6F0938113B0B695818534328EAF6686AE7B48A0F510BE5AF4FB29AE828138AF98B96B0FE3C413CE222E01D57F0600A5513B9F812DFFF3098C8D91B76CDAA40263A0AD1A3ED646F6CC5AABACD5C02290BB9FC5F6AD63EE3AB80B425806F4D5CE143F8C542BAD886AF5A62396E0C544321F72DF7EAB3DBE7AA9F399887677DCA36A3CC3E3ACEB66AAE74BCD295D134F324812E88B824FE4923212725EF69C254F8E7B2C56587DDC78F358EB9C067E84D4551AA699F60363887",false)
	OriginPkgnames := strings.Split(out, ",")
	for k, _ := range OriginPkgnames {
		fmt.Println("out: -- > ",OriginPkgnames[k])
	}

fmt.Println("errsss",err)
}
func AesEncryptss(aesKey, aesIv, inPut string, isEncrypt bool) (outPut string, err error) {

	if len(inPut) > 4096 {
		return "", fmt.Errorf("inPut size[%d] too big, encrypt will not process", len(inPut))
	}

	return AesEncryptWithBuffer(aesKey, aesIv, inPut, isEncrypt)
}

// AesEncryptWithBuffer ...
func AesEncryptWithBuffer(aesKey, aesIv, inPut string, isEncrypt bool) (outPut string, err error) {

	if len(aesKey) == 0 || len(aesIv) == 0 || len(inPut) == 0 {
		return "", errors.New("arguments error, aes_key or aes_iv, in_buf is NULL")
	}

	output, hasErr := hex2bin(inPut)
	if hasErr {
		return "", errors.New("aes encrypt error")
	}

	var outPutData []byte
	if isEncrypt == false {
		cry := &aesCryptor{
			key: []byte(aesKey),
			iv:  []byte(aesIv),
		}
		outPutData, err = cry.Decrypt([]byte(output))
		if err != nil {
			return "", err
		}
	}

	var outPutDataN []byte
	for _, out := range outPutData {
		if out == 32 {
			break
		}
		outPutDataN = append(outPutDataN, out)
	}

	return string(outPutDataN), nil
}

// AesEncry aes 加密
func AesEncry(aesKey, aesIv, inPut string) (outPut string, err error) {
	if len(aesKey) == 0 || len(aesIv) == 0 || len(inPut) == 0 {
		return "", errors.New("arguments error, aes_key or aes_iv, in_buf is NULL")
	}
	var outPutData string
	cry := &aesCryptor{
		key: []byte(aesKey),
		iv:  []byte(aesIv),
	}
	outPutData, err = cry.Encrypt([]byte(inPut))
	if err != nil {
		return "", err
	}
	return string(bin2hex(outPutData)), nil
}

func hex2bin(input string) (output []byte, hasErr bool) {
	if len(input) == 0 {
		return nil, true
	}
	outputData := make([]byte, len(input))
	var high byte
	var low byte
	for i := 0; i*2 < len(input); i++ {
		if i*2 >= len(input) || i*2+1 >= len(input) {
			return nil, true
		}
		high = hex2int(input[i*2])
		low = hex2int(input[i*2+1])
		if high == 0xff || low == 0xff {
			return nil, true
		}
		if ((high << 4) + (low & 0xff)) != 0 {
			outputData[i] = (high << 4) + (low & 0xff)
		}
	}

	for _, v := range outputData {
		output = append(output, v)
	}
	return
}

func bin2hex(input string) (output []byte) {
	if len(input) == 0 {
		return nil
	}
	output = make([]byte, len(input)*2)
	for i := 0; i < len(input); i++ {
		high := int2hex(input[i] >> 4)
		low := int2hex(input[i] & 0x0f)
		if high == '-' || low == '-' {
			return nil
		}
		output[i*2+0] = high
		output[i*2+1] = low
	}
	return
}

func hex2int(hex byte) byte {
	if hex >= 'A' && hex <= 'F' {
		return (10 + (hex - 'A'))
	} else if hex >= 'a' && hex <= 'f' {
		return (10 + (hex - 'a'))
	} else if hex >= '0' && hex <= '9' {
		return hex - '0'
	}
	return 0xff
}

func int2hex(data byte) byte {
	if data >= 10 && data <= 15 {
		return 'A' + (data - 10)
	} else if data >= 0 && data <= 9 {
		return data + '0'
	}
	return '-'
}

type aesCryptor struct {
	key []byte
	iv  []byte
}

// Decrypt aes decrypt
func (a *aesCryptor) Decrypt(src []byte) (data []byte, err error) {
	decrypted := make([]byte, len(src))

	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher(a.key)
	if err != nil {
		println(err.Error())
		return nil, err
	}
	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, a.iv)
	aesDecrypter.CryptBlocks(decrypted, src)
	return PKCS5Trimming(decrypted), nil
}

func (a *aesCryptor) Encrypt(str []byte) (string, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	str = ZeroPaddings(str, blockSize)
	var blockMode cipher.BlockMode
	blockMode = cipher.NewCBCEncrypter(block, a.iv)
	crypted := make([]byte, len(str))
	blockMode.CryptBlocks(crypted, str)
	return string(crypted), nil
}

// ZeroPaddings 0填充
func ZeroPaddings(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(ciphertext, padtext...)
}

// PKCS5Trimming 解包装
func PKCS5Trimming(encrypt []byte) []byte {

	return encrypt[:len(encrypt)]
}
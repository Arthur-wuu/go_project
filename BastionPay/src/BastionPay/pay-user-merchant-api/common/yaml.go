package common

////
////import (
////	"gopkg.in/yaml.v2"
////	"io/ioutil"
////	"log"
////)
////
////type Config struct {
////	Server struct {
////		Port  int
////		Debug bool
////	}
////	Db struct {
////		Type     string
////		Host     string
////		Port     int
////		User     string
////		Password string
////	}
////	Redis struct {
////		Host     string
////		Port     int
////		Password string
////	}
////	Token struct {
////		Secret    string
////		Algorithm string
////		PublicKey string
////	}
////	Exchange struct {
////		Host      string
////		Port      int
////		Appkey    string
////		Signature string
////	}
////	Log struct {
////		Type    string
////		Level   string
////		LogDir  string
////		LogFile string
////	}
////}
////
////func ParseYaml(file string, config *interface{}) error {
////	var(
////		err error
////		data []byte
////		//t interface{}
////		//tt Config
////	)
////	data, err = ioutil.ReadFile(file)
////	if err != nil {
////		log.Fatalf("Read file ERR:%v", err)
////		return err
////	}
////
////	err = yaml.Unmarshal([]byte(data), config)
////	if err != nil {
////		log.Fatalf("error: %v", err)
////		return err
////	}
////
////	return nil
////}
//
//package main
//
//import (
//"fmt"
//"log"
//
//"gopkg.in/yaml.v2"
//)
//
//var data = `
//language:
//  - zh-CN
//  - zh-TW
//  - en-US
//a: Easy!
//b:
//  c: 2
//  d: [3, 4]
//  e:
//    version: 1.0
//    name: >
//      Arch
//      Linux
//`
//
//// Note: struct fields must be public in order for unmarshal to
//// correctly populate the data.
//type T struct {
//	Language []string
//	A        string
//	B struct {
//		RenamedC int   `yaml:"c"`
//		D        []int `yaml:",flow"`
//		E struct {
//			Version string
//			Name    string
//		}
//	}
//}
//
//func main() {
//	t := T{}
//
//	err := yaml.Unmarshal([]byte(data), &t)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	fmt.Printf("--- t:\n%v\n\n", t)
//
//	d, err := yaml.Marshal(&t)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	fmt.Printf("--- t dump:\n%s\n\n", string(d))
//
//	m := make(map[interface{}]interface{})
//
//	err = yaml.Unmarshal([]byte(data), &m)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	fmt.Printf("--- m:\n%v\n\n", m)
//
//	d, err = yaml.Marshal(&m)
//	if err != nil {
//		log.Fatalf("error: %v", err)
//	}
//	fmt.Printf("--- m dump:\n%s\n\n", string(d))
//}

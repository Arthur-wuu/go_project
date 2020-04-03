package common

import (
	"archive/zip"
	"io/ioutil"
	"os"
)

func NewZip() *Zip {
	return &Zip{}
}

type Zip struct {
}

func (this *Zip) Compress(srcPath, desPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	fzip, err := os.Create(desPath)
	if err != nil {
		return err
	}
	w := zip.NewWriter(fzip)
	defer w.Close()

	fw, err := w.Create(file.Name())
	if err != nil {
		return err
	}
	filecontent, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}
	_, err = fw.Write(filecontent)
	return err
}

//func (this * Zip) CompressDir(srcDir, desDir string)  error {
//	f, err := ioutil.ReadDir(srcDir)
//	if err != nil {
//		return err
//	}
//	strings.Index()
//	fzip, _ := os.Create("img-50.zip")
//	w := zip.NewWriter(fzip)
//	defer w.Close()
//	for _, file := range f {
//		fw, _ := w.Create(file.Name())
//		filecontent, err := ioutil.ReadFile(dir + file.Name())
//		if err != nil {
//			fmt.Println(err)
//		}
//		n, err := fw.Write(filecontent)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(n)
//	}
//}

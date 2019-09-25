package utils

import (
	"testing"
	"os"
)

func TestGetRunDir(t *testing.T) {
	dir, err := GetRunDir()
	if err != nil || dir == ""{
		t.Error("TestGetRunDir failed")
	}
}

func TestGetCurrentDir(t *testing.T) {
	dir, err := GetCurrentDir()
	dir1, _ := os.Getwd()
	if err != nil{
		t.Error(err)
	}else if dir != dir1 {
		t.Error("TestGetCurrentDir failed")
	}
}

func TestGetAppDir(t *testing.T) {
	dir, err := GetAppDir()
	if err != nil || dir != "/Users/henly.liu/Library/Application Support"{
		t.Error("TestGetAppDir failed:", dir)
	}
}

func TestPathExists(t *testing.T) {
	type PathInfo struct{
		path string
		exist bool
	}
	var paths []PathInfo
	paths = append(paths, PathInfo{path:"/Users/henly.liu/workspace", exist:true})
	paths = append(paths, PathInfo{path:"/Users/henly.liu/work", exist:false})
	paths = append(paths, PathInfo{path:"/Users/henly.liu/workspace/private_wallet.pem", exist:true})
	paths = append(paths, PathInfo{path:"/Users/henly.liu/workspace/private.pem", exist:false})

	for _, v := range paths{
		b, _ := PathExists(v.path)
		if b != v.exist {
			t.Error("TestPathExists failed, ", v.path)
			break
		}
	}
}

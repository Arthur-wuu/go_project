package utils

import (
	"runtime/debug"
	l4g "github.com/alecthomas/log4go"
)

func PanicPrint()  {
	if err := recover(); err != nil {
		l4g.Error(string(debug.Stack()))
	}
}

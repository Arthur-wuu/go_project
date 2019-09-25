package utils

import (
	l4g "github.com/alecthomas/log4go"
	"runtime/debug"
)

func PanicPrint() {
	if err := recover(); err != nil {
		l4g.Error(string(debug.Stack()))
	}
}

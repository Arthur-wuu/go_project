package utils

import (
	. "BastionPay/bas-base/log/zap"
	"go.uber.org/zap"
	"runtime/debug"
)

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

package utils

import (
	"go.uber.org/zap"
	"runtime/debug"
	. "BastionPay/bas-base/log/zap"
)

func PanicPrint() {
	if err := recover(); err != nil {
		ZapLog().With(zap.Any("error", err)).Error(string(debug.Stack()))
	}
}

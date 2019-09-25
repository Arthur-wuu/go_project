package l4gmgr

var logWatchInt = 0
var logWacthStr = ""

func SetWatchInt(i int) {
	logWatchInt = i
}

func SetWatchString(str string) {
	logWacthStr = str
}

func IsWatchInt(condition int) bool {
	if condition == 0 {
		return false
	}
	return logWatchInt == condition
}

func IsWatchString(condition string) bool {
	return condition == logWacthStr
}

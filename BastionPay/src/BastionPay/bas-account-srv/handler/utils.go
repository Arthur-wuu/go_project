package handler

func IsBlankAuditeInfo(info string) bool {
	if len(info) == 0 {
		return true
	}
	for i := 0; i < len(info); i++ {
		if info[i] != ' ' {
			return false
		}
	}
	return true
}

package db

func GenRateFail2(succ, fail int)float32 {
	sum := succ + fail
	if sum == 0 {
		return 0
	}
	return float32(fail)/float32(sum)
}


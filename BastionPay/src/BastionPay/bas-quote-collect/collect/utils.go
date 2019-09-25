package collect

import "math/rand"

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

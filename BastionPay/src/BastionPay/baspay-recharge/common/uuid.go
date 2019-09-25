package common

import (
	"fmt"
	"github.com/satori/go.uuid"
	"strconv"
)

func (t *Tools) GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	return fmt.Sprintf("%s", ud)
}

func (t *Tools) GenerateUuidInt() int64 {
	ud := uuid.Must(uuid.NewV4())

	string := fmt.Sprintf("%s", ud)

	i, _ := strconv.ParseInt(string, 10, 64)
	return i
}

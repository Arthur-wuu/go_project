package common

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func (t *Tools) GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	return fmt.Sprintf("%s", ud)
}

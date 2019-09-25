package common

import (
	"github.com/satori/go.uuid"
	"fmt"
)

func (t *Tools) GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	return fmt.Sprintf("%s", ud)
}

func  GenerateUuid() string {
	ud := uuid.Must(uuid.NewV4())

	return fmt.Sprintf("%s", ud)
}
package main

import (
	"fmt"
	"github.com/BastionPay/bas-admin-api/common"
)

func main() {
	jj, err := common.NewJwt("nihao", "144000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(jj)
}

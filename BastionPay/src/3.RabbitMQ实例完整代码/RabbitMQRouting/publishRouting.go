package main

import (
	"fmt"
	"rabbitMq/3.RabbitMQ实例完整代码/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	imoocOne := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_one")
	imoocTwo := RabbitMQ.NewRabbitMQRouting("exImooc", "imooc_two")
	for i := 0; i <= 10; i++ {
		imoocOne.PublishRouting("Hello imooc one!" + strconv.Itoa(i))
		imoocTwo.PublishRouting("Hello imooc Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}

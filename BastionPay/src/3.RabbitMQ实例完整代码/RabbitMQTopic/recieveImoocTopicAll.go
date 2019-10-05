package main

import "rabbitMq/3.RabbitMQ实例完整代码/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQTopic("exImoocTopic","#")
	imoocOne.RecieveTopic()
}

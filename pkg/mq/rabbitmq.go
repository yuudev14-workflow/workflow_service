package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

var MQConn *amqp.Connection

func ConnectToMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logging.Logger.Panicf("%v", err)
	}

	MQConn = conn

	// defer mqConn.Close()
}

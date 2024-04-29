package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

var MQConn *amqp.Connection
var MQChannel *amqp.Channel

func ConnectToMQ() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logging.Logger.Panicf("%v", err)
	}
	logging.Logger.Info("connected to message queue...")

	MQConn = conn

	defer MQConn.Close()

	ch, chErr := MQConn.Channel()
	logging.Logger.Info("open channel in message queue...")

	if chErr != nil {
		logging.Logger.Panicf("%v", err)
	}

	MQChannel = ch

	defer MQChannel.Close()
}

package mq

import (
	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/environment"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

var (
	MQConn        *amqp.Connection
	MQChannel     *amqp.Channel
	SenderQueue   amqp.Queue
	ReceiverQueue amqp.Queue
)

func ConnectToMQ() {
	logging.Logger.Infof("connecting to message queue %v...", environment.Settings.MQ_URL)
	conn, err := amqp.Dial(environment.Settings.MQ_URL)
	if err != nil {
		logging.Logger.Panicf("%v", err)
	}
	logging.Logger.Info("connected to message queue...")

	MQConn = conn

	ch, chErr := MQConn.Channel()
	logging.Logger.Info("open channel in message queue...")

	if chErr != nil {
		logging.Logger.Panicf("%v", chErr)
	}

	MQChannel = ch
	declareQueues(MQChannel)

}

func declareQueues(ch *amqp.Channel) {
	logging.Logger.Info("Declaring queues")
	declareSenderQueue(ch)
	declareReceiverQueue(ch)
}

func declareSenderQueue(ch *amqp.Channel) {
	logging.Logger.Info("Declaring sender queue")
	// Declare a queue
	q, err := ch.QueueDeclare(
		environment.Settings.SenderQueueName, // name
		true,                                 // durable
		false,                                // delete when unused
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	if err != nil {
		logging.Logger.Panicf("%v", err)
	}
	SenderQueue = q
}

func declareReceiverQueue(ch *amqp.Channel) {
	logging.Logger.Info("Declaring receiver queue")
	// Declare a queue
	q, err := ch.QueueDeclare(
		environment.Settings.ReceiverQueueName, // name
		true,                                   // durable
		false,                                  // delete when unused
		false,                                  // exclusive
		false,                                  // no-wait
		nil,                                    // arguments
	)
	if err != nil {
		logging.Logger.Panicf("%v", err)
	}
	ReceiverQueue = q
}

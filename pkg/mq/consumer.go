package mq

import "github.com/yuudev14-workflow/workflow-service/pkg/logging"

func Listen() {
	msgs, err := MQChannel.Consume(
		ReceiverQueue.Name, // queue
		"",                 // consumer
		false,              // auto-acknowledge
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // arguments
	)

	if err != nil {
		panic("error in consuming a queue")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logging.Logger.Infof("Received a message: %s", d.Body)
		}
	}()

	logging.Logger.Info("Listening to message queue")
	<-forever
}

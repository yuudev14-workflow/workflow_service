package mq

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/utils"
)

func PrepareMessage(message utils.WorkflowData, currentNode string) {
	graph := message.Graph

	currentQueue := []string{
		"A",
	}

	// visited := message.Visited

	// queue := message.CurrentQueue

	// Publish a message to the queue
	body := utils.WorkflowData{
		Graph:        graph,
		CurrentNode:  currentNode,
		CurrentQueue: currentQueue,
		Visited:      currentQueue,
	}

	jsonData, jsonErr := json.Marshal(body)

	if jsonErr != nil {
		logging.Logger.Warnf("Error decoding JSON: %v", jsonErr)
	}
	err := MQChannel.Publish(
		"",               // exchange
		SenderQueue.Name, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jsonData),
		})
	if err != nil {
		logging.Logger.Errorf("MQ publish error: %v", jsonErr)
	}
}
func Listen() {
	msgs, err := MQChannel.Consume(
		ReceiverQueue.Name, // queue
		"",                 // consumer
		true,               // auto-acknowledge
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
			var data utils.WorkflowData

			err := json.Unmarshal(d.Body, &data)
			if err != nil {
				logging.Logger.Warnf("Error decoding JSON: %v", err)
			}
			logging.Logger.Infof("Received a message: %s", data)

		}
	}()

	logging.Logger.Info("Listening to message queue")
	<-forever
}

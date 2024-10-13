package mq

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/utils"
)

func PrepareMessage(message utils.WorkflowData) {
	for _, node := range message.Graph[message.CurrentNode] {
		// check if nodes with node destinatin in the database is already finished with success
		// if all is finish, publish the message
		logging.Sugar.Infof("Node: %s", node)
		body := utils.WorkflowData{
			Graph:        message.Graph,
			CurrentNode:  node,
			CurrentQueue: message.CurrentQueue,
			Visited:      message.Visited,
		}

		jsonData, jsonErr := json.Marshal(body)

		if jsonErr != nil {
			logging.Sugar.Warnf("Error decoding JSON: %v", jsonErr)
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
			logging.Sugar.Errorf("MQ publish error: %v", jsonErr)
		}
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
				logging.Sugar.Warnf("Error decoding JSON: %v", err)
			}
			logging.Sugar.Infof("Received a message: %s", data)
			// if all nodes in graph is finish with success dont prepare message
			// if status is failed, dont prepare message, remove all the message
			PrepareMessage(data)
		}
	}()

	logging.Sugar.Info("Listening to message queue")
	<-forever
}

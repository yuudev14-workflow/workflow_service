package mq

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
	"github.com/yuudev14-workflow/workflow-service/pkg/utils"
)

type TaskMessage struct {
	Nodes map[string][]string `json:"nodes"`
	Edges []repository.Edges  `json:"edges"`
}

func SendTaskMessage(graph TaskMessage) error {
	jsonData, jsonErr := json.Marshal(graph)

	if jsonErr != nil {
		return jsonErr
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
		return err
	}

	logging.Sugar.Infow("successfully pushed the message", "jsonData", string(jsonData))
	return nil

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
			// PrepareMessage(data)
		}
	}()

	logging.Sugar.Info("Listening to message queue")
	<-forever
}

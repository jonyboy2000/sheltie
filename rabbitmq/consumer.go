package rabbitmq

import (
	"time"

	p "../parser"
	util "../utility"
	"github.com/streadway/amqp"
)

// Consumer => Exported
// recieve a message
func Consumer(
	conStr string, // sample: "amqp://guest:guest@127.0.0.1:5672/"
	queueName string) {

	//============================================
	// Connect to Server
	//============================================
	conn, err := amqp.Dial(conStr)
	if err != nil {
		util.Log("(RabbitMQ=>Consumer) Failed to connect to RabbitMQ", err)

		defer conn.Close()
		return
	}

	//============================================
	// Create Channel
	//============================================
	ch, err := conn.Channel()
	if err != nil {
		util.Log("(RabbitMQ=>Consumer) Failed to open a channel", err)

		defer ch.Close()
		defer conn.Close()
		return
	}

	//============================================
	// Create Queue
	//============================================
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		util.Log("(RabbitMQ=>Consumer) Failed to declare a queue", err)

		defer ch.Close()
		defer conn.Close()
		return
	}

	//============================================
	// Get Message
	//============================================
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		util.Log("(RabbitMQ=>Consumer) Failed to register a consumer", err)

		defer ch.Close()
		defer conn.Close()
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			res := p.Parse(d)
			if res { // on success => send ack to rmq
				d.Ack(false)
			} else { // on fail => reject message & requeue (at first time)
				d.Reject(!d.Redelivered)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}()

	<-forever

	//============================================
	// Close
	//============================================
	defer ch.Close()
	defer conn.Close()
}

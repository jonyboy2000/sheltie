package rabbitmq

import (
	util "../utility"
	"github.com/streadway/amqp"
)

// Publisher => Exported
// send a message
func Publisher(
	conStr string, // sample: "amqp://guest:guest@127.0.0.1:5672/"
	queueName string,
	msg string) {

	//============================================
	// Connect to Server
	//============================================
	conn, err := amqp.Dial(conStr)
	if err != nil {
		util.Log("(RabbitMQ=>Publisher) Failed to connect to RabbitMQ", err)

		defer conn.Close()
		return
	}

	//============================================
	// Create Channel
	//============================================
	ch, err := conn.Channel()
	if err != nil {
		util.Log("(RabbitMQ=>Publisher) Failed to open a channel", err)

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
		util.Log("(RabbitMQ=>Publisher) Failed to declare a queue", err)

		defer ch.Close()
		defer conn.Close()
		return
	}

	//============================================
	// Send Message
	//============================================
	body := msg
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent, // 1=Transient(non-persistent), 2=Persistent
		})

	if err != nil {
		util.Log("(RabbitMQ=>Publisher) Failed to publish a message", err)
	}

	//============================================
	// Close
	//============================================
	defer ch.Close()
	defer conn.Close()
}

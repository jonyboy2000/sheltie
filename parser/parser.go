package parser

import (
	st "../structure"
	"github.com/streadway/amqp"
)

var Config st.Config

// Parse rabbitmq message & execute coresponding task => Exported
func Parse(msg amqp.Delivery) bool {

	// Do something

	return false
}

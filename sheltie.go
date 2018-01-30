package main

import (
	"encoding/json"
	"os"
	"time"

	p "./parser"
	rabCons "./rabbitmq/consumer"
	st "./structure"
	util "./utility"
	"github.com/streadway/amqp"
)

func main() {
	util.CreateDirectory("log")

	//============================================
	// Read Config
	//============================================
	content := util.ReadFile("config.json")
	if content == nil {
		util.Log("(Main) Config file is empty", nil)
		os.Exit(1)
	}

	var conf st.Config
	err := json.Unmarshal(content, &conf)
	if err != nil {
		util.Log("(Main) Error in parse 'config.json' file", err)
		os.Exit(1)
	}

	//============================================
	// Set Config
	//============================================
	conStr := "amqp://" + conf.Rabbit_User + ":" + conf.Rabbit_Pass + "@" + conf.Rabbit_Host + ":" + conf.Rabbit_Port + "/"

	p.Config = conf
	p.RabbitConnString = conStr

	var msg amqp.Delivery
	msg.Body = []byte("{\"type\": \"ReadFile\",\"file\": \"$POD_Directory/test.yaml\",\"commands\": [[]],\"args\": []}")
	p.Parse(msg)

	//============================================
	// Initialize RabbitMQ Listener
	//============================================
	rabCons.Consumer(conStr, conf.Rabbit_CommandQueue)

	//============================================
	// Main Loop
	//============================================
	for {
		time.Sleep(1 * time.Second)
	}
}

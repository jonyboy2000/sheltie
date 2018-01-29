package main

import (
	"encoding/json"
	"os"
	"time"

	p "./parser"
	rabbit "./rabbitmq"
	st "./structure"
	util "./utility"
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
		util.Log("(Main) Error in parse 'config.json' file", nil)
		os.Exit(1)
	}

	// Set config on parser package
	p.Config = conf

	//============================================
	// Initialize RabbitMQ Listener
	//============================================
	conStr := "amqp://" + conf.Rabbit_User + ":" + conf.Rabbit_Pass + "@" + conf.Rabbit_Host + ":" + conf.Rabbit_Port + "/"
	rabbit.Consumer(conStr, conf.Rabbit_CommandQueue)

	//============================================
	// Main Loop
	//============================================
	for {
		time.Sleep(1 * time.Second)
	}
}

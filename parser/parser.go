package parser

import (
	"encoding/json"
	"strings"

	k "../k8s"
	rabPub "../rabbitmq/publisher"
	st "../structure"
	util "../utility"
	"github.com/streadway/amqp"
)

// Config object
var Config st.Config

// RabbitConnString => RabbitMQ connection string
var RabbitConnString string

// ResolveParameter => Exported
// Replace existed alias in parameters with their values
func ResolveParameter(param string) string {
	alias := Config.Alias
	for i := 0; i < len(alias); i++ {
		param = strings.Replace(param, "$"+alias[i]["key"], alias[i]["value"], -1)
	}

	return param
}

// Parse rabbitmq message & execute coresponding task => Exported
func Parse(msg amqp.Delivery) bool {
	ret := false
	var err error

	//===================================================
	// Parse Command
	//===================================================
	var command st.Command
	err = json.Unmarshal([]byte(ResolveParameter(string(msg.Body))), &command)
	if err != nil {
		util.Log("(Parser=>Parse) Error in unmarshalling command message", err)
		return ret
	}

	//===================================================
	// Do
	//===================================================
	switch command.Type {
	//=====================
	// Create K8S POD
	//=====================
	case st.CommandTypeK8sPod:
		ret = k.CreatePOD(
			Config.K8s_ConfigFilePath, // kubeConfigPath
			command.File,              // podFilePath
			Config.K8s_NameSpace,      // namespace
			Config.K8s_NodeName,       // nodeName
			Config.K8s_Priority,       // priority
			Config.K8s_Labels,         // labels
			Config.K8s_Annotations,    // annotations
			command.Commands,          // commands
			command.Args)              // args
		break

	//=====================
	// Read File
	//=====================
	case st.CommandTypeReadFile:
		content := util.ReadFile(command.File)
		if content == nil {
			content = []byte("")
		}

		ret = rabPub.Publisher(
			RabbitConnString,          // conStr
			Config.Rabbit_ResultQueue, // queueName
			content)                   // msg
		break
	}

	//===================================================
	// Return
	//===================================================
	return ret
}

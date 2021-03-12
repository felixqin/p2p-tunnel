package contacts

import (
	"encoding/json"
	"log"
)

type messageFrom struct {
	User   string `json:"user"`
	Client string `json:"client"`
}

type message struct {
	Command string          `json:"command"`
	Data    interface{} `json:"data"`
	From    *messageFrom     `json:"from"`
}

func parseUserMessage(payload []byte) (string, string, interface{}, error) {
	var msg message
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		log.Println("unmarshal message failed!", err)
		return "", "", nil, err
	}

	return msg.From.Client, msg.Command, msg.Data, nil
}

func parseClientMessage(payload []byte) (string, string, interface{}, error) {
	var msg message
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		log.Println("unmarshal message failed!", err)
		return "", "", nil, err
	}

	return msg.From.User, msg.Command, msg.Data, nil
}

func sendMessageToUser(to string, cmd string, data interface{}) error {
	payload, err := json.Marshal(&message{
		Command: cmd,
		Data: data,
		From: &messageFrom{
			User:   configure.Username,
			Client: configure.clientId,
		},
	})
	if err != nil {
		return err
	}

	err = mqttSendToUser(to, payload)
	if err != nil {
		return err
	}

	return nil
}

func sendMessageToClient(to string, cmd string, data interface{}) error {
	payload, err := json.Marshal(&message{
		Command: cmd,
		Data: data,
		From: &messageFrom{
			User:   configure.Username,
			Client: configure.clientId,
		},
	})
	if err != nil {
		return err
	}

	err = mqttSendToClient(to, payload)
	if err != nil {
		return err
	}

	return nil
}

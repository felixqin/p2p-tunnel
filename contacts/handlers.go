package contacts

import (
	"encoding/json"
	"log"
)

var offerHandler func(fromClient string, offer *Offer)
var answerHandler func(fromClient string, answer *Answer)

func handleServerConnected() {
	log.Println("server connected!")
	sendMessageToUser(configure.Username, "queryClientInfo", nil)
}

func handleServerConnectionLost() {
	log.Println("server connection lost!")
}

func handleToUserMessage(uid string, payload []byte) {
	clientId, cmd, data, err := parseUserMessage(payload)
	log.Println("handleToUserMessage, from user:", uid, "client:", clientId)
	if err != nil {
		log.Println("parse user message failed!", err)
		return
	}

	switch cmd {
	case "queryClientInfo":
		err = onQueryClientInfo(uid)
	case "notifyClientInfo":
		err = onNotifyClientInfo(data)
	}

	if err != nil {
		log.Println("on message failed!", err)
		return
	}
}

func handleToClientMessage(clientId string, payload []byte) {
	uid, cmd, data, err := parseClientMessage(payload)
	log.Println("handleToClientMessage, from user:", uid, "client:", clientId)
	if err != nil {
		log.Println("parse user message failed!", err)
		return
	}

	switch cmd {
	case "sendOffer":
		err = onSendOffer(clientId, data)
	case "sendAnswer":
		err = onSendAnswer(clientId, data)
	}

	if err != nil {
		log.Println("on message failed!", err)
		return
	}
}

func onQueryClientInfo(toUser string) error {
	log.Println("on query client info:")
	return sendMessageToUser(toUser, "notifyClientInfo", &struct {
		Name   string `json:"name"`
		Client string `json:"client"`
		User   string `json:"user"`
	}{
		Name:   configure.Name,
		Client: configure.clientId,
		User:   configure.Username,
	})
}

func onNotifyClientInfo(data []byte) error {
	log.Println("on notify client info, data:", data)
	var info struct {
		Name   string `json:"name"`
		Client string `json:"client"`
		User   string `json:"user"`
	}
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}

	addContact(&Contact{
		Name:     info.Name,
		ClientId: info.Client,
		Owner:    info.User,
	})
	if err != nil {
		return err
	}

	return nil
}

func onSendOffer(fromClient string, data []byte) error {
	var offer Offer
	err := json.Unmarshal(data, &offer)
	if err != nil {
		return err
	}

	offerHandler(fromClient, &offer)
	return nil
}

func onSendAnswer(fromClient string, data []byte) error {
	var answer Answer
	err := json.Unmarshal(data, &answer)
	if err != nil {
		return err
	}

	answerHandler(fromClient, &answer)
	return nil
}

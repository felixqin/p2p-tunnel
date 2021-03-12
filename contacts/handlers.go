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

func handleToUserMessage(toUser string, payload []byte) {
	from, cmd, data, err := parseUserMessage(payload)
	log.Println("handleToUserMessage, from:", from)
	if err != nil {
		log.Println("parse user message failed!", err)
		return
	}

	switch cmd {
	case "queryClientInfo":
		err = onQueryClientInfo(from.User)
	case "notifyClientInfo":
		err = onNotifyClientInfo(data)
	}

	if err != nil {
		log.Println("on message failed!", err)
		return
	}
}

func handleToClientMessage(toClientId string, payload []byte) {
	from, cmd, data, err := parseClientMessage(payload)
	log.Println("handleToClientMessage, from:", from)
	if err != nil {
		log.Println("parse user message failed!", err)
		return
	}

	switch cmd {
	case "sendOffer":
		err = onSendOffer(from.Client, data)
	case "sendAnswer":
		err = onSendAnswer(from.Client, data)
	}

	if err != nil {
		log.Println("on message failed!", err)
		return
	}
}

func onQueryClientInfo(fromUser string) error {
	log.Println("on query client info")
	return sendMessageToUser(fromUser, "notifyClientInfo", &struct {
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
	// log.Printf("on notify client info, data: %s\n", data)
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
	// log.Printf("on send offer, data: %s\n", data)
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

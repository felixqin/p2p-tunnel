package contacts

import (
	"log"
)

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
	uid, _, _, err := parseClientMessage(payload)
	log.Println("handleToClientMessage, from user:", uid, "client:", clientId)
	if err != nil {
		log.Println("parse user message failed!", err)
		return
	}
}

func onQueryClientInfo(toUser string) error {
	log.Println("on query client info:")
	return sendMessageToUser(toUser, "notifyClientInfo", &struct{
		Name   string `json:"name"`
		Client string `json:"client"`
		User   string `json:"user"`
	}{
		Name:   configure.Name,
		Client: configure.clientId,
		User:   configure.Username,
	})
}

func onNotifyClientInfo(data interface{}) error {
	log.Println("on notify client info, data:", data)
	return nil
}

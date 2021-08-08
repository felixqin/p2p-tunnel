package main

import (
	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
)

var answerHandlers = map[string]func(sdp string){}

func handleAnswer(fromClient string, answer *contacts.Answer) {
	handler := answerHandlers[fromClient]
	if handler != nil {
		handler(answer.Sdp)
	}
}

func makeOfferSender(contact string) tunnel.OfferSender {
	return func(sdp string, answerHandler func(sdp string)) error {
		serverContact, err := contacts.FindContact(contact)
		if err != nil {
			return err
		}

		answerHandlers[serverContact.ClientId] = answerHandler
		return contacts.SendOffer(serverContact.ClientId, &contacts.Offer{
			Sdp: sdp,
		})
	}
}

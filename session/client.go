package session

import (
	"log"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
	"github.com/hashicorp/yamux"
)

type Client struct {
	tunnel  *tunnel.Client
	stream  *tunnel.Stream
	session *yamux.Session
}

var answerHandlers = map[string]func(sdp string){}
var clientSessions = []*Client{}

func init() {
	contacts.HandleAnswerFunc(handleAnswer)
}

func handleAnswer(fromClient string, answer *contacts.Answer) {
	handler := answerHandlers[fromClient]
	if handler != nil {
		handler(answer.Sdp)
	}
}

func makeOfferSender(nodeClientId string) tunnel.OfferSender {
	return func(sdp string, answerHandler func(sdp string)) error {
		log.Println("send offer to", nodeClientId)
		answerHandlers[nodeClientId] = answerHandler
		return contacts.SendOffer(nodeClientId, &contacts.Offer{
			Sdp: sdp,
		})
	}
}

// Connect create and start tunnel client
func Connect(nodeClientId string) error {
	client := &Client{}
	clientSessions = append(clientSessions, client)

	offerSender := makeOfferSender(nodeClientId)
	client.tunnel = tunnel.NewClient(&iceServers)
	return client.tunnel.Open(offerSender, func(stream *tunnel.Stream) {
		log.Println("proxy, to create yamux client ...")
		session, err := yamux.Client(stream, nil)
		if err != nil {
			log.Println("proxy, create yamux client failed!", err)
			return
		}

		client.session = session
		log.Println("client tunnel create success!!!")
	})
}

func (c *Client) Close() error {
	if c.session != nil {
		c.session.Close()
		c.session = nil
	}

	if c.stream != nil {
		c.stream.Close()
		c.stream = nil
	}

	if c.tunnel != nil {
		c.tunnel.Close()
		c.tunnel = nil
	}

	return nil
}

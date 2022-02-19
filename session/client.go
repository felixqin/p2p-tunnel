package session

import (
	"log"

	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/yamux"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
)

type Client struct {
	name    string
	tunnel  *tunnel.Client
	stream  *tunnel.Stream
	session *yamux.Session
}

var answerHandlers = map[string]func(sdp string){}
var clientSessions mapset.Set

func init() {
	clientSessions = mapset.NewSet()
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

func NewClient(name string) *Client {
	c := &Client{name: name}
	clientSessions.Add(c)
	return c
}

// Connect create and start tunnel client
func (c *Client) Connect(nodeClientId string) error {
	offerSender := makeOfferSender(nodeClientId)
	c.tunnel = tunnel.NewClient(&iceServers)
	return c.tunnel.Open(offerSender, func(stream *tunnel.Stream) {
		c.stream = stream
		log.Println("proxy, to create yamux client ...")
		session, err := yamux.Client(stream, nil)
		if err != nil {
			log.Println("proxy, create yamux client failed!", err)
			return
		}

		c.session = session
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

	clientSessions.Remove(c)
	return nil
}

package session

import (
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/yamux"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
)

type Client struct {
	Name          string
	Status        string
	toClientId    string ///< 对方的 Node Client ID
	answerHandler func(sdp string)
	tunnel        *tunnel.Client
	stream        *tunnel.Stream
	session       *yamux.Session
}

var answerHandlers = map[string]func(sdp string){}
var clientSessions mapset.Set

func init() {
	clientSessions = mapset.NewSet()
	contacts.HandleAnswerFunc(handleAnswer)
}

func NewClient(name string) *Client {
	c := &Client{Name: name, Status: "INIT"}
	clientSessions.Add(c)
	return c
}

func FindClient(name string) *Client {
	var client *Client
	clientSessions.Each(func(elem interface{}) bool {
		// fmt.Println("elem:", elem.(*Client))
		if elem.(*Client).Name == name {
			client = elem.(*Client)
			return true
		}

		return false
	})

	return client
}

func findClientByToNodeClientId(toNodeClientId string) *Client {
	var client *Client
	clientSessions.Each(func(elem interface{}) bool {
		if elem.(*Client).toClientId == toNodeClientId {
			client = elem.(*Client)
			return true
		}

		return false
	})

	return client
}

func DumpClients() {
	clientSessions.Each(func(elem interface{}) bool {
		// fmt.Println("elem:", elem.(*Client))
		client := elem.(*Client)
		fmt.Printf("%-12v%-12v\n", client.Name, client.Status)
		return false
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

// Connect create and start tunnel client
func (c *Client) Connect(toNodeClientId string) error {
	c.toClientId = toNodeClientId
	c.tunnel = tunnel.NewClient(&iceServers)
	return c.tunnel.Open(c.sendOffer, c.handleStreamOpen)
}

func (c *Client) sendOffer(sdp string, answerHandler func(sdp string)) error {
	log.Println("send offer to", c.toClientId)
	c.answerHandler = answerHandler
	c.Status = "OFFERING"
	return contacts.SendOffer(c.toClientId, &contacts.Offer{
		Sdp: sdp,
	})
}

func handleAnswer(fromClient string, answer *contacts.Answer) {
	client := findClientByToNodeClientId(fromClient)
	if client == nil {
		return
	}

	client.Status = "ANSWERED"
	if client.answerHandler != nil {
		client.answerHandler(answer.Sdp)
	}
}

func (c *Client) handleStreamOpen(stream *tunnel.Stream) {
	c.stream = stream
	c.Status = "STREAMED"
	log.Println("proxy, to create yamux client ...")
	session, err := yamux.Client(stream, nil)
	if err != nil {
		log.Println("proxy, create yamux client failed!", err)
		return
	}

	c.session = session
	c.Status = "CONNECTED"
	log.Println("client tunnel create success!!!")
}

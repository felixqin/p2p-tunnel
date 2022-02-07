package session

import (
	"log"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
	"github.com/hashicorp/yamux"
)

type Server struct {
	tunnel  *tunnel.Server
	stream  *tunnel.Stream
	session *yamux.Session
}

var serverSessions = []*Server{}

func init() {
	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		server := &Server{}
		serverSessions = append(serverSessions, server)

		answerSender := func(sdp string) error {
			log.Println("send answer to", fromClient)
			return contacts.SendAnswer(fromClient, &contacts.Answer{
				Sdp: sdp,
			})
		}

		server.tunnel = tunnel.NewServer(&iceServers) // 为每个 offer 创建一个 server 实例
		err := server.tunnel.Open(offer.Sdp, answerSender, func(stream *tunnel.Stream) {
			log.Println("stub, start yamux server ...")
			session, err := yamux.Server(stream, nil)
			if err != nil {
				log.Println("stub, create yamux server failed!", err)
				return
			}

			server.session = session
			log.Println("server tunnel create success!!!")
			go stubServe(session)
		})

		if err != nil {
			log.Println("open tunnel server failed!", err)
		}
	})
}

func (c *Server) Close() error {
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

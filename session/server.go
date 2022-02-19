package session

import (
	"log"

	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/yamux"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
)

type Server struct {
	tunnel  *tunnel.Server
	stream  *tunnel.Stream
	session *yamux.Session
}

var serverSessions mapset.Set

func init() {
	serverSessions = mapset.NewSet()

	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		answerSender := func(sdp string) error {
			log.Println("send answer to", fromClient)
			return contacts.SendAnswer(fromClient, &contacts.Answer{
				Sdp: sdp,
			})
		}

		s := NewServer()
		s.tunnel = tunnel.NewServer(&iceServers) // 为每个 offer 创建一个 server 实例
		err := s.tunnel.Open(offer.Sdp, answerSender, func(stream *tunnel.Stream) {
			s.stream = stream
			log.Println("stub, start yamux server ...")
			session, err := yamux.Server(stream, nil)
			if err != nil {
				log.Println("stub, create yamux server failed!", err)
				return
			}

			s.session = session
			log.Println("server tunnel create success!!!")
			go stubServe(session)
		})

		if err != nil {
			log.Println("open tunnel server failed!", err)
			s.Close()
			return
		}

	})
}

func NewServer() *Server {
	s := &Server{}
	serverSessions.Add(s)
	return s
}

func (s *Server) Close() error {
	if s.session != nil {
		s.session.Close()
		s.session = nil
	}

	if s.stream != nil {
		s.stream.Close()
		s.stream = nil
	}

	if s.tunnel != nil {
		s.tunnel.Close()
		s.tunnel = nil
	}

	serverSessions.Remove(s)
	return nil
}

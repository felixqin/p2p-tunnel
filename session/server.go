package session

import (
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/yamux"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
)

type Server struct {
	Status       string
	fromClientId string
	tunnel       *tunnel.Server
	stream       *tunnel.Stream
	session      *yamux.Session
}

var serverSessions mapset.Set

func init() {
	serverSessions = mapset.NewSet()
	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		s := NewServer()
		err := s.Serve(fromClient, offer.Sdp)
		if err != nil {
			log.Println("open tunnel server failed!", err)
			s.Close()
			return
		}
	})
}

func NewServer() *Server {
	s := &Server{Status: "INIT"}
	serverSessions.Add(s)
	return s
}

func DumpServers() {
	serverSessions.Each(func(elem interface{}) bool {
		server := elem.(*Server)
		fmt.Printf("%-12v%-12v\n", server.fromClientId, server.Status)
		return false
	})
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

func (s *Server) Serve(fromNodeClientId string, sdp string) error {
	s.fromClientId = fromNodeClientId
	s.tunnel = tunnel.NewServer(&iceServers) // 为每个 offer 创建一个 server 实例
	return s.tunnel.Open(sdp, s.sendAnswer, s.handleStreamOpen)
}

func (s *Server) sendAnswer(sdp string) error {
	s.Status = "ANSWERING"
	log.Println("send answer to", s.fromClientId)
	return contacts.SendAnswer(s.fromClientId, &contacts.Answer{
		Sdp: sdp,
	})
}

func (s *Server) handleStreamOpen(stream *tunnel.Stream) {
	s.stream = stream
	s.Status = "STREAMED"
	log.Println("stub, start yamux server ...")
	session, err := yamux.Server(stream, nil)
	if err != nil {
		log.Println("stub, create yamux server failed!", err)
		return
	}

	s.session = session
	s.Status = "CONNECTED"
	log.Println("server tunnel create success!!!")
	go stubServe(session)
}

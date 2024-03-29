package tunnel

import (
	"log"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

type Server struct {
	pc           *webrtc.RTCPeerConnection
	stream       *Stream
	onStreamOpen func(*Stream)
}

type AnswerSender func(sdp string) error

// Open(StubMessager) error

func NewServer(iceopts *IceServers) *Server {
	s := &Server{stream: newStream("stub")}

	pc, err := newWebRTC(iceopts)
	if err != nil {
		log.Println("stub, rtc error:", err)
		return nil
	}

	pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
		log.Print("stub, pc ice state change:", state)
	})

	pc.OnDataChannel(func(dc *webrtc.RTCDataChannel) {
		log.Print("stub, OnDataChannel", dc)
		s.stream.Open(dc, func() {
			if s.onStreamOpen != nil {
				s.onStreamOpen(s.stream)
			}
		})
	})

	s.pc = pc
	return s
}

func (s *Server) Open(sdp string, sender AnswerSender, onStreamOpen func(*Stream)) error {
	log.Println("stub to set remote sdp:", sdp)
	s.onStreamOpen = onStreamOpen
	err := s.pc.SetRemoteDescription(webrtc.RTCSessionDescription{
		Type: webrtc.RTCSdpTypeOffer,
		Sdp:  sdp,
	})
	if err != nil {
		log.Println("stub, set remote sdp failed!", err)
		return err
	}

	log.Println("stub set remote sdp success! then to send answer to proxy ...")
	answer, err := s.pc.CreateAnswer(nil)
	if err != nil {
		log.Println("stub, create answer failed!", err)
		return err
	}

	err = sender(answer.Sdp)
	if err != nil {
		log.Println("stub, send answer failed!", err)
		return err
	}

	log.Println("stub send answer success!")
	return nil
}

func (s *Server) Close() error {
	log.Println("stub close")
	s.stream.Close()
	s.pc.Close()
	return nil
}

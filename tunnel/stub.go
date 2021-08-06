package tunnel

import (
	"log"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

type Stub struct {
	pc     *webrtc.RTCPeerConnection
	stream *Stream
}

type AnswerSender func(sdp string) error

// Open(StubMessager) error
// Stream() *io.ReadWriteCloser

func NewStub(iceopts *IceOptions) *Stub {
	s := &Stub{stream: newStream("stub")}

	pc, err := newWebRTC(iceopts)
	if err != nil {
		log.Println("stub, rtc error:", err)
		return nil
	}

	pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
		log.Print("stub, pc ice state change:", state)
		if state == ice.ConnectionStateDisconnected {
			pc.Close()
		}
	})

	pc.OnDataChannel(func(dc *webrtc.RTCDataChannel) {
		log.Print("stub, OnDataChannel", dc)
		s.stream.Open(dc, func() {})
	})

	s.pc = pc
	return s
}

func (s *Stub) Open(sdp string, sender AnswerSender) error {
	log.Println("stub to set remote sdp:", sdp)
	err := s.pc.SetRemoteDescription(webrtc.RTCSessionDescription{
		Type: webrtc.RTCSdpTypeOffer,
		Sdp:  sdp,
	})
	if err != nil {
		log.Println("stub, set remote sdp failed!", err)
		return err
	}

	log.Println("proxy set remote sdp success! then to send answer to proxy ...")
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

	log.Println("proxy send answer success!")
	return nil
}

func (s *Stub) Stream() *Stream {
	return s.stream
}

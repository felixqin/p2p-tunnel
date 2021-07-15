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

// CreateAnswer() (string, error)
// ConnectProxy(sdp string) error
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
		s.stream.Open(dc)
	})

	s.pc = pc
	return s
}

func (s *Stub) CreateAnswer() (string, error) {
	answer, err := s.pc.CreateAnswer(nil)
	if err != nil {
		return "", err
	}

	return answer.Sdp, nil
}

func (s *Stub) ConnectProxy(sdp string) error {
	log.Println("stub to set remote sdp:", sdp)
	return s.pc.SetRemoteDescription(webrtc.RTCSessionDescription{
		Type: webrtc.RTCSdpTypeOffer,
		Sdp:  sdp,
	})
}

func (s *Stub) Stream() *Stream {
	return s.stream
}

// func (s *Stub) HandleOffer(fromClient string, offer *contacts.Offer) {
// 	go func() {
// 		log.Println("handler offer, sdp:", offer.Sdp)
// 		pc, err := newWebRTC(s.iceOptions)
// 		if err != nil {
// 			log.Println("rtc error:", err)
// 			return
// 		}

// 		sock, err := net.Dial("tcp", s.stubOptions.Addr)
// 		if err != nil {
// 			log.Println("sock dial filed:", err)
// 			pc.Close()
// 			return
// 		}

// 		pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
// 			log.Print("pc ice state change:", state)
// 			if state == ice.ConnectionStateDisconnected {
// 				pc.Close()
// 				sock.Close()
// 			}
// 		})

// 		pc.OnDataChannel(func(dc *webrtc.RTCDataChannel) {
// 			//dc.Lock()
// 			dc.OnOpen(func() {
// 				log.Print("dial:", s.stubOptions.Addr)
// 				// io.Copy(newWebRTCWriter(dc), sock)
// 				log.Println("disconnected")
// 			})

// 			dc.Onmessage(func(payload datachannel.Payload) {
// 				switch p := payload.(type) {
// 				case *datachannel.PayloadBinary:
// 					_, err := sock.Write(p.Data)
// 					if err != nil {
// 						log.Println("ssh write failed:", err)
// 						pc.Close()
// 						return
// 					}
// 				}
// 			})
// 			//dc.Unlock()
// 		})

// 		err = pc.SetRemoteDescription(webrtc.RTCSessionDescription{
// 			Type: webrtc.RTCSdpTypeOffer,
// 			Sdp:  offer.Sdp,
// 		})
// 		if err != nil {
// 			log.Println("rtc error:", err)
// 			pc.Close()
// 			sock.Close()
// 			return
// 		}

// 		answer, err := pc.CreateAnswer(nil)
// 		if err != nil {
// 			log.Println("rtc error:", err)
// 			pc.Close()
// 			sock.Close()
// 			return
// 		}

// 		err = contacts.SendAnswer(fromClient, &contacts.Answer{
// 			Sdp:  answer.Sdp,
// 			Stub: s.stubOptions.Name,
// 		})
// 		if err != nil {
// 			log.Println("rtc error:", err)
// 			pc.Close()
// 			sock.Close()
// 			return
// 		}
// 	}()
// }

package tunnel

import (
	"log"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

type Proxy struct {
	pc     *webrtc.RTCPeerConnection
	stream *Stream
}

// CreateOffer() (string, error)
// ConnectStub(sdp string) error
// Stream() *io.ReadWriteCloser

func NewProxy(iceopts *IceOptions) *Proxy {
	p := &Proxy{stream: newStream("proxy")}

	pc, err := newWebRTC(iceopts)
	if err != nil {
		log.Println("proxy, rtc error:", err)
		return nil
	}

	pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
		log.Print("proxy, pc ice state change:", state)
	})

	dc, err := pc.CreateDataChannel("data", nil)
	if err != nil {
		log.Println("proxy, create dc failed:", err)
		pc.Close()
		return nil
	}

	log.Print("proxy, DataChannel:", dc)
	p.stream.Open(dc)

	p.pc = pc
	return p
}

func (p *Proxy) CreateOffer() (string, error) {
	offer, err := p.pc.CreateOffer(nil)
	if err != nil {
		log.Println("create offer error:", err)
		return "", err
	}

	return offer.Sdp, nil
}

func (p *Proxy) ConnectStub(sdp string) error {
	log.Println("proxy to set remote sdp:", sdp)
	return p.pc.SetRemoteDescription(webrtc.RTCSessionDescription{
		Type: webrtc.RTCSdpTypeAnswer,
		Sdp:  sdp,
	})
}

func (p *Proxy) Stream() *Stream {
	return p.stream
}

// func (p *Proxy) HandleAnswer(fromClient string, answer *contacts.Answer) {
// 	if p.answerHandler != nil {
// 		p.answerHandler(fromClient, answer)
// 	}
// }

// func (p *Proxy) ListenAndServe() error {
// 	log.Println("listen", p.proxyOptions.Stub, "on", p.proxyOptions.Listen, "...")
// 	l, err := net.Listen("tcp", p.proxyOptions.Listen)
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		sock, err := l.Accept()
// 		if err != nil {
// 			log.Println("accept failed!", err)
// 			continue
// 		}

// 		go p.connectStub(sock)
// 	}
// }

// func (p *Proxy) connectStub(rw io.ReadWriter) {
// 	log.Println("connect stub ...")
// 	pc, err := newWebRTC(p.iceOptions)
// 	if err != nil {
// 		log.Println("rtc error:", err)
// 		return
// 	}

// 	pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
// 		log.Print("pc ice state change:", state)
// 	})

// 	dc, err := pc.CreateDataChannel("data", nil)
// 	if err != nil {
// 		log.Println("create dc failed:", err)
// 		pc.Close()
// 		return
// 	}

// 	//dc.Lock()
// 	dc.OnOpen(func() {
// 		io.Copy(newWebRTCWriter(dc), rw)
// 		pc.Close()
// 		log.Println("disconnected")
// 	})

// 	dc.OnMessage(func(payload datachannel.Payload) {
// 		switch p := payload.(type) {
// 		case *datachannel.PayloadBinary:
// 			_, err := rw.Write(p.Data)
// 			if err != nil {
// 				log.Println("write to stub failed:", err)
// 				pc.Close()
// 				return
// 			}
// 		}
// 	})
// 	//dc.Unlock()
// 	log.Print("DataChannel:", dc)

// 	// handle receive answer
// 	p.answerHandler = func(fromClient string, answer *contacts.Answer) {
// 		log.Println("handler answer, sdp:", answer.Sdp)
// 		err := pc.SetRemoteDescription(webrtc.RTCSessionDescription{
// 			Type: webrtc.RTCSdpTypeAnswer,
// 			Sdp:  answer.Sdp,
// 		})
// 		if err != nil {
// 			log.Println("set remote sdp failed!", err)
// 			pc.Close()
// 			return
// 		}
// 	}

// 	// offer, err := pc.CreateOffer(nil)
// 	// if err != nil {
// 	// 	log.Println("create offer error:", err)
// 	// 	pc.Close()
// 	// 	return
// 	// }

// 	// contact, err := contacts.FindContact(p.contact)
// 	// if err != nil {
// 	// 	log.Println("not found contact in contacts!", contact)
// 	// 	pc.Close()
// 	// 	return
// 	// }

// 	// err = contacts.SendOffer(contact.ClientId, &contacts.Offer{
// 	// 	Sdp:  offer.Sdp,
// 	// 	Stub: p.stub,
// 	// })
// 	// if err != nil {
// 	// 	log.Println("push error:", err)
// 	// 	pc.Close()
// 	// 	return
// 	// }
// }

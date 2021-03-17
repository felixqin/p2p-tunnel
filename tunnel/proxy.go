package tunnel

import (
	"io"
	"log"
	"net"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/pions/webrtc/pkg/ice"
)

type ProxyOptions struct {
	Listen  string `yaml:"listen"`
	Enable  bool   `yaml:"enable"`
	Contact string `yaml:"contact"`
	Stub    string `yaml:"stub"`
}

type Proxy struct {
	proxyOptions  *ProxyOptions
	iceOptions    *IceOptions
	answerHandler func(fromClient string, answer *contacts.Answer)
}

func NewProxy(opts *ProxyOptions, iceopts *IceOptions) *Proxy {
	return &Proxy{
		proxyOptions: opts,
		iceOptions:   iceopts,
	}
}

func (p *Proxy) HandleAnswer(fromClient string, answer *contacts.Answer) {
	if p.answerHandler != nil {
		p.answerHandler(fromClient, answer)
	}
}

func (p *Proxy) ListenAndServe() error {
	log.Println("listen", p.proxyOptions.Stub, "on", p.proxyOptions.Listen, "...")
	l, err := net.Listen("tcp", p.proxyOptions.Listen)
	if err != nil {
		return err
	}

	for {
		sock, err := l.Accept()
		if err != nil {
			log.Println("accept failed!", err)
			continue
		}

		go p.connectStub(sock)
	}
}

func (p *Proxy) connectStub(sock net.Conn) {
	log.Println("connect stub ...")
	pc, err := newWebRTC(p.iceOptions)
	if err != nil {
		log.Println("rtc error:", err)
		return
	}

	pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
		log.Print("pc ice state change:", state)
	})

	dc, err := pc.CreateDataChannel("data", nil)
	if err != nil {
		log.Println("create dc failed:", err)
		pc.Close()
		return
	}

	//dc.Lock()
	dc.OnOpen(func() {
		io.Copy(newWebRTCWriter(dc), sock)
		pc.Close()
		log.Println("disconnected")
	})

	dc.OnMessage(func(payload datachannel.Payload) {
		switch p := payload.(type) {
		case *datachannel.PayloadBinary:
			_, err := sock.Write(p.Data)
			if err != nil {
				log.Println("sock write failed:", err)
				pc.Close()
				return
			}
		}
	})
	//dc.Unlock()
	log.Print("DataChannel:", dc)

	// handle receive answer
	p.answerHandler = func(fromClient string, answer *contacts.Answer) {
		log.Println("handler answer, sdp:", answer.Sdp)
		err := pc.SetRemoteDescription(webrtc.RTCSessionDescription{
			Type: webrtc.RTCSdpTypeAnswer,
			Sdp:  answer.Sdp,
		})
		if err != nil {
			log.Println("set remote sdp failed!", err)
			pc.Close()
			return
		}
	}

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		log.Println("create offer error:", err)
		pc.Close()
		return
	}

	contact, err := contacts.FindContact(p.proxyOptions.Contact)
	if err != nil {
		log.Println("not found contact in contacts!", contact)
		pc.Close()
		return
	}

	err = contacts.SendOffer(contact.ClientId, &contacts.Offer{
		Sdp:  offer.Sdp,
		Stub: p.proxyOptions.Stub,
	})
	if err != nil {
		log.Println("push error:", err)
		pc.Close()
		return
	}
}

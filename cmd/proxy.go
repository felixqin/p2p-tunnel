package main

import (
	"io"
	"log"
	"net"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
	"github.com/pions/webrtc/pkg/ice"
)

type proxyConfigure struct {
	Listen string `yaml:"listen"`
	Stub   string `yaml:"stub"`
}

func proxyServe(proxyConf *proxyConfigure, iceConf *iceConfigure) error {
	log.Println("listen on", proxyConf.Listen, "...")
	l, err := net.Listen("tcp", proxyConf.Listen)
	if err != nil {
		return err
	}

	for {
		sock, err := l.Accept()
		if err != nil {
			log.Println("accept failed!", err)
			continue
		}

		go connectStub(proxyConf.Stub, iceConf, sock)
	}
}

func connectStub(stub string, iceConf *iceConfigure, sock net.Conn) {
	log.Println("connect stub ...")
	pc, err := newWebRTC(iceConf)
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
	contacts.HandleAnswerFunc(func(fromClient string, answer *contacts.Answer) {
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
	})

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		log.Println("create offer error:", err)
		pc.Close()
		return
	}

	contact, err := contacts.FindContact(stub)
	if err != nil {
		log.Println("not found stub in contacts!")
		pc.Close()
		return
	}

	err = contacts.SendOffer(contact.ClientId, &contacts.Offer{Sdp: offer.Sdp})
	if err != nil {
		log.Println("push error:", err)
		pc.Close()
		return
	}
}

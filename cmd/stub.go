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

type stubConfigure struct {
	Addr string `yaml:"addr"`
}

func stubServe(stubConf *stubConfigure, iceConf *iceConfigure) {
	// handle receive offer
	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		log.Println("handler offer, sdp:", offer.Sdp)
		pc, err := newWebRTC(iceConf)
		if err != nil {
			log.Println("rtc error:", err)
			return
		}

		sock, err := net.Dial("tcp", stubConf.Addr)
		if err != nil {
			log.Println("sock dial filed:", err)
			pc.Close()
			return
		}

		pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
			log.Print("pc ice state change:", state)
			if state == ice.ConnectionStateDisconnected {
				pc.Close()
				sock.Close()
			}
		})

		pc.OnDataChannel(func(dc *webrtc.RTCDataChannel) {
			//dc.Lock()
			dc.OnOpen(func() {
				log.Print("dial:", stubConf.Addr)
				io.Copy(newWebRTCWriter(dc), sock)
				log.Println("disconnected")
			})

			dc.Onmessage(func(payload datachannel.Payload) {
				switch p := payload.(type) {
				case *datachannel.PayloadBinary:
					_, err := sock.Write(p.Data)
					if err != nil {
						log.Println("ssh write failed:", err)
						pc.Close()
						return
					}
				}
			})
			//dc.Unlock()
		})

		err = pc.SetRemoteDescription(webrtc.RTCSessionDescription{
			Type: webrtc.RTCSdpTypeOffer,
			Sdp:  offer.Sdp,
		})
		if err != nil {
			log.Println("rtc error:", err)
			pc.Close()
			sock.Close()
			return
		}

		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			log.Println("rtc error:", err)
			pc.Close()
			sock.Close()
			return
		}

		err = contacts.SendAnswer(fromClient, &contacts.Answer{Sdp: answer.Sdp})
		if err != nil {
			log.Println("rtc error:", err)
			pc.Close()
			sock.Close()
			return
		}
	})
}

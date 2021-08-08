package tunnel

import (
	"log"
	"time"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/ice"
)

type Client struct {
	pc     *webrtc.RTCPeerConnection
	stream *Stream
}

type OfferSender func(sdp string, answerHandler func(sdp string)) error

// Open(ProxyMessager) error
// Stream() *io.ReadWriteCloser

func NewClient(iceopts *IceServers) *Client {
	p := &Client{stream: newStream("proxy")}

	pc, err := newWebRTC(iceopts)
	if err != nil {
		log.Println("proxy, rtc error:", err)
		return nil
	}

	pc.OnICEConnectionStateChange(func(state ice.ConnectionState) {
		log.Print("proxy, pc ice state change:", state)
	})

	p.pc = pc
	return p
}

func (p *Client) Open(sender OfferSender, onStreamOpen func(*Stream)) error {
	dc, err := p.pc.CreateDataChannel("data", nil)
	if err != nil {
		log.Println("proxy, create dc failed:", err)
		p.pc.Close()
		return err
	}

	log.Print("proxy, DataChannel:", dc)
	p.stream.Open(dc, func() { onStreamOpen(p.stream) })

	offer, err := p.pc.CreateOffer(nil)
	if err != nil {
		log.Println("proxy, create offer error:", err)
		p.stream.Close()
		p.pc.Close()
		return err
	}

	chAnswer := make(chan string)
	answerHandler := func(sdp string) {
		chAnswer <- sdp
	}

	go func() {
		for {
			log.Println("proxy, to send offer ...")
			err := sender(offer.Sdp, answerHandler)
			if err != nil {
				log.Println("proxy, send offer failed!", err)
			}

			select {
			case answer := <-chAnswer:
				log.Println("proxy to set remote sdp:", answer)
				err := p.pc.SetRemoteDescription(webrtc.RTCSessionDescription{
					Type: webrtc.RTCSdpTypeAnswer,
					Sdp:  answer,
				})
				if err != nil {
					log.Println("proxy, set remote sdp failed!", err)
					break
				}

				log.Println("proxy set remote sdp success!")
				return

			case <-time.After(10 * time.Second):
				log.Println("wait answer timeout!")
				break
			}
		}
	}()

	return nil
}

func (p *Client) Close() error {
	log.Println("proxy close")
	p.stream.Close()
	p.pc.Close()
	return nil
}

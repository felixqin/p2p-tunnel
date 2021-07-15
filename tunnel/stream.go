package tunnel

import (
	"fmt"
	"log"
	"time"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
)

type Stream struct {
	name string
	dc   *webrtc.RTCDataChannel
}

func newStream(name string) *Stream {
	return &Stream{name: name}
}

func (c *Stream) Open(dc *webrtc.RTCDataChannel) {
	dc.OnOpen(func() {
		log.Println(c.name, "dc OnOpen!")
		c.dc = dc
	})

	dc.OnMessage(func(payload datachannel.Payload) {
		// switch p := payload.(type) {
		// case *datachannel.PayloadBinary:
		// 	_, err := rw.Write(p.Data)
		// 	if err != nil {
		// 		log.Println("write to stub failed:", err)
		// 		pc.Close()
		// 		return
		// 	}
		// }
	})
}

func (c *Stream) Close() error {
	log.Println(c.name, "close")
	return nil
}

func (c *Stream) Read(b []byte) (int, error) {
	log.Println(c.name, "read data from dc")
	time.Sleep(5 * time.Minute)
	return 0, fmt.Errorf("%s, read not implemented", c.name)
}

func (c *Stream) Write(b []byte) (int, error) {
	log.Println(c.name, "write data to dc, len:", len(b))
	if c.dc == nil {
		return 0, fmt.Errorf("not open")
	}

	err := c.dc.Send(datachannel.PayloadBinary{Data: b})
	log.Println(c.name, "send error:", err)

	return len(b), err
}

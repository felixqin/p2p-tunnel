package tunnel

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/pions/webrtc"
	"github.com/pions/webrtc/pkg/datachannel"
)

type Stream struct {
	name     string
	dc       *webrtc.RTCDataChannel
	chaInbuf chan *bytes.Buffer
	curInbuf *bytes.Buffer
}

func newStream(name string) *Stream {
	inbufc := make(chan *bytes.Buffer, 200) // 接收队列
	return &Stream{name: name, chaInbuf: inbufc}
}

func (c *Stream) Open(dc *webrtc.RTCDataChannel, onOpen func()) {
	dc.OnOpen(func() {
		log.Println(c.name, "dc OnOpen!")
		c.dc = dc
		onOpen()
	})

	dc.OnMessage(func(payload datachannel.Payload) {
		switch p := payload.(type) {
		case *datachannel.PayloadBinary:
			log.Println(c.name, "on message, data len:", len(p.Data))
			c.chaInbuf <- bytes.NewBuffer(p.Data)
		}
	})
}

func (c *Stream) Close() error {
	log.Println(c.name, "stream close")
	return nil
}

func (c *Stream) Read(b []byte) (int, error) {
	if c.curInbuf == nil || c.curInbuf.Len() == 0 {
		select {
		case buf := <-c.chaInbuf:
			c.curInbuf = buf

		case <-time.After(30 * time.Second):
			return 0, fmt.Errorf("timeout")
		}
	}

	return c.curInbuf.Read(b)
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

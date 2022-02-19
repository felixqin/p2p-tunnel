package session

import (
	"io"
	"log"
	"net"

	mapset "github.com/deckarep/golang-set"
	"github.com/hashicorp/yamux"
)

type Proxy struct {
	stub   string
	client *Client
}

var proxies mapset.Set

func init() {
	proxies = mapset.NewSet()
}

func NewProxy(client *Client, stub string) *Proxy {
	p := &Proxy{stub: stub, client: client}
	proxies.Add(p)
	return p
}

func (p *Proxy) ListenAndServe(listenPort string) error {
	return proxyListenAndServe(listenPort, p.client.session, p.stub)
}

func (p *Proxy) Close() error {
	// TODO: stop listen and serve
	proxies.Remove(p)
	return nil
}

func proxyListenAndServe(listenPort string, session *yamux.Session, stub string) error {
	log.Println("listen tcp", listenPort, "...")
	l, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		log.Println("proxy, wait accept ...")
		c, err := l.Accept()
		if err != nil {
			log.Println("proxy, sock accept failed!", err)
			continue
		}
		// c := newLogReadWriteCloser("proxy_server", c1)

		go func() {
			log.Println("proxy, accepted!")
			defer c.Close()

			log.Println("proxy, to open yamux session ...")
			s, err := session.Open()
			if err != nil {
				log.Println("proxy, open yamux stream failed!", err)
				return
			}
			defer s.Close()

			// send stub name to tunnel server
			err = writeConnectRequest(s, stub)
			if err != nil {
				log.Println("proxy, write handshake request failed!", err)
				return
			}

			// check response
			code, err := readConnectResponse(s)
			if err != nil || code != 0 {
				log.Println("proxy, read handshake response failed!", err, "code", code)
				return
			}

			log.Println("proxy session opened! start io copy ...")
			go io.Copy(s, c)
			io.Copy(c, s)
			log.Println("proxy, to disconnect ...")
		}()
	}
}

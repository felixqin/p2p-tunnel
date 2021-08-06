package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
	"github.com/hashicorp/yamux"
)

func main() {
	log.Println("welcome to p2p tunnel")

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		terminateError := fmt.Errorf("%s", <-c)

		// Place whatever shutdown handling you want here

		errc <- terminateError
	}()

	// stubopts := make(map[string]*StubOptions)
	proxys := make(map[string]*tunnel.Proxy)
	answerHandlers := make(map[string]func(sdp string))
	// stubs := make(map[string]*tunnel.Stub)

	// Handle contacts message
	contacts.Open(configure.Contact)
	defer contacts.Close()

	contacts.HandleAnswerFunc(func(fromClient string, answer *contacts.Answer) {
		handler := answerHandlers[answer.Stub]
		if handler != nil {
			handler(answer.Sdp)
		}
	})

	// Create and start proxy services
	for _, opts := range configure.Proxys {
		var (
			port        = opts.Listen
			stub        = opts.Stub
			stubContact = opts.Contact
		)

		proxy := tunnel.NewProxy(configure.Ices)
		proxys[opts.Stub] = proxy

		offerSender := func(sdp string, answerHandler func(sdp string)) error {
			contact, err := contacts.FindContact(stubContact)
			if err != nil {
				return err
			}

			answerHandlers[stub] = answerHandler
			contacts.SendOffer(contact.ClientId, &contacts.Offer{
				Sdp:  sdp,
				Stub: stub,
			})

			return nil
		}

		err := proxy.Open(offerSender, func(stream *tunnel.Stream) {
			// 启动代理侦听服务和建立与Stub的连接
			go func() {
				log.Println("listen proxy on", port, "for stub", stub, "...")
				errc <- proxyListenAndServe(port, stream)
			}()
		})
		if err != nil {
			errc <- err
		}

		defer proxy.Close()
	}

	// Run!
	log.Println("exit:", <-errc)
}

func proxyListenAndServe(listenPort string, stream io.ReadWriteCloser) error {
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

		go func() {
			log.Println("proxy, accepted!")
			defer c.Close()

			log.Println("proxy, to create yamux client ...")
			session, err := yamux.Client(stream, nil)
			if err != nil {
				log.Println("proxy, create yamux client failed!", err)
				return
			}
			defer session.Close()

			log.Println("proxy, to open yamux session ...")
			s, err := session.Open()
			if err != nil {
				log.Println("proxy, open yamux stream failed!", err)
				return
			}
			defer s.Close()

			log.Println("proxy session opened!")
			io.Copy(s, c)
			go io.Copy(c, s)
			log.Println("proxy, to disconnect ...")
		}()
	}
}

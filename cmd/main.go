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

	proxys := make(map[string]*tunnel.Proxy)
	answerHandlers := make(map[string]func(sdp string))
	stubs := make(map[string]*tunnel.Stub)
	stubopts := make(map[string]*StubOptions)

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
			return contacts.SendOffer(contact.ClientId, &contacts.Offer{
				Sdp:  sdp,
				Stub: stub,
			})
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
			break
		}

		defer proxy.Close()
	}

	// Create stubs
	for _, opts := range configure.Stubs {
		stubopts[opts.Name] = opts
		stub := tunnel.NewStub(configure.Ices)
		stubs[opts.Name] = stub
	}

	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		var (
			stubContact = offer.Stub
			stub        = stubs[stubContact]
			upstream    = stubopts[stubContact].Upstream
		)

		if stub != nil {
			answerSender := func(sdp string) error {
				return contacts.SendAnswer(fromClient, &contacts.Answer{
					Sdp:  sdp,
					Stub: stubContact,
				})
			}

			err := stub.Open(offer.Sdp, answerSender, func(stream *tunnel.Stream) {
				// 启动Stub与upstream的连接
				go func() {
					log.Println("stub", stubContact, "serve and dail upstream", upstream, "...")
					errc <- stubDailAndServe(stream, upstream)
				}()
			})

			if err != nil {
				errc <- err
			}
		}
	})

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
		// c := newLogReadWriteCloser("proxy_server", c1)

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

			log.Println("proxy session opened! start io copy ...")
			go io.Copy(s, c)
			io.Copy(c, s)
			log.Println("proxy, to disconnect ...")
		}()
	}
}

func stubDailAndServe(stream io.ReadWriteCloser, upstream string) error {
	log.Println("stub, start yamux server ...")
	session, err := yamux.Server(stream, nil)
	if err != nil {
		log.Println("stub, create yamux server failed!", err)
		return err
	}
	defer session.Close()

	for {
		log.Println("stub, wait accept ...")
		s, err := session.Accept()
		if err != nil {
			log.Println("stub, yamux accept failed!", err)
			continue
		}

		go func() {
			log.Println("stub, accepted!")
			defer s.Close()

			log.Println("stub, dial upstream", upstream, "...")
			c, err := net.Dial("tcp", upstream)
			if err != nil {
				log.Println("stub, sock dial failed!", err)
				return
			}
			// c := newLogReadWriteCloser("stub_client", c1)
			defer c.Close()

			log.Println("stub upstream dailed! start io copy ...")
			go io.Copy(c, s)
			io.Copy(s, c)
		}()
	}
}

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

	stubopts := make(map[string]*StubOptions)
	proxys := make(map[string]*tunnel.Proxy)
	stubs := make(map[string]*tunnel.Stub)

	// Create and start proxy services
	for _, opts := range configure.Proxys {
		proxy := tunnel.NewProxy(configure.Ices)
		proxys[opts.Stub] = proxy

		// 启动代理侦听服务和建立与Stub的连接
		port := opts.Listen
		stub := opts.Stub
		contact := opts.Contact
		go func() {
			log.Println("listen proxy on", port, "for stub", stub, "...")
			errc <- proxyListenAndServe(proxy, port, stub, contact)
		}()
	}

	// Create stubs
	for _, opts := range configure.Stubs {
		stubopts[opts.Name] = opts
		stub := tunnel.NewStub(configure.Ices)
		stubs[opts.Name] = stub
	}

	// Handle contacts message
	contacts.Open(configure.Contact)
	defer contacts.Close()

	contacts.HandleAnswerFunc(func(fromClient string, answer *contacts.Answer) {
		proxy := proxys[answer.Stub]
		if proxy != nil {
			proxy.ConnectStub(answer.Sdp)
		}
	})

	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		stub := stubs[offer.Stub]
		opts := stubopts[offer.Stub]
		if stub != nil {
			stub.ConnectProxy(offer.Sdp)

			// 启动Stub与upstream的连接
			go func() {
				errc <- stubDailAndServe(stub, opts.Upstream, offer.Stub, fromClient)
			}()
		}
	})

	// Run!
	log.Println("exit:", <-errc)
}

func proxyListenAndServe(proxy *tunnel.Proxy, listenPort string, stub, stubContact string) error {
	l, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer l.Close()

	sdp, err := proxy.CreateOffer()
	if err != nil {
		return err
	}

	offerSent := false
	for {
		log.Println("proxy, wait accept ...")
		c, err := l.Accept()
		if err != nil {
			log.Println("proxy, sock accept failed!", err)
			continue
		}

		// 如果没有连接成功，发送offer给stub
		if !offerSent {
			log.Println("proxy, to send offer ...")
			contact, err := contacts.FindContact(stubContact)
			if err != nil {
				log.Println("proxy, not found contact in contacts!", contact)
				c.Close()
				continue
			}

			err = contacts.SendOffer(contact.ClientId, &contacts.Offer{
				Sdp:  sdp,
				Stub: stub,
			})
			if err != nil {
				log.Println("proxy, send offer failed!", err)
				c.Close()
				continue
			}

			// FIXME: 判断接收到answer才停止发送offer
			offerSent = true
		}

		go func() {
			log.Println("proxy, accepted!")
			defer c.Close()

			log.Println("proxy, to create yamux client ...")
			session, err := yamux.Client(proxy.Stream(), nil)
			if err != nil {
				log.Println("proxy, create yamux client failed!", err)
				return
			}
			defer session.Close()

			log.Println("proxy, to open yamux session ...")
			stream, err := session.Open()
			if err != nil {
				log.Println("proxy, open yamux stream failed!", err)
				return
			}
			defer stream.Close()

			log.Println("proxy session opened!")
			io.Copy(stream, c)
			go io.Copy(c, stream)
			log.Println("proxy, to disconnect ...")
		}()
	}
}

func stubDailAndServe(stub *tunnel.Stub, upstream string, stubName, toClient string) error {
	sdp, err := stub.CreateAnswer()
	if err != nil {
		return err
	}

	err = contacts.SendAnswer(toClient, &contacts.Answer{
		Sdp:  sdp,
		Stub: stubName,
	})
	if err != nil {
		log.Println("stub send answer failed!", err)
		return err
	}

	log.Println("stub, start yamux server ...")
	session, err := yamux.Server(stub.Stream(), nil)
	if err != nil {
		log.Println("stub, create yamux server failed!", err)
		return err
	}
	defer session.Close()

	for {
		log.Println("stub, wait accept ...")
		stream, err := session.Accept()
		if err != nil {
			log.Println("stub, yamux accept failed!", err)
			continue
		}

		go func() {
			log.Println("stub, accepted!")
			defer stream.Close()

			log.Println("stub, dial upstream ...")
			c, err := net.Dial("tcp", upstream)
			if err != nil {
				log.Println("stub, sock dial failed!", err)
				return
			}
			defer c.Close()

			log.Println("stub upstream dailed!")
			io.Copy(c, stream)
			go io.Copy(stream, c)
		}()
	}
}

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

	// Handle contacts message
	contacts.Open(configure.Contact)
	defer contacts.Close()

	// create and start tunnel server
	tunnelServer := tunnel.NewServer(configure.Ices)
	defer tunnelServer.Close()
	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		answerSender := func(sdp string) error {
			return contacts.SendAnswer(fromClient, &contacts.Answer{
				Sdp: sdp,
			})
		}

		err := tunnelServer.Open(offer.Sdp, answerSender, func(stream *tunnel.Stream) {
			go func() {
				log.Println("stub serve ...")
				errc <- stubServe(stream)
			}()
		})

		if err != nil {
			errc <- err
		}
	})

	// create clients for all proxy contacts
	contacts.HandleAnswerFunc(handleAnswer)
	tunnelClients := make(map[string]*tunnel.Client)
	for _, clientOpts := range configure.Proxys {
		if !clientOpts.Enable {
			continue
		}

		serverContactName := clientOpts.Contact
		if tunnelClients[serverContactName] != nil {
			// tunnel client 创建过无需再创建，直接跳过
			continue
		}

		// create and start tunnel client
		tunnelClient := tunnel.NewClient(configure.Ices)
		defer tunnelClient.Close()
		tunnelClients[serverContactName] = tunnelClient

		offerSender := makeOfferSender(serverContactName)
		err := tunnelClient.Open(offerSender, func(stream *tunnel.Stream) {
			for _, proxyOpts := range configure.Proxys {
				if !proxyOpts.Enable || proxyOpts.Contact != serverContactName {
					// 跳过不通过此contact连接的代理服务
					continue
				}

				port := proxyOpts.Listen
				stub := proxyOpts.Stub
				go func() {
					log.Println("listen proxy on", port, "for stub", stub, "...")
					errc <- proxyListenAndServe(port, stream, stub)
				}()
			}
		})
		if err != nil {
			errc <- err
			break
		}
	}

	// Run!
	log.Println("exit:", <-errc)
}

func proxyListenAndServe(listenPort string, stream io.ReadWriteCloser, stub string) error {
	l, err := net.Listen("tcp", listenPort)
	if err != nil {
		return err
	}
	defer l.Close()

	log.Println("proxy, to create yamux client ...")
	session, err := yamux.Client(stream, nil)
	if err != nil {
		log.Println("proxy, create yamux client failed!", err)
		return err
	}
	defer session.Close()

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

func stubServe(stream io.ReadWriteCloser) error {
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
			// time.Sleep(5 * time.Second)
			continue
		}

		go func() {
			log.Println("stub, accepted!")
			defer s.Close()

			stub, err := readConnectRequest(s)
			if err != nil {
				log.Println("stub, read handshake request failed!", err)
				return
			}

			opt, err := findStubOption(stub)
			if err != nil {
				log.Println("stub, find stub", stub, "option failed!", err)
				return
			}

			log.Println("stub, dial upstream", opt.Upstream, "...")
			c, err := net.Dial("tcp", opt.Upstream)
			if err != nil {
				log.Println("stub, sock dial failed!", err)
				writeConnectResponse(s, -1)
				return
			}
			// c := newLogReadWriteCloser("stub_client", c1)
			defer c.Close()

			err = writeConnectResponse(s, 0)
			if err != nil {
				log.Println("stub, write handshake response failed!", err)
				return
			}

			log.Println("stub upstream dailed! start io copy ...")
			go io.Copy(c, s)
			io.Copy(s, c)
		}()
	}
}

func findStubOption(stub string) (*StubOption, error) {
	for _, opt := range configure.Stubs {
		if opt.Name == stub {
			return opt, nil
		}
	}

	return nil, fmt.Errorf("stub option not found")
}

package session

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

type StubOption struct {
	Name     string
	Upstream string
}

var stubOptions []*StubOption

func EnableStub(name, upstream string) error {
	stubOptions = append(stubOptions, &StubOption{
		Name:     name,
		Upstream: upstream,
	})

	return nil
}

func DisableStub(name string) error {
	return nil
}

func DumpStubs() {
	for _, opt := range stubOptions {
		fmt.Printf("%-10v%v\n", opt.Name, opt.Upstream)
	}
}

func findStubOption(stub string) (*StubOption, error) {
	for _, opt := range stubOptions {
		if opt.Name == stub {
			return opt, nil
		}
	}

	return nil, fmt.Errorf("stub option not found")
}

func stubServe(session *yamux.Session) error {
	for {
		log.Println("stub, wait accept ...")
		s, err := session.Accept()
		if err != nil {
			log.Println("stub, yamux accept failed!", err)
			return err
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

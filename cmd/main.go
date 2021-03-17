package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
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

	stubs := make(map[string]*tunnel.Stub)
	proxys := make(map[string]*tunnel.Proxy)

	// Create stubs
	for _, opts := range configure.Stubs {
		stub := tunnel.NewStub(opts, configure.Ices)
		stubs[opts.Name] = stub
	}

	// Create proxys
	for _, opts := range configure.Proxys {
		proxy := tunnel.NewProxy(opts, configure.Ices)
		proxys[opts.Stub] = proxy
	}

	// Handle contacts message
	contacts.Open(configure.Contact)
	defer contacts.Close()

	contacts.HandleOfferFunc(func(fromClient string, offer *contacts.Offer) {
		stub := stubs[offer.Stub]
		if stub != nil {
			stub.HandleOffer(fromClient, offer)
		}
	})

	contacts.HandleAnswerFunc(func(fromClient string, answer *contacts.Answer) {
		proxy := proxys[answer.Stub]
		if proxy != nil {
			proxy.HandleAnswer(fromClient, answer)
		}
	})

	// Start proxys service
	for _, proxy := range proxys {
		p := proxy
		go func() {
			errc <- p.ListenAndServe()
		}()
	}

	// Run!
	log.Println("exit:", <-errc)
}

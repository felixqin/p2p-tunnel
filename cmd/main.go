package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/felixqin/p2p-tunnel/contacts"
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

	contacts.Open(configure.Contacts)
	defer contacts.Close()

	if configure.Proxy != nil {
		go func() {
			errc <- proxyServe(configure.Proxy, configure.Ice)
		}()
	}

	if configure.Stub != nil {
		go func() {
			errc <- stubServe(configure.Stub)
		}()
	}

	// Run!
	log.Println("exit:", <-errc)
}

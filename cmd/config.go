package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
	"gopkg.in/yaml.v2"
)

var configure struct {
	Contact *contacts.Options      `yaml:"contact"`
	Ices    *tunnel.IceOptions     `yaml:"ices"`
	Proxys  []*tunnel.ProxyOptions `yaml:"proxys"`
	Stubs   []*tunnel.StubOptions  `yaml:"stubs"`
}

func init() {
	conf := flag.String("conf", "/etc/p2p-tunnel/config.proxy.yaml", "configure file")
	flag.Parse()

	txt, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatalln("read configure file failed!", err)
	}

	err = yaml.Unmarshal(txt, &configure)
	if err != nil {
		log.Fatalln("unmarshal config failed!", err)
	}
}

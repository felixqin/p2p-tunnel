package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/felixqin/p2p-tunnel/contacts"
	"gopkg.in/yaml.v2"
)

var configure struct {
	Contacts contacts.Configure `yaml:"contacts`
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

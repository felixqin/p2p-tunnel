package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/felixqin/p2p-tunnel/contacts"
	"github.com/felixqin/p2p-tunnel/tunnel"
	"gopkg.in/yaml.v2"
)

type ProxyOption struct {
	Listen  string `yaml:"listen"`
	Enable  bool   `yaml:"enable"`
	Stub    string `yaml:"stub"`    // 要连接到的对方 Stub 名称
	Contact string `yaml:"contact"` // Stub 所在在联系人名称
}

type StubOption struct {
	Name     string `yaml:"name"`
	Enable   bool   `yaml:"enable"`
	Upstream string `yaml:"upstream"`
}

var configure struct {
	Contact *contacts.Option   `yaml:"contact"`
	Ices    *tunnel.IceServers `yaml:"ices"`
	Proxys  []*ProxyOption     `yaml:"proxys"`
	Stubs   []*StubOption      `yaml:"stubs"`
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

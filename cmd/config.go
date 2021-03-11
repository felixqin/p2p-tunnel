package main

// import (
// 	"flag"
// 	"io/ioutil"
// 	"log"

// 	"gopkg.in/yaml.v2"
// )

// var configure struct {
// 	GBS struct {
// 		Listen struct {
// 			Debug   string
// 			Http    string
// 			Admin   string
// 			Private string
// 			Udp     string
// 			Tcp     string
// 		}
// 	} `yaml:"gbs`
// }

// var (
// 	Listen = &configure.GBS.Listen
// )

// func init() {
// 	conf := flag.String("conf", "/etc/gbs/config.yaml", "configure file")
// 	flag.Parse()

// 	txt, err := ioutil.ReadFile(*conf)
// 	if err != nil {
// 		log.Fatalln("read configure file failed!", err)
// 	}

// 	err = yaml.Unmarshal(txt, &configure)
// 	if err != nil {
// 		log.Fatalln("unmarshal config failed!", err)
// 	}
// }

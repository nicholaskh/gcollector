package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nicholaskh/etclib"
	"github.com/nicholaskh/golib/ip"
	log "github.com/nicholaskh/log4go"
)

var localAddr string

func init() {
	loadLocalAddr()
}

func RegisterEtc(etcServers []string) error {
	err := etclib.Dial(etcServers)
	if err != nil {
		return err
	}
	err = etclib.BootService(localAddr, etclib.SERVICE_GCOLLECTOR)
	return err
}

func UnregisterEtc() error {
	return etclib.ShutdownService(localAddr, etclib.SERVICE_GCOLLECTOR)
}

func GetPiped() (string, error) {
	addrs, err := etclib.Children(fmt.Sprintf("/%s", etclib.SERVICE_PIPED))
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		panic("No piped found")
	}
	log.Info(addrs)
	time.Sleep(time.Second)
	return addrs[rand.Intn(len(addrs))], nil
}

func loadLocalAddr() {
	localIps := ip.LocalIpv4Addrs()
	if len(localIps) == 0 {
		panic("No local ip address found")
	}
	localAddr = localIps[0]
}

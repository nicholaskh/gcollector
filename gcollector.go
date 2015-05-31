package main

import (
	"github.com/nicholaskh/golib/server"
)

type Gcollector struct {
	config *GcollectorConfig
}

func NewGcollector(config *GcollectorConfig) *Gcollector {
	this := new(Gcollector)
	this.config = config

	return this
}

func (this *Gcollector) RunForever() {
	go server.StartPingServer(this.config.UdpPort)

	forwarder := NewForwarder(this.config.Forwarder)
	for _, inputConfig := range this.config.Inputs {
		poller := NewPoller(inputConfig, forwarder)
		go poller.Poll()
	}

	forwarder.Send()
}

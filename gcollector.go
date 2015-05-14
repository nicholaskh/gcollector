package main

import (
//	log "github.com/nicholaskh/log4go"
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
	forwarder := NewForwarder(this.config.ToAddr)
	for _, inputConfig := range this.config.Inputs {
		poller := NewPoller(inputConfig, forwarder)
		go poller.Poll()
	}

	forwarder.Send()
}

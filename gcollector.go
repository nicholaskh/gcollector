package main

//	log "github.com/nicholaskh/log4go"

type Gcollector struct {
	config *GcollectorConfig
}

func NewGcollector(config *GcollectorConfig) *Gcollector {
	this := new(Gcollector)
	this.config = config

	return this
}

func (this *Gcollector) RunForever() {
	startUdpServer(this.config.UdpPort)

	forwarder := NewForwarder(this.config.Forwarder)
	for _, inputConfig := range this.config.Inputs {
		poller := NewPoller(inputConfig, forwarder)
		go poller.Poll()
	}

	forwarder.Send()
}

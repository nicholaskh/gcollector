package main

import (
	"net"

	log "github.com/nicholaskh/log4go"
)

type Forwarder struct {
	queue chan string
	net.Conn
}

func NewForwarder(toAddr string) *Forwarder {
	this := new(Forwarder)
	var err error
	this.Conn, err = net.Dial("tcp", toAddr)
	if err != nil {
		log.Error(err)
	}
	this.queue = make(chan string, 100000)

	return this
}

func (this *Forwarder) Enqueue(line string) {
	this.queue <- line
}

func (this *Forwarder) Send() {
	for line := range this.queue {
		log.Debug(line)
		this.Write([]byte(line))
	}
}

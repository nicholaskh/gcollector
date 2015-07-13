package main

import (
	"net"

	log "github.com/nicholaskh/log4go"
)

type Forwarder struct {
	config *ForwarderConfig
	queue  chan string
	proto  *Protocol
	net.Conn
}

func NewForwarder(config *ForwarderConfig, app string) *Forwarder {
	this := new(Forwarder)
	this.config = config
	this.proto = NewProtocol(app)
	this.reconnect()
	this.queue = make(chan string, this.config.Backlog)

	return this
}

func (this *Forwarder) reconnect() {
	if this.Conn != nil {
		this.Conn.Close()
	}
	conn, err := net.Dial("tcp", this.config.ToAddr)
	if err != nil {
		log.Error(err)
	}
	this.Conn = conn

	if conn != nil {
		go func() {
			for {
				_, err := this.Conn.Read(make([]byte, 1000))
				if err != nil {
					log.Warn(err)
				}
				this.Conn.Close()
				break
			}
		}()
	}
}

func (this *Forwarder) Enqueue(line string) {
	log.Info(line)
	this.queue <- line
}

func (this *Forwarder) Send() {
	for line := range this.queue {
		data := this.proto.Marshal([]byte(line))
		if this.Conn == nil {
			this.reconnect()
		}
		if this.Conn != nil {
			_, err := this.Write(data)
			if err != nil {
				log.Error("write error: %s", err.Error())
				// retry
				for i := 0; i < 3; i++ {
					this.reconnect()
					if this.Conn != nil {
						_, err = this.Write(data)
						if err == nil {
							break
						}
					}
				}
			}
		}
	}
}

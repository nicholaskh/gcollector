package main

import (
	"net"

	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
)

type Forwarder struct {
	config *ForwarderConfig
	queue  chan string
	proto  *server.Protocol
}

func NewForwarder(config *ForwarderConfig) *Forwarder {
	this := new(Forwarder)
	this.config = config
	this.proto = server.NewProtocol()
	this.reconnect()
	this.queue = make(chan string, this.config.Backlog)

	return this
}

func (this *Forwarder) reconnect() {
	if this.proto.Conn != nil {
		this.proto.Conn.Close()
	}
	conn, err := net.Dial("tcp", this.config.ToAddr)
	if err != nil {
		log.Error(err)
	}
	this.proto.SetConn(conn)

	if conn != nil {
		go func() {
			for {
				_, err := this.proto.Conn.Read(make([]byte, 1000))
				if err != nil {
					log.Warn(err)
				}
				this.proto.Conn.Close()
				break
			}
		}()
	}
}

func (this *Forwarder) Enqueue(line string) {
	this.queue <- line
}

func (this *Forwarder) Send() {
	for line := range this.queue {
		log.Debug(line)
		if this.proto.Conn == nil {
			this.reconnect()
		}
		if this.proto.Conn != nil {
			_, err := this.proto.Write([]byte(line))
			if err != nil {
				log.Error("write error: %s", err.Error())
				// retry
				for i := 0; i < 3; i++ {
					this.reconnect()
					if this.proto.Conn != nil {
						_, err = this.proto.Write([]byte(line))
						if err == nil {
							break
						}
					}
				}
			}
		}
	}
}

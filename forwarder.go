package main

import (
	"net"

	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
)

type Forwarder struct {
	queue  chan string
	toAddr string
	proto  *server.Protocol
}

func NewForwarder(toAddr string) *Forwarder {
	this := new(Forwarder)
	this.toAddr = toAddr
	this.proto = server.NewProtocol()
	this.reconnect()
	this.queue = make(chan string, 100000)

	return this
}

func (this *Forwarder) reconnect() {
	if this.proto.Conn != nil {
		this.proto.Conn.Close()
	}
	conn, err := net.Dial("tcp", this.toAddr)
	if err != nil {
		log.Error(err)
	}
	this.proto.SetConn(conn)

	if conn != nil {
		go func() {
			for {
				_, err := this.proto.Conn.Read(make([]byte, 1000))
				if err != nil {
					this.proto.Conn.Close()
					break
				}
			}
		}()
	}
}

func (this *Forwarder) Enqueue(line string) {
	this.queue <- line
}

func (this *Forwarder) Send() {
	i := 0
	for line := range this.queue {
		log.Debug(line)
		i++
		if i%10000 == 0 {
			log.Info("sent %d lines", i)
		}
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

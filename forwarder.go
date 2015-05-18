package main

import (
	"net"

	log "github.com/nicholaskh/log4go"
)

type Forwarder struct {
	queue chan string
	net.Conn
	toAddr string
}

func NewForwarder(toAddr string) *Forwarder {
	this := new(Forwarder)
	this.toAddr = toAddr
	this.reconnect()
	this.queue = make(chan string, 100000)

	return this
}

func (this *Forwarder) reconnect() {
	var err error
	this.Conn, err = net.Dial("tcp", this.toAddr)
	if err != nil {
		log.Error(err)
	}

	if this.Conn != nil {
		go func() {
			for {
				_, err := this.Conn.Read(make([]byte, 1000))
				if err != nil {
					this.Conn.Close()
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
	for line := range this.queue {
		log.Debug(line)
		if this.Conn == nil {
			this.reconnect()
		}
		if this.Conn != nil {
			_, err := this.Write([]byte(line))
			if err != nil {
				log.Error("write error: %s", err.Error())
				// retry
				for i := 0; i < 3; i++ {
					this.reconnect()
					if this.Conn != nil {
						_, err = this.Write([]byte(line))
						if err == nil {
							break
						}
					}
				}
			}
		}
	}
}

package main

import (
	"net"

	log "github.com/nicholaskh/log4go"
)

const (
	PACKET_PING = "PING"
	PACKET_PONG = "PONG"
)

func startUdpServer(port int) {
	listenConn, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	if err != nil {
		panic(err)
	}
	for {
		data := make([]byte, 1000)
		n, remoteAddr, err := listenConn.ReadFromUDP(data)
		if err != nil {
			log.Error("Read from %s error", remoteAddr)
			continue
		}
		rev := string(data[:n])
		if rev == PACKET_PING {
			_, err := listenConn.WriteTo([]byte(PACKET_PONG), remoteAddr)
			if err != nil {
				log.Error("Write to %s error: %s", remoteAddr, err.Error())
			}
		} else {
			log.Warn("Received unknown packet: %s", rev)
			continue
		}
	}
}

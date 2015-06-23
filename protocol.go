package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	log "github.com/nicholaskh/log4go"
)

const (
	HEAD_LENGTH = 8
)

type Protocol struct {
	net.Conn
	app string
}

func NewProtocol(app string) *Protocol {
	this := new(Protocol)
	this.app = app
	log.Info(this.app)

	return this
}

func (this *Protocol) SetConn(conn net.Conn) {
	this.Conn = conn
}

//len+appLength+app+payload
func (this *Protocol) Marshal(payload []byte) []byte {
	buf := bytes.NewBuffer([]byte{})
	tmpBuff := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, int32(len(payload)))
	binary.Write(tmpBuff, binary.BigEndian, int32(len(this.app)))
	buf.Write(tmpBuff.Bytes())
	buf.Write([]byte(this.app))
	buf.Write(payload)
	log.Info(string(buf.Bytes()))
	return buf.Bytes()
}

func (this *Protocol) Read() ([]byte, []byte, error) {
	buf := make([]byte, HEAD_LENGTH)
	err := this.ReadN(this.Conn, buf, HEAD_LENGTH)
	if err != nil {
		log.Error("[Protocol] Read data length error: %s", err.Error())
		return []byte{}, []byte{}, err
	}
	//data length
	b_buf := bytes.NewBuffer(buf[:4])
	var dataLength int32
	binary.Read(b_buf, binary.BigEndian, &dataLength)
	log.Info("dataLength: %d", dataLength)
	//app length
	var appLength int32
	b_buf.Write(buf[4:8])
	binary.Read(b_buf, binary.BigEndian, &appLength)
	log.Info("appLength: %d", appLength)

	//app + data
	payloadLength := int(dataLength + appLength)
	payload := make([]byte, payloadLength)
	err = this.ReadN(this.Conn, payload, payloadLength)
	if err != nil && err != io.EOF {
		log.Error("[Protocol] Read data error: %s", err.Error())
		return []byte{}, []byte{}, err
	}

	return payload[:appLength], payload[appLength:payloadLength], nil
}

func (this *Protocol) ReadN(conn net.Conn, buf []byte, n int) error {
	buffer := bytes.NewBuffer([]byte{})
	for n > 0 {
		b_buf := make([]byte, n)
		readN, err := conn.Read(b_buf)
		if err != nil {
			return err
		}
		n -= readN
		buffer.Write(b_buf)
	}
	copy(buf, buffer.Bytes())
	return nil
}

package main

import (
	"fmt"
	"os"

	log "github.com/nicholaskh/log4go"
	"github.com/nicholaskh/tail"
)

type Poller struct {
	config    *InputConfig
	forwarder *Forwarder
	parser    *Parser
}

func NewPoller(config *InputConfig, forwarder *Forwarder) *Poller {
	this := new(Poller)
	this.config = config
	this.forwarder = forwarder
	this.parser = NewParser()

	return this
}

func (this *Poller) Poll() {
	log.Info("Tail file: %s", this.config.File)
	t, err := tail.TailFile(this.config.File, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Error("tail file[%s] errer: %s", this.config.File, err.Error())
	}
	for line := range t.Lines {
		txt := line.Text
		if tag := this.filter(txt); tag != "" {
			this.forwarder.Enqueue(fmt.Sprintf("%s|%s", tag, txt))
		}
	}
}

func (this *Poller) filter(txt string) (tag string) {
	for _, tp := range this.config.Types {
		switch tp {
		case LOG_TYPE_NGINX_500:
			logPart := this.parser.parse(txt, tp)
			if logPart[5] == "500" {
				return LOG_TYPE_NGINX_500
			}
		case LOG_TYPE_NGINX_404:
			logPart := this.parser.parse(txt, tp)
			if logPart[5] == "404" {
				return LOG_TYPE_NGINX_404
			}
		case LOG_TYPE_PHP_ERROR:
			if this.parser.match(txt, tp) {
				return LOG_TYPE_PHP_ERROR
			}
		}
	}
	return ""
}

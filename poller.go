package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

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
	if strings.Contains(this.config.File, "*") {
		last_sep := strings.LastIndex(this.config.File, PATH_SEP)
		dir := this.config.File[0:last_sep]
		dir_list, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Error("read dir error: %s", err.Error())
			return
		}

		filename := this.config.File[last_sep+1:]
		sourceReg := fmt.Sprintf("^%s$", strings.Replace(filename, "*", ".*?", -1))
		reg := regexp.MustCompile(sourceReg)

		for _, path := range dir_list {
			if path.IsDir() == true {
				continue
			}

			if matchFile(reg, path.Name()) {
				go this.tailFile(fmt.Sprintf("%s%s%s", dir, PATH_SEP, path.Name()))
			}

		}

	} else {
		go this.tailFile(this.config.File)
	}
}

func matchFile(sourceReg *regexp.Regexp, destFile string) bool {
	return sourceReg.MatchString(destFile)
}

func (this *Poller) tailFile(filename string) {
	log.Info("Tail file: %s", filename)
	t, err := tail.TailFile(filename, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Error("tail file[%s] errer: %s", filename, err.Error())
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
		case LOG_TYPE_NGINX_500, LOG_TYPE_APACHE_500:
			logPart := this.parser.parse(txt, tp)
			if logPart[5] == "500" {
				return tp
			}
		case LOG_TYPE_NGINX_404, LOG_TYPE_APACHE_404:
			logPart := this.parser.parse(txt, tp)
			if logPart[5] == "404" {
				return tp
			}
		case LOG_TYPE_PHP_ERROR:
			if this.parser.match(txt, tp) {
				return tp
			}
		case LOG_TYPE_ELAPSED:
			return tp

		case LOG_TYPE_ANY:
			return tp
		}
	}
	return ""
}

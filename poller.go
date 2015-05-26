package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/howeyc/fsnotify"
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
	if this.config.File[len(this.config.File)-2:] == "**" {
		last_sep := strings.LastIndex(this.config.File, PATH_SEP)
		dir := this.config.File[0:last_sep]
		this.tailFilesInDirRecursive(dir)
	} else if strings.Contains(this.config.File, "**") {
		panic("'**' pattern must in the end")
	} else if strings.Contains(this.config.File, "*") {
		last_sep := strings.LastIndex(this.config.File, PATH_SEP)
		dir := this.config.File[0:last_sep]
		filename := this.config.File[last_sep+1:]
		this.tailFilesInDir(dir, filename)
	} else {
		go this.tailFile(this.config.File, false)
	}
}

func matchFile(sourceReg *regexp.Regexp, destFile string) bool {
	return sourceReg.MatchString(destFile)
}

func (this *Poller) tailFile(filename string, isNew bool) {
	log.Info("Tail file: %s", filename)
	var location *tail.SeekInfo
	if isNew {
		location = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_SET}
	} else {
		location = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	}
	t, err := tail.TailFile(filename, tail.Config{Follow: true, Location: location})
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
			if len(logPart) < 6 {
				return ""
			}
			if logPart[6] == "500" {
				return tp
			}
		case LOG_TYPE_NGINX_404, LOG_TYPE_APACHE_404:
			logPart := this.parser.parse(txt, tp)
			if len(logPart) < 6 {
				return ""
			}
			if logPart[6] == "404" {
				return tp
			}
		case LOG_TYPE_PHP_ERROR:
			if this.parser.match(txt, tp) {
				return tp
			}
		case LOG_TYPE_APP:
			return tp
		}
	}
	return ""
}

func (this *Poller) tailFilesInDir(dir string, filenameReg string) {
	sourceReg := fmt.Sprintf("^%s$", strings.Replace(filenameReg, "*", ".*?", -1))
	dir_list, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("read dir error: %s", err.Error())
		return
	}

	reg := regexp.MustCompile(sourceReg)

	for _, path := range dir_list {
		if path.IsDir() == true {
			continue
		}

		if matchFile(reg, path.Name()) {
			go this.tailFile(fmt.Sprintf("%s%s%s", dir, PATH_SEP, path.Name()), false)
		}

	}

	go this.watchDir(dir, false)
}

func (this *Poller) watchDir(dir string, followDir bool) {
	watcher, err := fsnotify.NewWatcher()
	err = watcher.Watch(dir)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := watcher.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Process events
	for {
		select {
		case ev := <-watcher.Event:
			if ev.IsCreate() {
				finfo, err := os.Stat(ev.Name)
				if err != nil {
					panic(err)
					return
				}
				if finfo.IsDir() {
					if followDir {
						this.tailFilesInDir(ev.Name, "*")
					} else {
						continue
					}
				} else {
					go this.tailFile(ev.Name, true)
				}
			}
		// TODO -- when delete file, stop tail
		case err := <-watcher.Error:
			log.Error("error:", err)
		}
	}

}

func (this *Poller) tailFilesInDirRecursive(dir string) {
	dir_list, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("read dir error: %s", err.Error())
		return
	}

	for _, path := range dir_list {
		fullPath := fmt.Sprintf("%s%s%s", dir, PATH_SEP, path.Name())
		if path.IsDir() == true {
			this.tailFilesInDirRecursive(fullPath)
		} else {
			go this.tailFile(fullPath, false)
		}
	}

	go this.watchDir(dir, true)
}

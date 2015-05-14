package main

import (
	"regexp"

	log "github.com/nicholaskh/log4go"
)

var (
	nginxLogReg = regexp.MustCompile(`^([^ ]*) ([^ ]*) ([^ ]*) \[([^\]]*)\] "([^"]*)" ([^ ]*) ([^ ]*) "([^"]*)" "([^"]*)"$`)
	parser      *Parser
)

type Parser struct {
	Regexps map[string]*regexp.Regexp
}

func NewParser() *Parser {
	this := new(Parser)
	this.Regexps = make(map[string]*regexp.Regexp)
	this.Regexps[LOG_TYPE_NGINX_404] = nginxLogReg
	this.Regexps[LOG_TYPE_NGINX_500] = nginxLogReg

	return this
}

func (this *Parser) parse(txt string, tp string) []string {
	re, exists := this.Regexps[tp]
	if !exists {
		log.Warn("regexp not found for %s", tp)
	}
	return re.FindAllStringSubmatch(txt, -1)[0][1:]
}

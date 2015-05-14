package main

import (
	"fmt"
	"strings"

	conf "github.com/nicholaskh/jsconf"
)

const (
	LOG_TYPE_NGINX_500 = "nginx_500"
	LOG_TYPE_NGINX_404 = "nginx_404"
)

type GcollectorConfig struct {
	ToAddr string

	Inputs []*InputConfig
}

func (this *GcollectorConfig) LoadConfig(cf *conf.Conf) {
	this.ToAddr = cf.String("to_addr", ":5687")

	this.Inputs = make([]*InputConfig, 0)
	for i, _ := range cf.List("inputs", nil) {
		section, err := cf.Section(fmt.Sprintf("inputs[%d]", i))
		if err != nil {
			panic(err)
		}
		input := new(InputConfig)
		input.File = section.String("file", "")
		types := strings.Split(section.String("types", ""), ",")
		for _, tp := range types {
			input.Types = append(input.Types, strings.Trim(tp, " "))
		}
		this.Inputs = append(this.Inputs, input)
	}
}

type InputConfig struct {
	File  string
	Types []string
}

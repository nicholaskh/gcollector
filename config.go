package main

import (
	"fmt"
	"strings"

	conf "github.com/nicholaskh/jsconf"
)

type GcollectorConfig struct {
	UdpPort int

	Forwarder *ForwarderConfig

	Inputs []*InputConfig
}

func (this *GcollectorConfig) LoadConfig(cf *conf.Conf) {
	this.UdpPort = cf.Int("udp_port", 14570)

	section, err := cf.Section("forwarder")
	if err != nil {
		panic("no forwarder config found")
	}
	this.Forwarder = new(ForwarderConfig)
	this.Forwarder.ToAddr = section.String("to_addr", ":5687")
	this.Forwarder.Backlog = section.Int("backlog", 1000)

	this.Inputs = make([]*InputConfig, 0)
	for i, _ := range cf.List("inputs", nil) {
		section, err = cf.Section(fmt.Sprintf("inputs[%d]", i))
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

type ForwarderConfig struct {
	ToAddr  string
	Backlog int
}

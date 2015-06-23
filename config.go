package main

import (
	"fmt"
	"strings"

	conf "github.com/nicholaskh/jsconf"
)

type GcollectorConfig struct {
	EtcServers []string

	App string

	UdpPort int

	Forwarder *ForwarderConfig

	Inputs []*InputConfig
}

func (this *GcollectorConfig) LoadConfig(cf *conf.Conf) {
	this.EtcServers = cf.StringList("etc_servers", nil)
	if this.EtcServers == nil {
		panic("No etc servers found")
	}

	this.App = cf.String("app", "")
	if this.App == "" {
		panic("No app specified")
	}

	this.UdpPort = cf.Int("udp_port", 14570)

	this.Inputs = make([]*InputConfig, 0)
	for i, _ := range cf.List("inputs", []interface{}{}) {
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

func (this *GcollectorConfig) LoadForwarder(cf *conf.Conf) {
	this.Forwarder = new(ForwarderConfig)
	var err error
	this.Forwarder.ToAddr, err = GetPiped()
	if err != nil {
		panic(err)
	}
	this.Forwarder.Backlog = cf.Int("forwarder_backlog", 1000)
}

type InputConfig struct {
	File  string
	Types []string
}

type ForwarderConfig struct {
	ToAddr  string
	Backlog int
}

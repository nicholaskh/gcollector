package main

import (
	"github.com/nicholaskh/golib/server"
)

var (
	GcollectorConf *GcollectorConfig
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

	conf := server.LoadConfig(options.configFile)
	GcollectorConf = new(GcollectorConfig)
	GcollectorConf.LoadConfig(conf)
}

func main() {
	gcollector := NewGcollector(GcollectorConf)
	gcollector.RunForever()
}

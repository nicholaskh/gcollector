package main

import (
	"fmt"
	"os"

	"github.com/nicholaskh/golib/locking"
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

	if options.kill {
		if err := server.KillProcess(options.lockFile); err != nil {
			fmt.Fprintf(os.Stderr, "stop failed: %s\n", err)
		}

		os.Exit(0)
	}

	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

	if options.lockFile != "" {
		if locking.InstanceLocked(options.lockFile) {
			fmt.Fprintf(os.Stderr, "Another gcollector is running, exit...\n")
			os.Exit(1)
		}

		locking.LockInstance(options.lockFile)
	}

	conf := server.LoadConfig(options.configFile)
	GcollectorConf = new(GcollectorConfig)
	GcollectorConf.LoadConfig(conf)
}

func main() {
	gcollector := NewGcollector(GcollectorConf)
	gcollector.RunForever()
}

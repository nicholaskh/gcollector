package main

import (
	"flag"
)

var (
	options struct {
		configFile   string
		logFile      string
		logLevel     string
		crashLogFile string
		showVersion  bool
		tick         int
		lockFile     string
		kill         bool
	}
)

func parseFlags() {
	flag.BoolVar(&options.kill, "k", false, "kill gcollectord")
	flag.StringVar(&options.lockFile, "lockfile", "gcollectord.lock", "lock file")
	flag.StringVar(&options.configFile, "conf", "etc/gcollector.cf", "config file")
	flag.BoolVar(&options.showVersion, "v", false, "show version and exit")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "info", "log level")
	flag.StringVar(&options.crashLogFile, "crashlog", "panic.dump", "crash log file")
	flag.IntVar(&options.tick, "tick", 60*10, "watchdog ticker length in seconds")

	flag.Parse()

	if options.tick <= 0 {
		panic("tick must be possitive")
	}
}

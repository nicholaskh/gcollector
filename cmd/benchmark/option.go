package main

import "flag"

var (
	options struct {
		concurrency  int
		requests     int
		last         int
		addr         string
		log          string
		logFile      string
		logLevel     string
		crashLogFile string
		showVersion  bool
	}
)

func parseFlags() {
	flag.IntVar(&options.concurrency, "c", 10, "how many goroutines")
	flag.IntVar(&options.requests, "n", 1000, "how many logs one goroutine generates")
	flag.IntVar(&options.last, "t", 10, "how long(in seconds) will last")
	flag.StringVar(&options.log, "f", "test.log", "which log file to put log in")
	flag.BoolVar(&options.showVersion, "v", false, "show version and exit")
	flag.StringVar(&options.logFile, "log", "stdout", "log file")
	flag.StringVar(&options.logLevel, "level", "info", "log level")
	flag.StringVar(&options.crashLogFile, "crashlog", "panic.dump", "crash log file")

	flag.Parse()

}

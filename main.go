package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/nicholaskh/etclib"
	"github.com/nicholaskh/golib/locking"
	"github.com/nicholaskh/golib/server"
	"github.com/nicholaskh/golib/signal"
	log "github.com/nicholaskh/log4go"
)

var (
	GcollectorConf *GcollectorConfig
)

func init() {
	parseFlags()

	if options.concurrency != 0 {
		runtime.GOMAXPROCS(options.concurrency)
	}

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	if options.kill {
		if err := server.KillProcess(options.lockFile); err != nil {
			fmt.Fprintf(os.Stderr, "stop failed: %s\n", err)
			os.Exit(1)
		}
		etclib.Dial(GcollectorConf.EtcServers)
		loadLocalAddr()
		UnregisterEtc()

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

	signal.RegisterSignalHandler(syscall.SIGINT, func(sig os.Signal) {
		shutdown()
	})

	conf := server.LoadConfig(options.configFile)
	GcollectorConf = new(GcollectorConfig)
	GcollectorConf.LoadConfig(conf)

	err := RegisterEtc(GcollectorConf.EtcServers)
	if err != nil {
		panic(err)
	}

	GcollectorConf.LoadForwarder(conf)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer func() {
		cleanup()

		if err := recover(); err != nil {
			fmt.Println(err)
			debug.PrintStack()
		}
	}()

	go server.RunSysStats(time.Now(), time.Duration(options.tick)*time.Second)

	gcollector := NewGcollector(GcollectorConf)
	gcollector.RunForever()
}

func shutdown() {
	cleanup()
	log.Info("Terminated")
	os.Exit(0)
}

func cleanup() {
	if options.lockFile != "" {
		locking.UnlockInstance(options.lockFile)
		log.Debug("Cleanup lock %s", options.lockFile)
	}
	UnregisterEtc()
}

package main

import (
	"io"
	"os"
	"time"

	"github.com/nicholaskh/golib/server"
	log "github.com/nicholaskh/log4go"
)

func init() {
	parseFlags()

	if options.showVersion {
		server.ShowVersionAndExit()
	}

	server.SetupLogging(options.logFile, options.logLevel, options.crashLogFile)

}

func main() {
	log_ := "NOTICE: 15-05-25 15:32:26 errno[0] client[10.77.141.87] uri[/commoncomponent?header=1&footer=1&menu=1&hlflag=child] user[15000000000209285] refer[] cookie[SESSIONID=02bef8adbvrquan5j1gjcsd9u3;U_UID=02bef8adbvrquan5j1gjcsd9u3;tempNoticeClosed=1;CITY_ID=110100;uid=15000000000209285;PHPSESSID=02bef8adbvrquan5j1gjcsd9u3] post[] ts[0.045583009719849]\n"
	var err error
	if checkFileIsExist(options.log) {
		log.Info("Log file: %s", options.log)
	} else {
		_, err = os.Create(options.log)
		if err != nil {
			log.Error("file not exists")
		}
		os.Exit(1)
	}
	waitGroup := make(chan interface{})
	for i := 0; i < options.concurrency; i++ {
		go func() {
			f, err := os.OpenFile(options.log, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				log.Error(err)
			}
			for j := 0; j < options.last; j++ {
				for k := 0; k < options.requests; k++ {
					appendLog(f, log_)
				}
				time.Sleep(time.Second)
			}
			waitGroup <- true
		}()
	}
	for i := 0; i < options.concurrency; i++ {
		<-waitGroup
	}
}

func appendLog(f *os.File, log_ string) {
	_, err := io.WriteString(f, log_)
	if err != nil {
		log.Error(err)
	}
}

func checkFileIsExist(filename string) bool {
	finfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if finfo.IsDir() {
		return false
	} else {
		return true
	}

}

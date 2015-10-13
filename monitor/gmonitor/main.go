package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"io"

	"flag"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/typingincolor/go-galen/monitor/monitor"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	var mongoHost = flag.String("mongo.host", "localhost", "Mongodb hostname")
	var influxHost = flag.String("influx.host", "localhost", "Influxdb hostname")
	var influxPort = flag.Int("influx.port", 8086, "Influxdb port")
	flag.Parse()

	Info.Println("Starting...")

	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		Info.Println("Stopping...")
		stopChan <- struct{}{}
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	database, err := mongo.Db(*mongoHost)
	if err != nil {
		Error.Println("failed to dial MongoDB", err)
		os.Exit(1)
	}
	defer database.Close()

	// start things
	monitorchan := make(chan monitor.Result)

	saveMonitorResultStoppedChan := monitor.InfluxSaver(monitorchan, *influxHost, *influxPort).Save()
	monitor := monitor.HTTPMonitor(stopChan, monitorchan, database)
	monitorStoppedChan := monitor.Start()

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()

	<-monitorStoppedChan
	close(monitorchan)
	<-saveMonitorResultStoppedChan
}

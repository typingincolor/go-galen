package main

import (
	log "gopkg.in/inconshreveable/log15.v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"flag"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/typingincolor/go-galen/monitor/monitor"
)

var logger = log.New(log.Ctx{"module": "main"})

func main() {
	var mongoHost = flag.String("mongo.host", "localhost", "Mongodb hostname")
	var influxHost = flag.String("influx.host", "localhost", "Influxdb hostname")
	var influxPort = flag.Int("influx.port", 8086, "Influxdb port")
	flag.Parse()

	logger.Info("Starting...")

	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		logger.Info("Stopping...")
		stopChan <- struct{}{}
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	database, err := mongo.Db(*mongoHost)
	if err != nil {
		logger.Error("failed to dial MongoDB", log.Ctx{"error": err})
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

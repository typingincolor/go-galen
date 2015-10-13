package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"flag"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/typingincolor/go-galen/monitor/monitor"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	var mongoHost = flag.String("mongo.host", "localhost", "Mongodb hostname")
	var influxHost = flag.String("influx.host", "localhost", "Influxdb hostname")
	var influxPort = flag.Int("influx.port", 8086, "Influxdb port")
	flag.Parse()

	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("Stopping...")
		stopChan <- struct{}{}
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	database, err := mongo.Db(*mongoHost)
	if err != nil {
		log.WithError(err).Fatal("failed to dial MongoDB")
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

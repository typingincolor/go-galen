package monitor

import (
	log "github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb/client"
	"github.com/typingincolor/go-galen/monitor/influx"
)

// Saver saves a result
type Saver interface {
	Save() <-chan struct{}
}

type consolesaver struct {
	monitorchan <-chan Result
}

type influxsaver struct {
	repo        influx.HealthCheckRepository
	monitorchan <-chan Result
}

func (consolesaver *consolesaver) Save() <-chan struct{} {
	stopchan := make(chan struct{}, 1)

	go func() {
		for monitor := range consolesaver.monitorchan {
			log.WithFields(log.Fields{"status_code": monitor.StatusCode}).Debug("saving")
		}
		stopchan <- struct{}{}
	}()

	return stopchan
}

func (influxsaver *influxsaver) Save() <-chan struct{} {
	stopchan := make(chan struct{}, 1)

	go func() {
		for monitor := range influxsaver.monitorchan {
			influxsaver.repo.Save(influx.HealthCheck{StatusCode: monitor.StatusCode})
		}
		stopchan <- struct{}{}
	}()

	return stopchan
}

// ConsoleSaver - write the output of a check to the console
func ConsoleSaver(monitorchan <-chan Result) Saver {
	return &consolesaver{monitorchan: monitorchan}
}

// InfluxSaver - save result to influxdb
func InfluxSaver(monitorchan <-chan Result, cfg client.Config) Saver {
	repo := influx.HealthCheckRepo(cfg)

	return &influxsaver{monitorchan: monitorchan, repo: repo}
}

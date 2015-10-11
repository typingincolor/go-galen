package monitor

import (
	log "github.com/Sirupsen/logrus"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"time"
)

// Result of a check
type Result struct {
	StatusCode int
}

type monitor struct {
	runner      APIClient
	stopchan    <-chan struct{}
	monitorchan chan<- Result
	db          mongo.Database
}

// Monitor - interface for something that monitors
type Monitor interface {
	Start() <-chan struct{}
}

// Start a monitor
func (m *monitor) Start() <-chan struct{} {
	log.Info("starting monitors")
	stoppedchan := make(chan struct{}, 1)

	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()

		for {
			select {
			case <-m.stopchan:
				log.Info("stopping monitor...")
				return
			default:
				m.monitor()
				log.Debug("  (waiting)")
				time.Sleep(10 * time.Second)
			}
		}
	}()
	return stoppedchan
}

// DummyMonitor that writes a canned result to the console
func DummyMonitor(stopChan <-chan struct{}, monitorchan chan<- Result, database mongo.Database) Monitor {
	return &monitor{runner: &dummyAPIClient{}, stopchan: stopChan, monitorchan: monitorchan, db: database}
}

func (m *monitor) loadMonitors() ([]mongo.HealthCheck, error) {
	return m.db.GetMonitors()
}

func (m *monitor) monitor() {
	log.Debug("running monitors")
	monitors, err := m.loadMonitors()

	if err != nil {
		log.WithError(err).Fatal("failed to load monitors")
		return
	}

	for _, mon := range monitors {
		res, _ := m.runner.Call(mon)
		m.monitorchan <- res
	}
}

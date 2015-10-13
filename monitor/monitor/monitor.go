package monitor

import (
	"github.com/typingincolor/go-galen/monitor/mongo"
	log "gopkg.in/inconshreveable/log15.v2"
	"time"
)

var logger = log.New(log.Ctx{"module": "monitor"})

// Result of a check
type Result struct {
	ID         string
	StatusCode int
	Elapsed    time.Duration
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
	logger.Info("starting monitor")
	stoppedchan := make(chan struct{}, 1)

	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()

		for {
			select {
			case <-m.stopchan:
				logger.Info("stopping monitor...")
				return
			default:
				m.monitor()
				logger.Debug("  (waiting)")
				time.Sleep(10 * time.Second)
			}
		}
	}()
	return stoppedchan
}

// DummyMonitor that writes a canned result selected result Saver
func DummyMonitor(stopChan <-chan struct{}, monitorchan chan<- Result, database mongo.Database) Monitor {
	return &monitor{runner: DummyAPIClient(), stopchan: stopChan, monitorchan: monitorchan, db: database}
}

// HTTPMonitor that uses a real API client to get the results
func HTTPMonitor(stopChan <-chan struct{}, monitorchan chan<- Result, database mongo.Database) Monitor {
	return &monitor{runner: HTTPAPIClient(), stopchan: stopChan, monitorchan: monitorchan, db: database}
}

func (m *monitor) loadMonitors() ([]mongo.HealthCheck, error) {
	return m.db.GetMonitors()
}

func (m *monitor) monitor() {
	logger.Debug("monitoring...")
	monitors, err := m.loadMonitors()

	if err != nil {
		logger.Error("failed to load monitors", log.Ctx{"error": err})
		return
	}

	for _, mon := range monitors {
		go func(mon mongo.HealthCheck) {
			res, _ := m.runner.Call(mon)
			m.monitorchan <- res
		}(mon)
	}
}

package monitor

import (
	"github.com/typingincolor/go-galen/monitor/mongo"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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
	log.Println("starting monitor")
	stoppedchan := make(chan struct{}, 1)

	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()

		for {
			select {
			case <-m.stopchan:
				log.Println("stopping monitor...")
				return
			default:
				m.monitor()
				log.Println("  (waiting)")
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
	log.Println("monitoring...")
	monitors, err := m.loadMonitors()

	if err != nil {
		log.Println("failed to load monitors", err)
		return
	}

	for _, mon := range monitors {
		go func(mon mongo.HealthCheck) {
			res, _ := m.runner.Call(mon)
			m.monitorchan <- res
		}(mon)
	}
}

package monitor

import (
	log "github.com/Sirupsen/logrus"
)

// Saver saves a result
type Saver interface {
	Save() <-chan struct{}
}

type consolesaver struct {
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

// ConsoleSaver - write the output of a check to the console
func ConsoleSaver(monitorchan <-chan Result) Saver {
	return &consolesaver{monitorchan: monitorchan}
}

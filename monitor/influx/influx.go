package influx

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb/client"
	"net/url"
	"time"
)

const (
	database        = "galen"
	retentionPolicy = "default"
)

// HealthCheck representation in InfluxDB
type HealthCheck struct {
	StatusCode int
	Elapsed    time.Duration
}

// HealthCheckRepository interface
type HealthCheckRepository interface {
	Save(HealthCheck) error
}

type healthCheckRepository struct {
	connection *client.Client
}

func (repo *healthCheckRepository) Save(h HealthCheck) error {
	log.WithFields(log.Fields{"status_code": h.StatusCode, "elapsed": h.Elapsed}).Debug("saving to influx")
	point := client.Point{
		Measurement: "healthcheck",
		Fields: map[string]interface{}{
			"status_code": h.StatusCode,
			"elapsed":     h.Elapsed.Seconds() * 1e3,
		},
		Time:      time.Now(),
		Precision: "s",
	}

	bps := client.BatchPoints{
		Points:          []client.Point{point},
		Database:        database,
		RetentionPolicy: retentionPolicy,
	}

	if _, err := repo.connection.Write(bps); err != nil {
		log.WithError(err).Fatal("unable to write to influxdb")
	}

	return nil
}

// HealthCheckRepo - create one...
func HealthCheckRepo(hostname string, port int) HealthCheckRepository {
	influxURL := fmt.Sprintf("http://%s:%d", hostname, port)

	log.WithField("url", influxURL).Info("Connecting to influxdb")
	u, err := url.Parse(influxURL)
	if err != nil {
		log.Fatal(err)
	}
	con, _ := client.NewClient(client.Config{URL: *u})
	return &healthCheckRepository{connection: con}
}

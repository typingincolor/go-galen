package influx

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"gopkg.in/inconshreveable/log15.v2"
	"net/url"
	"os"
	"time"
)

const (
	database        = "galen"
	retentionPolicy = "default"
)

var logger = log15.New(log15.Ctx{"module": "influx"})

// HealthCheck representation in InfluxDB
type HealthCheck struct {
	ID         string
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
	logger.Debug("saving to influx", log15.Ctx{"point": h})
	point := client.Point{
		Measurement: "healthcheck",
		Tags: map[string]string{
			"monitor": h.ID,
		},
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
		logger.Crit("unable to write to influxdb", log15.Ctx{"error": err})
		return err
	}

	return nil
}

// HealthCheckRepo - create one...
func HealthCheckRepo(hostname string, port int) HealthCheckRepository {
	influxURL := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Info("Connecting to influxdb", log15.Ctx{"url": influxURL})
	u, err := url.Parse(influxURL)
	if err != nil {
		logger.Crit("error parsing influx url", log15.Ctx{"error": err})
		os.Exit(1)
	}
	con, _ := client.NewClient(client.Config{URL: *u})
	return &healthCheckRepository{connection: con}
}

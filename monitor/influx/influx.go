package influx

import (
	log "github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb/client"
	"time"
)

const (
	database        = "galen"
	retentionPolicy = "default"
)

// HealthCheck representation in InfluxDB
type HealthCheck struct {
	StatusCode int
}

// HealthCheckRepository interface
type HealthCheckRepository interface {
	Save(HealthCheck) error
}

type healthCheckRepository struct {
	connection *client.Client
}

func (repo *healthCheckRepository) Save(h HealthCheck) error {
	log.WithFields(log.Fields{"status_code": h.StatusCode}).Debug("saving to influx")
	point := client.Point{
		Measurement: "healthcheck",
		Fields: map[string]interface{}{
			"status_code": h.StatusCode,
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

// HealthCheckRepository - create one...
func HealthCheckRepo(cfg client.Config) HealthCheckRepository {
	con, _ := client.NewClient(cfg)
	return &healthCheckRepository{connection: con}
}

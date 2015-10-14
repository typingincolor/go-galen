package influx

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	database        = "galen"
	retentionPolicy = "default"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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
	log.Printf("saving to influx point: %+v", h)
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
		log.Println("unable to write to influxdb:", err)
		return err
	}

	return nil
}

// HealthCheckRepo - create one...
func HealthCheckRepo(hostname string, port int) HealthCheckRepository {
	influxURL := fmt.Sprintf("http://%s:%d", hostname, port)

	log.Println("Connecting to influxdb url:", influxURL)
	u, err := url.Parse(influxURL)
	if err != nil {
		log.Fatalln("error parsing influx url: ", err)
		os.Exit(1)
	}
	con, _ := client.NewClient(client.Config{URL: *u})
	return &healthCheckRepository{connection: con}
}

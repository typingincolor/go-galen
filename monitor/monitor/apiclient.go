package monitor

import (
	"errors"
	"fmt"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"gopkg.in/inconshreveable/log15.v2"
	"net/http"
	"strings"
	"time"
)

// APIClient calls the HealthCheck
type APIClient interface {
	Call(monitor mongo.HealthCheck) (Result, error)
}

type dummyAPIClient struct{}

func (m *dummyAPIClient) Call(monitor mongo.HealthCheck) (Result, error) {
	logger.Debug("calling", log15.Ctx{"id": monitor.ID.Hex()})
	return Result{StatusCode: 200}, nil
}

type apiClient struct{}

func (client *apiClient) Call(monitor mongo.HealthCheck) (Result, error) {
	logger.Debug("calling", log15.Ctx{"id": monitor.ID.Hex(), "url": monitor.URL})
	if strings.ToUpper(monitor.Method) == "GET" {
		start := time.Now()
		resp, err := http.Get(monitor.URL)

		if err != nil {
			return Result{}, errors.New(err.Error())
		}
		defer resp.Body.Close()

		return Result{ID: monitor.ID.Hex(), StatusCode: resp.StatusCode, Elapsed: time.Since(start)}, nil
	}

	return Result{}, fmt.Errorf("API client has not implemented method %s", monitor.Method)
}

// HTTPAPIClient - returns a real API client
func HTTPAPIClient() APIClient {
	return &apiClient{}
}

// DummyAPIClient - returns a dummy API client
func DummyAPIClient() APIClient {
	return &dummyAPIClient{}
}

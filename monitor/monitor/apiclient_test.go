package monitor_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/typingincolor/go-galen/monitor/monitor"
	"testing"
)

var client = monitor.HTTPAPIClient()

func TestDummyApiClient(t *testing.T) {
	dummy := monitor.DummyAPIClient()

	res, err := dummy.Call(mongo.HealthCheck{Method: "GET", URL: "http://999.999.999.999"})

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestItFailsWithUnknownHttpMethod(t *testing.T) {
	expectedError := errors.New("API client has not implemented method POST")

	_, err := client.Call(mongo.HealthCheck{Method: "POST", URL: "http://example.com"})

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, expectedError, err)
	}
}

func TestItHandlesAnUnknownHost(t *testing.T) {
	_, err := client.Call(mongo.HealthCheck{Method: "GET", URL: "http://999.999.999.999"})

	if assert.Error(t, err, "An error was expected") {
		assert.Contains(t, err.Error(), "Get http://999.999.999.999: dial tcp: lookup 999.999.999.999")
	}
}

func TestItCanCallAnUrl(t *testing.T) {
	res, err := client.Call(mongo.HealthCheck{Method: "GET", URL: "http://echo.jsontest.com/key/value/one/two"})

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestItCanHandleANon200Response(t *testing.T) {
	res, err := client.Call(mongo.HealthCheck{Method: "GET", URL: "http://httpstat.us/500"})

	assert.Nil(t, err)
	assert.Equal(t, 500, res.StatusCode)
}

package monitor_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/typingincolor/go-galen/monitor/monitor"
	"testing"
	"log"
	"io"
	"os"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

var client = monitor.HTTPAPIClient()

func Init(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func TestDummyApiClient(t *testing.T) {
	Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	Info.Println("hello")

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
	res, err := client.Call(mongo.HealthCheck{Method: "GET", URL: "http://jsonplaceholder.typicode.com/posts/1"})

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestItCanHandleANon200Response(t *testing.T) {
	res, err := client.Call(mongo.HealthCheck{Method: "GET", URL: "http://httpstat.us/500"})

	assert.Nil(t, err)
	assert.Equal(t, 500, res.StatusCode)
}

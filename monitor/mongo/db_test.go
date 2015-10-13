package mongo_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"os"
	"testing"
)

func TestDial(t *testing.T) {
	if os.Getenv("TEST_MONGO") != "true" {
		t.Skip("skipping test - TEST_MONGO not set")
	}

	db, err := mongo.Db("localhost")
	defer db.Close()

	assert.Nil(t, err)
}

func TestGetMonitors(t *testing.T) {
	if os.Getenv("TEST_MONGO") != "true" {
		t.Skip("skipping test - TEST_MONGO not set")
	}

	db, err := mongo.Db("localhost")
	defer db.Close()

	assert.Nil(t, err)

	healthchecks, err := db.GetMonitors()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(healthchecks))
}

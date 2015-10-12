package mongo_test
import (
	"testing"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/stretchr/testify/assert"
	"os"
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
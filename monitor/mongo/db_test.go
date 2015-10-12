package mongo_test
import (
	"testing"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestDial(t *testing.T) {
	if os.Getenv("TRAVIS") != "true" {
		t.Skip("skipping test - only runs on travis")
	}

	db, err := mongo.Db("localhost")
	defer db.Close()

	assert.Nil(t, err)
}
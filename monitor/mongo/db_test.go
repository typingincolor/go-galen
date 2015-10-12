package mongo_test
import (
	"testing"
	"github.com/typingincolor/go-galen/monitor/mongo"
	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	db, err := mongo.Db("localhost")
	defer db.Close()

	assert.Nil(t, err)
}
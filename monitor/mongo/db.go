package mongo

import (
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var logger = log15.New(log15.Ctx{"module": "mongo"})

type db struct {
	host    string
	session *mgo.Session
}

// HealthCheck to run
type HealthCheck struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	URL    string
	Method string
}

// Database - Interface defining database operations
type Database interface {
	Close()
	GetMonitors() ([]HealthCheck, error)
}

func (d *db) dial() error {
	var err error
	logger.Info("dialing mongodb", log15.Ctx{"mongo_host": d.host})
	d.session, err = mgo.Dial(d.host)
	return err
}

func (d *db) Close() {
	d.session.Close()
	logger.Info("closed mongodb connection")
}

func (d *db) GetMonitors() ([]HealthCheck, error) {
	var monitors []HealthCheck
	iter := d.session.DB("monitors").C("apis").Find(nil).Iter()
	var cfg HealthCheck
	for iter.Next(&cfg) {
		monitors = append(monitors, cfg)
	}
	iter.Close()

	return monitors, iter.Err()
}

// Db - initialise a Database implementation
func Db(host string) (Database, error) {
	res := &db{host: host}
	err := res.dial()
	return res, err
}

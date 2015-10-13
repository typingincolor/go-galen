package mongo

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

type db struct {
	host    string
	session *mgo.Session
}

// HealthCheck to run
type HealthCheck struct {
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
	log.WithField("mongo_host", d.host).Info("dialing mongodb")
	d.session, err = mgo.Dial(d.host)
	return err
}

func (d *db) Close() {
	d.session.Close()
	log.Info("closed mongodb connection")
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

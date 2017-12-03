package db

import (
	"github.com/micro-company/go-auth/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	log = logrus.New()

	Session *mgo.Session
	Mongo   *mgo.DialInfo
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func ConnectToMongo() {
	// Get configuration
	MONGO_URL := utils.Getenv("MONGO_URL", "localhost/auth")

	s, err := mgo.Dial(MONGO_URL)
	if err != nil {
		log.Panic("Fail connect to Mongo")
		panic(err)
	}
	s.SetSafe(&mgo.Safe{})
	log.Info("Success connect to MongoDB")
	Session = s
}

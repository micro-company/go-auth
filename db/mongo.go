package db

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"github.com/batazor/go-auth/utils"
)

var (
	log = logrus.New()

	Session *mgo.Session
	Mongo *mgo.DialInfo
)

func Connect() {
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
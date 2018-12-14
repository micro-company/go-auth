package mongodb

import (
	"context"
	"github.com/micro-company/go-auth/utils"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()

	Session *mongo.Client
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
	MONGO_URL := utils.Getenv("MONGO_URL", "mongodb://localhost/auth")

	client, err := mongo.Connect(context.TODO(), MONGO_URL)
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
	}

	log.Info("Success connect to MongoDB")
	Session = client
}

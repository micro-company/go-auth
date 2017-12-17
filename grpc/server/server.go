package server

import (
	"github.com/micro-company/go-auth/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log = logrus.New()

	err        error
	GrpcClient *grpc.ClientConn
)

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

func init() {
	// Get configuration
	API_MAIL_ADDRESS := utils.Getenv("API_MAIL_ADDRESS", "localhost:50051")

	// Connect to API by gRPC
	GrpcClient, err = grpc.Dial(API_MAIL_ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Error("did not connect: ", err)
	}

	log.Info("Success connect to Mail API by gRPC")
}

func GetConnClient() *grpc.ClientConn {
	return GrpcClient
}

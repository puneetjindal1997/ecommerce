package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"ecommerce/constant"
	"ecommerce/types"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type manager struct {
	connection *mongo.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

var Mgr Manager

type Manager interface {
	Insert(interface{}, string) error
	GetSingleRecordByEmail(string, string) *types.Verification
	UpdateVerification(types.Verification, string) error
	UpdateEmailVerifiedStatus(types.Verification, string) error
}

func ConnectDb() {
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = constant.MDBUri
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("%s%s", "mongodb://", uri)))

	if err != nil {
		ConnectDb()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Println("Unable to initialize database connectors. Retrying...")
		ConnectDb()
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Println("Unable to connect to the database. Retrying...")
		ConnectDb()
	}
	log.Println("Successfully connected to the database at %s", uri)

	Mgr = &manager{connection: client, ctx: ctx, cancel: cancel}
}

func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

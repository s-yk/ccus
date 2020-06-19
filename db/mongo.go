package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type TypeInfo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Project   string             `json:"project" bson:"Project,omitempty"`
	Namespace string             `json:"namespace" bson:"Namespace,omitempty"`
	Class     string             `json:"class" bson:"Class,omitempty"`
	Type      string             `json:"type" bson:"Type,omitempty"`
}

var Client *mongo.Client

func GetConnection() error {

	host := os.Getenv("CCUS_MONGO_HOST")
	if host == "" {
		return fmt.Errorf("Empty: %s", "CCUS_MONGO_HOST")
	}

	port := os.Getenv("CCUS_MONGO_PORT")
	if port == "" {
		return fmt.Errorf("Empty: %s", "CCUS_MONGO_PORT")
	}

	user := os.Getenv("CCUS_MONGO_USER")
	if user == "" {
		return fmt.Errorf("Empty: %s", "CCUS_MONGO_USER")
	}

	pwd := os.Getenv("CCUS_MONGO_PASSWORD")
	if pwd == "" {
		return fmt.Errorf("Empty: %s", "CCU_MONGO_PASSWORD")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/ccu", user, pwd, host, port)))
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("connection error: %s", err)
	}

	return nil
}

func GetData(filter interface{}) ([]bson.M, error) {

	d := Client.Database("ccu")
	cl := d.Collection("ccu")

	opt := options.Find()
	opt.SetSort(bson.D{{Key: "_id", Value: 1}})

	c, err := cl.Find(context.Background(), filter, opt)
	if err != nil {
		return nil, err
	}
	defer c.Close(context.Background())

	var ts []bson.M
	if err := c.All(context.Background(), &ts); err != nil {
		return nil, err
	}

	return ts, nil
}

func Disconnect() error {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(c); err != nil {
		return err
	}

	return nil
}

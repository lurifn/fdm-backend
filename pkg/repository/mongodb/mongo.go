package mongodb

import (
	"context"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoDB struct {
	URI        string
	DB         string
	Collection string
	Username   string
	Password   string
}

func unmarshalBson(obj string) (bson.M, error) {
	doc := bson.M{}
	err := bson.UnmarshalExtJSON([]byte(obj), true, &doc)
	return doc, err
}

func (config MongoDB) Save(message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	credential := options.Credential{Username: config.Username, Password: config.Password}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI).SetAuth(credential))

	defer func() {
		cancel()

		if err = client.Disconnect(ctx); err != nil {
			log.Error.Panic(err)
		}
	}()

	if err != nil {
		log.Error.Println("Error connecting with mongo:", err.Error())
		return "", err
	}

	// verify that we are really connected
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error.Println("Error connecting with mongo:", err.Error())
		return "", err
	}

	log.Debug.Println("Successfully connected with mongo")

	doc, err := unmarshalBson(message)
	if err != nil {
		log.Error.Println("Error unmarshalling", message, err.Error())
		return "", err
	}

	log.Debug.Println("Order bson prepared")

	db := client.Database(config.DB)
	collection := db.Collection(config.Collection)
	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		log.Error.Println("Error inserting ", doc, err.Error())

		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	log.Info.Println("Document with ID ", id, " stored in the DB collection ", config.Collection)

	return id, nil
}

package mongodb

import (
	"context"
	"fmt"
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

func unmarshalBsonMap(obj string) (map[string]bson.M, error) {
	doc := make(map[string]bson.M)
	err := bson.UnmarshalExtJSON([]byte(obj), true, &doc)
	if err != nil {
		log.Error.Println("Error unmarshalling bson map:", err)
	}

	return doc, err
}

func unmarshalBson(obj string) (bson.M, error) {
	doc := bson.M{}
	err := bson.UnmarshalExtJSON([]byte(obj), true, &doc)
	if err != nil {
		log.Error.Println("Error unmarshalling bson:", err)
	}

	return doc, err
}

type operationFunc func(record string, config MongoDB, client *mongo.Client, ctx context.Context) (string, error)

func transactionWrapper(f operationFunc) func(record string, config MongoDB) (string, error) {
	return func(record string, config MongoDB) (string, error) {
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

		return f(record, config, client, ctx)
	}
}

func save(record string, config MongoDB, client *mongo.Client, ctx context.Context) (string, error) {
	doc, err := unmarshalBson(record)
	if err != nil {
		log.Error.Println("Error unmarshalling", record, err.Error())
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

func update(record string, config MongoDB, client *mongo.Client, ctx context.Context) (string, error) {
	doc, err := unmarshalBsonMap(record)
	if err != nil {
		log.Error.Println("Error unmarshalling", record, err.Error())
		return "", err
	}

	if len(doc) != 1 {
		return "validation error", fmt.Errorf("update func only accepts one document at a time")
	}

	log.Debug.Println("Order bson prepared")

	db := client.Database(config.DB)
	collection := db.Collection(config.Collection)

	for id, content := range doc {
		hex, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Error.Println("Invalid ID", id, err)
			return "", err
		}

		res, err := collection.ReplaceOne(ctx, bson.M{"_id": hex}, content)
		if err != nil {
			log.Error.Println("Error inserting ", doc, err.Error())

			return "", err
		}

		log.Debug.Println("Document updated with result", res)

		// prepares obj to be returned
		if res.UpsertedCount+res.ModifiedCount == 1 {
			content["_id"] = id
			bytes, err := bson.Marshal(content)
			if err != nil {
				log.Error.Println("Error preparing response", err)
				return "", err
			}

			return string(bytes), nil
		}
	}

	return "ok", nil
}

func (config MongoDB) Save(record string) (string, error) {
	f := transactionWrapper(save)
	return f(record, config)
}

func (config MongoDB) Update(record string) (string, error) {
	f := transactionWrapper(update)
	return f(record, config)
}

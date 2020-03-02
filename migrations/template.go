package migrations

import (
	"context"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		opt := options.Index().SetName("index")
		keys := bson.D{{"key", 1}}
		model := mongo.IndexModel{Keys: keys, Options: opt}
		_, err := db.Collection("coll").Indexes().CreateOne(context.TODO(), model)

		if err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		_, err := db.Collection("coll").Indexes().DropOne(context.TODO(), "index")
		if err != nil {
			return err
		}
		return nil
	})
}

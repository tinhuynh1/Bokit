package database

import (
	"context"
	"fmt"
	"quiz-svc/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func NewMongoDB(cfg *config.MongoDBConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	fmt.Println("c", c)
	clientOptions := options.Client().ApplyURI(c)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client.Database(cfg.DBName), nil
}

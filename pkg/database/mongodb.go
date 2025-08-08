package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fampay-youtube-api/internal/config"
)

func NewMongoDB(cfg config.MongoDBConfig) (*mongo.Database, error) {
	// Create MongoDB client
	clientOptions := options.Client().ApplyURI(cfg.URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)

	// Create indexes
	if err := createIndexes(db); err != nil {
		log.Printf("Warning: Failed to create some indexes: %v", err)
	}

	log.Println("MongoDB connected successfully")
	return db, nil
}

func createIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	videosCollection := db.Collection("videos")

	// Create indexes for performance
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{"video_id", 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{"published_at", -1}},
		},
		{
			Keys: bson.D{{"title", "text"}, {"description", "text"}},
		},
		{
			Keys: bson.D{
				{"published_at", -1},
				{"title", 1},
				{"description", 1},
			},
		},
	}

	_, err := videosCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("MongoDB indexes created successfully")
	return nil
}

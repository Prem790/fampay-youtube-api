package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fampay-youtube-api/internal/models"
	"fampay-youtube-api/internal/repository"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	db := client.Database("fampay_youtube")
	repo := repository.NewVideoRepository(db)

	// Test videos
	testVideos := []*models.Video{
		{
			VideoID:      "test_cricket_1",
			Title:        "India vs Australia Cricket Highlights",
			Description:  "Amazing cricket match with spectacular batting and bowling",
			PublishedAt:  time.Now().Add(-1 * time.Hour),
			ChannelTitle: "Cricket World",
			ChannelID:    "UC_cricket_1",
			ThumbnailURL: models.Thumbnail{
				Default: "https://i.ytimg.com/vi/test1/default.jpg",
				Medium:  "https://i.ytimg.com/vi/test1/mqdefault.jpg",
				High:    "https://i.ytimg.com/vi/test1/hqdefault.jpg",
			},
		},
		{
			VideoID:      "test_football_1",
			Title:        "Football Skills and Techniques",
			Description:  "Learn professional football moves and tricks",
			PublishedAt:  time.Now().Add(-2 * time.Hour),
			ChannelTitle: "Football Pro",
			ChannelID:    "UC_football_1",
			ThumbnailURL: models.Thumbnail{
				Default: "https://i.ytimg.com/vi/test2/default.jpg",
				Medium:  "https://i.ytimg.com/vi/test2/mqdefault.jpg",
				High:    "https://i.ytimg.com/vi/test2/hqdefault.jpg",
			},
		},
		{
			VideoID:      "test_cricket_2",
			Title:        "Cricket World Cup Final Analysis",
			Description:  "Detailed analysis of the cricket championship final match",
			PublishedAt:  time.Now().Add(-3 * time.Hour),
			ChannelTitle: "Sports Analytics",
			ChannelID:    "UC_sports_1",
			ThumbnailURL: models.Thumbnail{
				Default: "https://i.ytimg.com/vi/test3/default.jpg",
				Medium:  "https://i.ytimg.com/vi/test3/mqdefault.jpg",
				High:    "https://i.ytimg.com/vi/test3/hqdefault.jpg",
			},
		},
		{
			VideoID:      "test_tutorial_1",
			Title:        "How to improve cricket batting technique",
			Description:  "Step by step tutorial for better cricket batting",
			PublishedAt:  time.Now().Add(-4 * time.Hour),
			ChannelTitle: "Cricket Coach",
			ChannelID:    "UC_coach_1",
			ThumbnailURL: models.Thumbnail{
				Default: "https://i.ytimg.com/vi/test4/default.jpg",
				Medium:  "https://i.ytimg.com/vi/test4/mqdefault.jpg",
				High:    "https://i.ytimg.com/vi/test4/hqdefault.jpg",
			},
		},
	}

	// Insert test videos
	inserted := 0
	for _, video := range testVideos {
		// Check if video already exists
		existing, _ := repo.GetByVideoID(video.VideoID)
		if existing == nil {
			if err := repo.Create(video); err != nil {
				log.Printf("Failed to insert video %s: %v", video.VideoID, err)
				continue
			}
			inserted++
		}
	}

	log.Printf("Successfully loaded %d test videos", inserted)
}

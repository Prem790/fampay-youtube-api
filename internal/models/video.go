package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VideoID      string             `json:"video_id" bson:"video_id"`
	Title        string             `json:"title" bson:"title"`
	Description  string             `json:"description" bson:"description"`
	PublishedAt  time.Time          `json:"published_at" bson:"published_at"`
	ChannelTitle string             `json:"channel_title" bson:"channel_title"`
	ChannelID    string             `json:"channel_id" bson:"channel_id"`
	SearchQuery  string             `json:"search_query" bson:"search_query"` // Track origin query
	ThumbnailURL Thumbnail          `json:"thumbnails" bson:"thumbnails"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type Thumbnail struct {
	Default string `json:"default" bson:"default"`
	Medium  string `json:"medium" bson:"medium"`
	High    string `json:"high" bson:"high"`
}

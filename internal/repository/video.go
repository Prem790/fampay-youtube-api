package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fampay-youtube-api/internal/models"
)

type VideoRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewVideoRepository(db *mongo.Database) *VideoRepository {
	return &VideoRepository{
		db:         db,
		collection: db.Collection("videos"),
	}
}

// Enhanced search with better partial matching
func (r *VideoRepository) Search(query string, page, pageSize int, sortBy string) ([]models.Video, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build search filter for partial matching
	filter := r.buildSearchFilter(query)

	// Count total matching documents
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Build sort options
	sortOptions := r.buildSortOptions(sortBy)

	// Calculate skip
	skip := (page - 1) * pageSize

	// Find options
	findOptions := options.Find()
	findOptions.SetSort(sortOptions)
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	// Find videos
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search videos: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []models.Video
	if err = cursor.All(ctx, &videos); err != nil {
		return nil, 0, fmt.Errorf("failed to decode search results: %w", err)
	}

	return videos, total, nil
}

func (r *VideoRepository) buildSearchFilter(query string) bson.M {
	if strings.TrimSpace(query) == "" {
		return bson.M{}
	}

	// Split query into words for better partial matching
	words := strings.Fields(strings.ToLower(query))

	var orConditions []bson.M

	// For each word, create regex patterns for title and description
	for _, word := range words {
		// Escape special regex characters
		escapedWord := strings.ReplaceAll(word, `\`, `\\`)
		escapedWord = strings.ReplaceAll(escapedWord, `.`, `\.`)

		wordConditions := []bson.M{
			{"title": bson.M{"$regex": escapedWord, "$options": "i"}},
			{"description": bson.M{"$regex": escapedWord, "$options": "i"}},
		}
		orConditions = append(orConditions, bson.M{"$or": wordConditions})
	}

	// All words must match (AND logic), but can match in either title OR description
	if len(orConditions) == 1 {
		return orConditions[0]
	}

	return bson.M{"$and": orConditions}
}

func (r *VideoRepository) buildSortOptions(sortBy string) bson.D {
	switch sortBy {
	case "oldest":
		return bson.D{{"published_at", 1}}
	case "title":
		return bson.D{{"title", 1}}
	case "channel":
		return bson.D{{"channel_title", 1}, {"published_at", -1}}
	case "latest":
		return bson.D{{"published_at", -1}}
	default:
		return bson.D{{"published_at", -1}}
	}
}

// Enhanced get paginated with sorting
func (r *VideoRepository) GetPaginated(page, pageSize int, sortBy string) ([]models.Video, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count videos: %w", err)
	}

	// Build sort options
	sortOptions := r.buildSortOptions(sortBy)

	// Calculate skip
	skip := (page - 1) * pageSize

	// Find options
	findOptions := options.Find()
	findOptions.SetSort(sortOptions)
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64(skip))

	// Find videos
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find videos: %w", err)
	}
	defer cursor.Close(ctx)

	var videos []models.Video
	if err = cursor.All(ctx, &videos); err != nil {
		return nil, 0, fmt.Errorf("failed to decode videos: %w", err)
	}

	return videos, total, nil
}

func (r *VideoRepository) Create(video *models.Video) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	video.CreatedAt = time.Now()
	video.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, video)
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}

	return nil
}

func (r *VideoRepository) GetByVideoID(videoID string) (*models.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var video models.Video
	err := r.collection.FindOne(ctx, bson.M{"video_id": videoID}).Decode(&video)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Video not found
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return &video, nil
}

func (r *VideoRepository) GetLatest() (*models.Video, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"published_at", -1}})

	var video models.Video
	err := r.collection.FindOne(ctx, bson.M{}, findOptions).Decode(&video)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No videos found
		}
		return nil, fmt.Errorf("failed to get latest video: %w", err)
	}

	return &video, nil
}

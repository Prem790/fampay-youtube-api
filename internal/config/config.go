package config

import (
    "os"
    "strconv"
    "strings"
    "github.com/joho/godotenv"
)

type Config struct {
    Server   ServerConfig
    MongoDB  MongoDBConfig
    Redis    RedisConfig
    YouTube  YouTubeConfig
}

type ServerConfig struct {
    Port string
    Host string
    Mode string
}

type MongoDBConfig struct {
    URI      string
    Database string
}

type RedisConfig struct {
    Host     string
    Port     string
    Password string
    DB       int
}

type YouTubeConfig struct {
    APIKeys            []string
    SearchQueries      []string
    FetchInterval      int
    MaxResultsPerQuery int
    RegionCode         string
    RelevanceLanguage  string
}

func Load() (*Config, error) {
    godotenv.Load()

    // Parse and clean search queries
    searchQueriesStr := getEnv("YOUTUBE_SEARCH_QUERIES", "cricket,football,technology,music,gaming")
    searchQueries := strings.Split(searchQueriesStr, ",")
    for i, query := range searchQueries {
        searchQueries[i] = strings.TrimSpace(query)
    }

    // Parse and clean API keys
    apiKeysStr := getEnv("YOUTUBE_API_KEYS", "")
    apiKeys := strings.Split(apiKeysStr, ",")
    var cleanKeys []string
    for _, key := range apiKeys {
        if trimmed := strings.TrimSpace(key); trimmed != "" {
            cleanKeys = append(cleanKeys, trimmed)
        }
    }

    config := &Config{
        Server: ServerConfig{
            Port: getEnv("PORT", "8080"),
            Host: getEnv("HOST", "localhost"),
            Mode: getEnv("GIN_MODE", "debug"),
        },
        MongoDB: MongoDBConfig{
            URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
            Database: getEnv("MONGODB_DATABASE", "fampay_youtube"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnv("REDIS_PORT", "6379"),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvInt("REDIS_DB", 0),
        },
        YouTube: YouTubeConfig{
            APIKeys:            cleanKeys,
            SearchQueries:      searchQueries,
            FetchInterval:      getEnvInt("FETCH_INTERVAL", 10),
            MaxResultsPerQuery: getEnvInt("MAX_RESULTS_PER_QUERY", 50),
            RegionCode:         getEnv("REGION_CODE", "IN"),
            RelevanceLanguage:  getEnv("RELEVANCE_LANGUAGE", "en"),
        },
    }

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

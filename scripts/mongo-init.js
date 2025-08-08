// MongoDB initialization script
db = db.getSiblingDB('fampay_youtube');

// Create the videos collection
db.createCollection('videos');

// Create indexes for better performance
db.videos.createIndex({ "video_id": 1 }, { unique: true });
db.videos.createIndex({ "published_at": -1 });
db.videos.createIndex({ "title": "text", "description": "text" });
db.videos.createIndex({ "channel_title": 1, "published_at": -1 });
db.videos.createIndex({ "search_query": 1, "published_at": -1 });

print('✅ MongoDB initialization completed - indexes created');
print('📊 Database: fampay_youtube');
print('📚 Collections: videos');
print('🔍 Indexes: video_id, published_at, text_search, channel_title, search_query');
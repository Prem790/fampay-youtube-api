#!/bin/bash

echo "Setting up FamPay YouTube API project..."

# Create .env file from example
if [ ! -f .env ]; then
    cp .env.example .env
    echo "Created .env file. Please update it with your YouTube API keys."
fi

# Install Go dependencies
echo "Installing Go dependencies..."
go mod download

# Start Docker services
echo "Starting Docker services..."
docker-compose up -d postgres redis

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 10

# Run migrations
echo "Running database migrations..."
docker-compose exec postgres psql -U postgres -d fampay_youtube -f /docker-entrypoint-initdb.d/001_create_videos_table.sql

echo "Setup completed! You can now run 'make run' to start the server."

#!/bin/bash

# LoveApp Backend Startup Script
# This script helps you get started with the LoveApp Backend

set -e

echo "🚀 LoveApp Backend Startup Script"
echo "=================================="

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from example.env..."
    cp example.env .env
    echo "✅ .env file created. You can modify it if needed."
else
    echo "✅ .env file already exists."
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "✅ Docker is running."

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

echo "✅ docker-compose is available."

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
docker-compose --env-file .env down > /dev/null 2>&1 || true

# Build and start containers
echo "🏗️  Building and starting containers..."
docker-compose --env-file .env up -d --build

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if backend is healthy
echo "🏥 Checking backend health..."
for i in {1..30}; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Backend is healthy!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Backend health check failed after 30 attempts."
        echo "📋 Showing backend logs:"
        docker-compose --env-file .env logs loveapp-backend
        exit 1
    fi
    echo "⏳ Attempt $i/30: Backend not ready yet, waiting..."
    sleep 2
done

echo ""
echo "🎉 LoveApp Backend is now running!"
echo "=================================="
echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
echo "🏥 Health Check: http://localhost:8080/health"
echo "🗄️  Database: PostgreSQL running on localhost:5432"
echo ""
echo "👥 Default Users:"
echo "   Username: anyel    | Password: password"
echo "   Username: alexis   | Password: password"
echo ""
echo "🔧 Useful Commands:"
echo "   View logs: make logs"
echo "   Stop services: make docker-stop"
echo "   Restart services: make docker-rebuild"
echo ""
echo "Happy coding! 💕"
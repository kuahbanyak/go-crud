#!/bin/bash
# Production Deployment Script for Docker

set -e

echo "🚀 Starting Production Deployment..."

# Check if .env.prod exists
if [ ! -f .env.prod ]; then
    echo "❌ .env.prod file not found!"
    echo "📋 Please copy .env.prod.template to .env.prod and fill in your values"
    exit 1
fi

# Load environment variables
export $(grep -v '^#' .env.prod | xargs)

echo "📦 Building production Docker images..."
docker-compose -f docker-compose.prod.yml build --no-cache

echo "🔄 Stopping existing containers..."
docker-compose -f docker-compose.prod.yml down

echo "🚀 Starting production services..."
docker-compose -f docker-compose.prod.yml up -d

echo "⏳ Waiting for services to be ready..."
sleep 30

echo "🔍 Checking service health..."
docker-compose -f docker-compose.prod.yml ps

echo "✅ Production deployment completed!"
echo "🌐 Your API should be available at: http://your-domain.com"
echo "📊 Check logs with: docker-compose -f docker-compose.prod.yml logs -f"

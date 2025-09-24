@echo off
REM Production Deployment Script for Windows

echo 🚀 Starting Production Deployment...

REM Check if .env.prod exists
if not exist .env.prod (
    echo ❌ .env.prod file not found!
    echo 📋 Please copy .env.prod.template to .env.prod and fill in your values
    pause
    exit /b 1
)

echo 📦 Building production Docker images...
docker-compose -f docker-compose.prod.yml build --no-cache

echo 🔄 Stopping existing containers...
docker-compose -f docker-compose.prod.yml down

echo 🚀 Starting production services...
docker-compose -f docker-compose.prod.yml up -d

echo ⏳ Waiting for services to be ready...
timeout /t 30 /nobreak

echo 🔍 Checking service health...
docker-compose -f docker-compose.prod.yml ps

echo ✅ Production deployment completed!
echo 🌐 Your API should be available at: http://your-domain.com
echo 📊 Check logs with: docker-compose -f docker-compose.prod.yml logs -f
pause

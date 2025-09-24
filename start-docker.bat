@echo off
echo Building and starting Go CRUD service with Docker...

REM Stop any existing containers
docker-compose down

REM Build and start the services
docker-compose up --build -d

echo.
echo Services are starting up...
echo API will be available at: http://localhost:8081
echo SQL Server will be available at: localhost:1433
echo.
echo To view logs, run: docker-compose logs -f
echo To stop services, run: docker-compose down
echo.

pause

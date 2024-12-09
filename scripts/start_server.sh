#!/bin/bash

docker-compose down -v
docker-compose up -d

# Wait for the database to be up
until nc -z -v -w30 localhost 5432
do
    echo "Waiting for database connection..."
    sleep 1
done

echo "Database is up!"

# Build and run the Docker container
docker build -t servicedashboard . && docker run -it -d -p 8080:8080 servicedashboard:latest
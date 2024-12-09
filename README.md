# kong-service-dashboard
  Thought about project description for a second Services &amp; Users API in Go This repository provides a sample Go (Golang) application that implements a basic HTTP API for managing and retrieving users and services.

# Start Dependencies using docker compose
```sh
docker-compose down -v
```
Starting up the service and DB
```sh
docker-compose up -d
```

# Start containerized application
```sh
docker build -t test-server . && docker run -it test-server:latest
```

# How to run in local environment only
```sh
go run ./cmd
```


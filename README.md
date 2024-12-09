# kong-service-dashboard
  Thought about project description for a second Services &amp; Users API in Go This repository provides a sample Go (Golang) application that implements a basic HTTP API for managing and retrieving users and services.

# How to test with docker compose (WIP)
This starts up the database only, still working on the server
```sh
docker-compose up -d
```

# How to run in local environment only
```sh
go run ./cmd
```

# How to test this project without docker compose
```sh
docker build -t test-server . && docker run -it test-server:latest
```
Run `./app` to start the server
```sh
./app
```
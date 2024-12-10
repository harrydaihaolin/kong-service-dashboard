## Start Dependencies using Docker Compose

To build the dependent database, run:
```sh
docker-compose down -v
docker-compose up -d
```

## Start Containerized Application

To build and run the containerized application, execute:
```sh
./scripts/start_server.sh
```

## Run in Local Environment

To run the application locally without containerization, use:
```sh
go run ./cmd
```

## How to run unit tests

To run tests and generate coverage reports, execute:
```sh
./scripts/test_report.sh
```

To view coverage reports
```sh
open ./coverage/coverage.html
```

## How to run E2E Integration Test

To run integration tests
```sh
./scripts/e2e_tests.sh
```

## How to migrate

Ensure the `migrate` binary is installed:
```sh
migrate create -ext sql -dir ./migrations -seq create_users
```
Two SQL files will be generated under `/migrations`. Test the SQL plan first.

Use GORM in `./cmd/orm.go` to redefine the schema. Remove the migration file to verify if GORM matches the SQL plan.
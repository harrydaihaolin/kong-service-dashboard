# Kong Service Dashboard

This repository provides a sample Go (Golang) application that implements a basic HTTP API for managing and retrieving users and services.

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
Make sure `migrate` go binary is installed properly
```sh
migrate create -ext sql -dir ./migrations -seq create_users
```
Two SQL files are generated under `/migrations`. It is recommended to first test the SQL plan to ensure it works as expected. 

Then, use GORM in `./cmd/orm.go` to redefine the schema. 

Remove the migration file to verify if GORM matches the SQL plan. This provides an additional reference to ensure the database table is generated correctly.

## Configuration

The application can be configured using environment variables. The following variables are available:

- `PORT`: The port on which the application will run (default: `8080`).
- `DB_HOST`: The hostname of the database server (default: `localhost`).
- `DB_PORT`: The port of the database server (default: `5432`).
- `DB_USER`: The username for the database (default: `postgres`).
- `DB_PASSWORD`: The password for the database (default: `example`).
- `DB_NAME`: The name of the database (default: `postgres`).

### Example Configuration

To set up the environment variables, you can use the following commands:

```sh
export SERVICE_DASHBOARD_DB_PORT=5432
export SERVICE_DASHBOARD_DB_USER=postgres
export SERVICE_DASHBOARD_DB_PASSWORD=example
export SERVICE_DASHBOARD_DB_NAME=postgres
export SERVICE_DASHBOARD_DB_HOST=localhost
```

## API Endpoints

The application provides the following API endpoints:

### public endpoint
- `POST /v1/auth`: Retrieve the JWT token given username and password

### protected endpoints
- `GET /v1/users`: Retrieve a list of users.
- `POST /v1/users`: Create a new user.
- `GET /v1/users/{id}`: Retrieve a specific user by ID.
- `PUT /v1/users/{id}`: Update a specific user by ID.
- `DELETE /v1/users/{id}`: Delete a specific user by ID.
- `GET /v1/services`: Retrieve a list of services.
- `POST /v1/services`: Create a new service.
- `GET /v1/services/{id}`: Retrieve a specific service by ID.
- `PUT /v1/services/{id}`: Update a specific service by ID.
- `DELETE /v1/services/{id}`: Delete a specific service by ID.

## Example: User Authentication

Here is an example of how to authenticate a user using the `POST /v1/auth` endpoint:

### Request Payload

```json
{
  "username": "exampleUser",
  "password": "examplePassword"
}
```

### Example Request

```sh
curl -X POST "http://localhost:8080/v1/auth" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "password": "password"
  }'
```

### Example Response

```json
{
  "token": "generated.jwt.token"
}
```

If the request payload is invalid or the credentials are incorrect, the response will include an appropriate HTTP error status.
## Role-Based Access Control (RBAC) Middleware

The `RoleBasedMiddleware` is a middleware that checks if the request has a valid JWT token with the required role. It allows requests to whitelisted paths without a token. For other paths, it expects an Authorization header with a Bearer token. The token is then parsed and validated. If the token is valid and contains the required role, the request is allowed to proceed. Otherwise, an appropriate error response is returned.

### Example Usage
1. First of all, fetch the JWT token by calling /v1/auth endpoint with the credentials (see above examples)
2. JWT token contains the role information inferred by the user identity information
## Permissions

Permissions define the access control for different roles. It is structured as a nested map where the outer map's keys are role names (e.g., "admin", "user"), and the values are inner maps that specify the HTTP methods (e.g., "GET", "POST", "PUT", "DELETE") and whether each method is allowed (true) or not (false) for that role.

For example:
- The "admin" role has permissions to perform all HTTP methods: GET, POST, PUT, and DELETE.
- The "user" role is only allowed to perform the GET method.

```sh
curl -X GET "http://localhost:8080/v1/services" \
  -H "Authorization: Bearer <token>"
```

### Role-Based Access Control (RBAC)

RBAC is a method of regulating access to resources based on the roles of individual users within an organization. In this middleware, roles are extracted from the JWT token and checked against the required permissions for the requested HTTP method. If the role has the required permission, the request is allowed to proceed.

### JWT Token and Roles

Authorization is done via a JWT token obtained from the authentication process. The JWT token contains the role information of the registered user. Based on the role information, the system determines whether the user is an admin or a regular user.

### Example: User Authentication

Here is an example of how to authenticate a user using the `POST /v1/auth` endpoint:

#### Example Request

```sh
curl -X POST "http://localhost:8080/v1/authenticate" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user1",
    "password": "password"
  }'
```

#### Example Response

```json
{
  "token": "generated.jwt.token"
}
```

If the request payload is invalid or the credentials are incorrect, the response will include an appropriate HTTP error status.

## Example: Retrieve Services

Here is an example of how to retrieve services using the `GET /v1/services` endpoint:

```sh
curl -X GET "http://localhost:8080/v1/services?page=1&limit=10&sort_by=id&order=asc"
```

This request will retrieve a paginated list of services, sorted by `id` in ascending order. You can adjust the query parameters to customize the pagination and sorting:

- `page`: The page number to retrieve (default: `1`).
- `limit`: The number of services per page (default: `10`).
- `sort_by`: The field to sort by (default: `id`). Valid options are `id`, `service_name`, and `created_at`.
- `order`: The sort order (default: `asc`). Valid options are `asc` and `desc`.

You can also search for services by name or retrieve a specific service by ID:

```sh
# Search for services by name
curl -X GET "http://localhost:8080/v1/services?search_mode=true&name=example"

# Search specific service by name
curl -X GET "http://localhost:8080/v1/services?name=example"

# Retrieve a specific service by ID
curl -X GET "http://localhost:8080/v1/services?id=1"
```


## Example: Retrieve Users

Here is an example of how to retrieve users using the `GET /v1/users` endpoint:

```sh
curl -X GET "http://localhost:8080/v1/users"
```

This request will retrieve a list of all users. You can also search for a specific user by username:

```sh
# Search for a user by username
curl -X GET "http://localhost:8080/v1/users?username=example"
```

In the above example, replace `example` with the actual username you want to search for.
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

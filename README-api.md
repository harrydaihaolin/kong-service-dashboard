# API Usage Documentation

This document provides examples of how to interact with the API endpoints for user authentication, retrieving services, and managing users. Each section includes example requests, payloads, and responses to help you understand how to use the API effectively.

## Order of Contents
- [User Authentication](#example-user-authentication)
- [Retrieve Services](#example-retrieve-services)
- [Retrieve Users](#example-retrieve-users)
- [Create Users](#example-create-users)
- [Update Users](#example-update-users)
- [Delete Users](#example-delete-users)
- [Create Service Versions](#example-create-service-versions)
- [Update Service Versions](#example-update-service-versions)
- [Delete Service Versions](#example-delete-service-versions)

## Example: User Authentication

Here is an example of how to authenticate a user using the `POST /v1/auth` endpoint:

### Role-Based Access Control (RBAC) Middleware
The `RoleBasedMiddleware` checks for a valid JWT token with the required role. It allows requests to whitelisted paths without a token. For other paths, it expects an Authorization header with a Bearer token. If the token is valid and contains the required role, the request proceeds; otherwise, an error response is returned.

### Example Usage
1. First of all, fetch the JWT token by calling /v1/auth endpoint with the credentials (see above examples)
2. JWT token contains the role information inferred by the user identity information

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
    "token": "<your_jwt_token>"
}
```

If the request payload is invalid or the credentials are incorrect, the response will include an appropriate HTTP error status.

## Example: Retrieve Services

Here is an example of how to retrieve services using the `GET /v1/services` endpoint:

```sh
curl -X GET "http://localhost:8080/v1/services?page=1&limit=10&sort_by=id&order=asc" \
    -H "Authorization: Bearer <your_jwt_token>"
```

This request will retrieve a paginated list of services, sorted by `id` in ascending order. You can adjust the query parameters to customize the pagination and sorting:

### Query Parameters

- `page`: The page number to retrieve (default: `1`).
    - Example: `?page=2`
- `limit`: The number of services per page (default: `10`).
    - Example: `?limit=20`
- `sort_by`: The field to sort by (default: `id`). Valid options are `id`, `service_name`, and `created_at`.
    - Example: `?sort_by=service_name`
- `order`: The sort order (default: `asc`). Valid options are `asc` and `desc`.
    - Example: `?order=desc`
- `search_mode`: A flag to indicate if search mode is enabled (default: `false`).
    - Example: `?search_mode=true`
- `name`: The name of the service to search for.
    - Example: `?name=example-service`
- `id`: The ID of the service to retrieve.
    - Example: `?id=123`
- `load_version`: A flag to indicate if service versions should be loaded (default: `false`).
    - Example: `?load_version=true`

The function handles the following scenarios:
- Fetching a specific service by ID.
- Searching for services by name.
- Fetching a single service by name.
- Fetching paginated and sorted results.
You can also search for services by name or retrieve a specific service by ID:

```sh
# Search for services by name
curl -X GET "http://localhost:8080/v1/services?search_mode=true&name=example" \
    -H "Authorization: Bearer <your_jwt_token>"

# Search specific service by name
curl -X GET "http://localhost:8080/v1/services?name=example" \
    -H "Authorization: Bearer <your_jwt_token>"

# Retrieve a specific service by ID
curl -X GET "http://localhost:8080/v1/services?id=1" \
    -H "Authorization: Bearer <your_jwt_token>"

# Retrieve services with pagination and sorting
curl -X GET "http://localhost:8080/v1/services?page=2&limit=20&sort_by=service_name&order=desc" \
    -H "Authorization: Bearer <your_jwt_token>"

# Retrieve services with search mode enabled
curl -X GET "http://localhost:8080/v1/services?search_mode=true&name=example-service" \
    -H "Authorization: Bearer <your_jwt_token>"

# Retrieve a specific service by ID and load versions
curl -X GET "http://localhost:8080/v1/services?id=123&load_version=true" \
    -H "Authorization: Bearer <your_jwt_token>"
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

## Example: Create Users

Here is an example of how to create a new user using the `POST /v1/users` endpoint:

### Example Request

```sh
curl -X POST "http://localhost:8080/v1/users" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <your_jwt_token>" \
    -d '{
        "username": "newUser",
        "password": "newPassword",
        "role": "user",
        "user_profile": {
            "first_name": "New",
            "last_name": "User",
            "email": "newUser@example.com"
    }'
```

### Example Response

```json
{
  "ID": 3,
  "CreatedAt": "2024-12-10T08:37:22.798819462Z",
  "UpdatedAt": "2024-12-10T08:37:22.798819462Z",
  "DeletedAt": null,
  "service_name": "",
  "service_description": ""
}
```

## Example: Update Users
Here is an example of how to update an existing user using the `PUT /v1/users/{id}` endpoint:

### Example Request

```sh
curl -X PUT "http://localhost:8080/v1/users/3" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <your_jwt_token>" \
    -d '{
        "username": "updatedUser",
        "password": "updatedPassword",
        "role": "admin",
        "user_profile": {
            "first_name": "Updated",
            "last_name": "User",
            "email": "updatedUser@example.com"
    }'
```

### Example Response

```json
{
  "ID": 3,
  "CreatedAt": "2024-12-10T08:37:22.798819462Z",
  "UpdatedAt": "2024-12-11T09:45:22.798819462Z",
  "DeletedAt": null,
  "username": "updatedUser",
  "role": "admin",
  "user_profile": {
    "first_name": "Updated",
    "last_name": "User",
    "email": "updatedUser@example.com"
  }
}
```

## Example: Delete Users

Here is an example of how to delete an existing user using the `DELETE /v1/users/{id}` endpoint:

### Example Request

```sh
curl -X DELETE "http://localhost:8080/v1/users/3" \
    -H "Authorization: Bearer <your_jwt_token>"
```

### Example Response

```json
{
  "message": "User deleted successfully"
}
```

If the user ID does not exist, the response will include an appropriate HTTP error status.

## Example: Create Service Versions

Here is an example of how to create a new service version using the `POST /v1/service_versions` endpoint:

### Example Request

```sh
curl -X POST "http://localhost:8080/v1/service_versions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <your_jwt_token>" \
    -d '{
        "service_id": 1,
        "service_version_name": "v1.0.0",
        "service_version_description": "Initial release",
        "service_version_url": "http://example.com/v1.0.0"
    }'
```

### Example Response

```json
{
  "ID": 1,
  "ServiceID": 1,
  "ServiceVersionName": "v1.0.0",
  "ServiceVersionDescription": "Initial release",
  "ServiceVersionURL": "http://example.com/v1.0.0",
  "CreatedAt": "2024-12-10T08:37:22.798819462Z",
  "UpdatedAt": "2024-12-10T08:37:22.798819462Z"
}
```

## Example: Update Service Versions

Here is an example of how to update an existing service version using the `PUT /v1/service_versions` endpoint:

### Example Request

```sh
curl -X PUT "http://localhost:8080/v1/service_versions" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer <your_jwt_token>" \
    -d '{
        "id": 1,
        "service_version_name": "v1.0.1",
        "service_version_description": "Bug fixes",
        "service_version_url": "http://example.com/v1.0.1"
    }'
```

### Example Response

```json
{
  "ID": 1,
  "ServiceID": 1,
  "ServiceVersionName": "v1.0.1",
  "ServiceVersionDescription": "Bug fixes",
  "ServiceVersionURL": "http://example.com/v1.0.1",
  "CreatedAt": "2024-12-10T08:37:22.798819462Z",
  "UpdatedAt": "2024-12-11T09:45:22.798819462Z"
}
```

## Example: Delete Service Versions

Here is an example of how to delete an existing service version using the `DELETE /v1/service_versions` endpoint:

### Example Request

```sh
curl -X DELETE "http://localhost:8080/v1/service_versions/{id}" \
    -H "Authorization: Bearer <your_jwt_token>"
```

### Example Response

```json
{
  "message": "Service version deleted successfully"
}
```

If the service version ID does not exist, the response will include an appropriate HTTP error status.

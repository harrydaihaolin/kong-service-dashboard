# Kong Service Dashboard

This repository provides a sample Go (Golang) application that implements a basic HTTP API for managing and retrieving users and services.

## API Endpoints

The application provides the following API endpoints:

### Public Endpoint
- `POST /v1/auth`: Retrieve the JWT token given username and password

### Protected Endpoints
- `PUT /v1/services`: Update an existing service.
- `DELETE /v1/services`: Delete an existing service.
- `GET /v1/services`: Get an existing service.
- `POST /v1/services`: Create a new service.
- `GET /v1/users`: Retrieve a list of users.
- `POST /v1/users`: Create a new user.
- `PUT /v1/users`: Update an existing user.
- `DELETE /v1/users`: Delete an existing user.

## Links
- [API Documentation](README-api.md)
- [Developer Documentation](README-dev.md)
## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

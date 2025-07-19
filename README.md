# Go REST API III

This repository contains a sample Go REST API using **Chi** for routing and **Gorm** for database access. It demonstrates user authentication with JWT, product management and role based authorization. Swagger documentation is included for exploring the available endpoints.

## Project structure

- **cmd/** – entry points for the application. The `server` directory contains `main.go` which boots the web server.
- **configs/** – configuration loader that reads environment variables into a struct.
- **internal/** – private application packages such as entities, database implementations and HTTP handlers.
- **pkg/** – reusable helper packages that can be imported by other modules.

```
cmd/
  server/main.go        application entry point
configs/                configuration helpers
internal/
  dto/                  request/response DTOs
  entity/               domain entities
  infra/                infrastructure code (database, webserver)
  ...
pkg/                    shared utilities
```

## Getting started

Clone the repository and create a `.env` file based on `.env-example`:

```bash
cp .env-example .env
```

Edit the values if necessary. The main variables are:

- `DB_DRIVER` – database driver (e.g. `mysql`)
- `DB_HOST` – database host
- `DB_PORT` – database port
- `DB_USER` – database user
- `DB_PASSWORD` – database password
- `DB_NAME` – database name
- `WEB_SERVER_PORT` – port where the API will run
- `JWT_SECRET` – secret used to sign JWT tokens
- `JWT_EXPIRESIN` – token expiration time in seconds

### Running with Docker

The repository provides a `Dockerfile` and a `docker-compose.yaml` with MySQL and phpMyAdmin configured. Build and start the stack with:

```bash
docker-compose up --build
```

The API will be available on port `8080` (or the value from `WEB_SERVER_PORT`). The database container exposes MySQL on port `3307` by default.

### Running locally

To run without Docker make sure you have Go and a MySQL server installed. Create the database defined in `.env` and then execute:

```bash
go run cmd/server/main.go
```

When the server starts it performs the database migrations and seeds sample data.
If the user table is empty, three accounts are created for testing:

- **admin@example.com** / `1234` (role: `admin`)
- **manager@example.com** / `1234` (role: `manager`)
- **customer@example.com** / `1234` (role: `customer`)

### Swagger documentation

Swagger files are located in the `docs/` folder. If you modify the API you can regenerate them using [swag](https://github.com/swaggo/swag):

```bash
# install swag if not present
go install github.com/swaggo/swag/cmd/swag@latest

# generate documentation
swag init -g cmd/server/main.go -o docs
```

After running the server access `http://localhost:8080/docs/index.html` to explore the API.


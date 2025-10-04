# Nodehaus AI Platform

The Nodehaus AI Platform to finetune and evaluate models with custom training data.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

## Setup PostgreSQL database

Start a postgres shell:

```
$ sudo -u postgres psql

```

Run the following commands

```
CREATE DATABASE ai_platform;
CREATE USER nodehaus with password 'SecurePassword;
GRANT ALL PRIVILEGES ON DATABASE ai_platform TO nodehaus;
\connect ai_platform
GRANT ALL PRIVILEGES ON SCHEMA public TO nodehaus;
```

Create `.env` file (see below for Docker) and run migrations:

```
$ ./run_migrations.sh
```

Add a test user:

```
$ sudo -u postgres psql

```

```
INSERT INTO users (email, password) VALUES ('test@example.com', '$2a$12$k2WRsfc9868pKseoXaGAf.YdtXrp8uXumJiWoTxq1UxBWQ5m0df96');
```

## Docker Deployment

The application includes a Dockerfile for containerized deployment. Here's how to deploy it on a server with automatic
restart capabilities:

### Building the Docker Image

```bash
# Build the Docker image
docker build --platform linux/amd64 -t ai-platform .

# Or with a specific tag
docker build -t ai-platform:latest .
```

```bash
docker tag ai-platform pebouda/ai-platform
docker push pebouda/ai-platform
```

### Running with Docker

#### Using Environment File

Create a `.env` file:

```env
PORT=8081
API_BASE_URL=https://ai.peterbouda.eu
BLUEPRINT_DB_HOST=host.docker.internal
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=ai_platform
BLUEPRINT_DB_USERNAME=nodehaus
BLUEPRINT_DB_PASSWORD=your-password
BLUEPRINT_DB_SCHEMA=public
JWT_SECRET_KEY=your-jwt-secret
AWS_ENDPOINT_URL=http://s3.peterbouda.eu:3900
AWS_DEFAULT_REGION=garage
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
APP_S3_BUCKET=nodehaus
APP_EXTERNAL_API_KEY=VerySecureKey
RUNPOD_API_KEY=
RUNPOD_POD_ID_FINETUNE=
RUNPOD_POD_ID_OLLAMA=
GIN_MODE=release
```

Then run:

```bash
docker run -d \
  --name ai-platform \
  --restart always \
  -p 8081:8081 \
  --env-file .env \
  --add-host=host.docker.internal:host-gateway \
  pebouda/ai-platform
```

Add `--add-host=host.docker.internal:host-gateway` if you want to access postgres on the host.

### Restart Policies

-   `--restart unless-stopped`: Restart unless manually stopped
-   `--restart always`: Always restart (even after system reboot)
-   `--restart on-failure`: Restart only on failure
-   `--restart no`: Never restart automatically (default)

# Microservices Assignment

A small microservices-based backend system built using Go, PostgreSQL, NATS JetStream, JWT authentication, Docker, and Swagger.


## Architecture

![Microservices Architecture](docs/architecture.png)

## Communication

The API Gateway communicates with the User Service using HTTP.

The User Service and Notification Service **do not communicate using REST APIs or WebSockets**.

Instead, the User Service publishes a `user.created` event to NATS JetStream. The Notification Service asynchronously consumes this event.

## Services

### 1. API Gateway

Port: `8080`

Responsibilities:

* Single entry point for clients
* Routes requests to backend services
* JWT authentication middleware
* Swagger API documentation

Endpoints:

```text
GET  /health
POST /api/auth/login
POST /api/users
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

### 2. User Service

Port: `8081`

Responsibilities:

* User registration
* User login
* Password hashing using bcrypt
* PostgreSQL persistence
* Publishing `user.created` events

Endpoints:

```text
GET  /health
POST /users
POST /login
```

### 3. Notification Service

Port: `8082`

Responsibilities:

* Consume `user.created` events from NATS JetStream
* Create welcome notifications
* Store notifications in PostgreSQL
* Acknowledge successfully processed messages

The Notification Service does not expose a REST API for communication with the User Service.

## Technology Stack

* Go
* PostgreSQL (Neon)
* NATS JetStream
* JWT
* bcrypt
* Docker & Docker Compose
* Swagger / OpenAPI

## Event Flow

When a new user is created:

```text
Client
  |
  v
API Gateway
  |
  v
User Service
  |
  | Save user
  v
PostgreSQL

User Service
  |
  | Publish user.created
  v
NATS JetStream
  |
  | Async delivery
  v
Notification Service
  |
  | Save notification
  v
PostgreSQL
```

Example event:

```json
{
  "event_id": "unique-event-id",
  "user_id": "user-id",
  "name": "Ritik",
  "email": "user@example.com"
}
```

## Reliability

NATS JetStream provides persistent messaging.

The Notification Service uses:

* Manual message acknowledgement
* ACK after successful database processing
* Message redelivery when processing fails
* Maximum delivery attempts
* Durable consumer
* Unique `event_id` to prevent duplicate notifications

This allows failed messages to be retried instead of being silently lost.

## Authentication

JWT-based authentication is used.

### Login

```text
POST /api/auth/login
```

Request:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

The response contains a JWT token.

Protected endpoints require:

```text
Authorization: Bearer <JWT_TOKEN>
```

The API Gateway validates the JWT before forwarding protected requests.

## Environment Variables

Sensitive configuration is stored using environment variables.

Create a `.env` file in the project root:

```env
DATABASE_URL=your-neon-database-url
NATS_URL=nats://localhost:4222
JWT_SECRET=your-secret-key
```

Do not commit `.env` to GitHub.

A `.env.example` file is provided as a template.

## Running with Docker

Make sure Docker Desktop is running.

From the project root:

```bash
docker compose build
```

Start all services:

```bash
docker compose up -d
```

Check running containers:

```bash
docker compose ps
```

Expected services:

```text
gateway
user-service
notification-service
nats
```

View logs:

```bash
docker compose logs -f user-service
```

```bash
docker compose logs -f notification-service
```

```bash
docker compose logs -f gateway
```

Stop services:

```bash
docker compose down
```

## Testing

### Gateway Health

```bash
curl http://localhost:8080/health
```

### Login

Use Swagger:

```text
http://localhost:8080/swagger/index.html
```

First create a user, then login to receive a JWT token.

### Protected User Creation

Send the JWT using:

```text
Authorization: Bearer <JWT_TOKEN>
```

Then call:

```text
POST /api/users
```

After successful user creation, the User Service publishes a `user.created` event.

The Notification Service consumes the event and stores a notification.

## Project Structure

```text
microservices-assignment/
│
├── gateway/
│   ├── main.go
│   ├── auth.go
│   ├── proxy.go
│   ├── docs/
│   ├── .dockerignore
│   └── go.mod
│
├── user-service/
│   ├── main.go
│   ├── auth.go
│   ├── handlers/
│   ├── models/
│   ├── repositries/
│   ├── events/
│   ├── .dockerignore
│   └── go.mod
│
├── notification-service/
│   ├── main.go
│   ├── models/
│   ├── repositries/
│   ├── events/
│   ├── .dockerignore
│   └── go.mod
│
├── docker-compose.yml
├── .env.example
├── .gitignore
└── README.md
```

## Security Considerations

* Passwords are stored as bcrypt hashes rather than plaintext.
* JWT is used for authentication.
* JWT secret is stored in an environment variable.
* `.env` files are excluded from Docker build contexts.
* `.env` is excluded from Git.
* Protected API routes require JWT authentication.
* User Service and Notification Service communicate asynchronously through NATS JetStream.

## Future Improvements

Possible production improvements include:

* HTTPS/TLS between services
* Refresh tokens
* Rate limiting
* Centralized logging
* Distributed tracing
* Health/readiness checks
* Kubernetes deployment

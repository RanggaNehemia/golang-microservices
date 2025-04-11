# Golang Microservices: Auth, Data, Trade (OAuth2-compliant)

A scalable microservices architecture built in Go.\
Utilizes OAuth2-compliant authentication and inter-service communication using JWTs and PostgreSQL.

---

## Microservices Overview

| Service         | Port | Description                                                 |
| --------------- | ---- | ----------------------------------------------------------- |
| `auth-service`  | 8080 | Handles user and machine authentication with token issuance |
| `data-service`  | 8081 | Periodically generates and stores random pricing data       |
| `trade-service` | 8082 | Allows users to place validated trades based on price data  |

---

Port 8080, 8081, and 8082 are the default ports for each services.

## Getting Started

### Prerequisites

- [Go 1.20+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [pgAdmin 4](https://www.pgadmin.org/) (optional GUI)
- [`golang-migrate`](https://github.com/golang-migrate/migrate) (optional)
- Git, curl, Postman or other API testing tools

---

## Project Structure

```
golang-microservices/
├── auth-service/
├── data-service/
├── trade-service/
└── README.md
```

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/RanggaNehemia/golang-microservices.git
cd golang-microservices
```

### 2. Create PostgreSQL Databases

Create three databases:

- Authentication Database
- Data Database
- Trade Database

### 3. Configure `.env` for Each Service

Each microservice needs a `.env` file in root folder.
These are the example of the env:

#### `auth-service/.env`

```env
PORT=8080
GORM_DATABASE_URL="host=<host> user=<user> password=<password> dbname=<database> port=<port> sslmode=disable TimeZone=UTC"
PGX_DATABASE_URL="postgres://<user>:<password>@<host>:<port>/<database>?sslmode=disable&timezone=UTC"
SECRET_KEY=<secret_key>
```

#### `data-service/.env`

```env
PORT=8081
DATABASE_URL="host=localhost user=<user> password=<password> dbname=<data_database_name> port=<port> sslmode=disable TimeZone=UTC"
SECRET_KEY=<secret_key>
TRADE_SERVICE_CLIENT_ID="trade-service"
AUTH_URL="<auth-service-url>"
```

#### `trade-service/.env`

```env
PORT=8082
DATABASE_URL="host=localhost user=<user> password=<password> dbname=<trade_database_name> port=<port> sslmode=disable TimeZone=UTC"
SECRET_KEY=<secret_key>
DATA_SERVICE_URL="<data_service_url:port>"
TRADE_SERVICE_TOKEN="<trade service token>"
TRADE_SERVICE_CLIENT_ID="<trade-service client id>"
TRADE_SERVICE_CLIENT_SECRET="<trade-service client secret>"
WEB_CLIENT_ID="<web-client id>"
AUTH_URL="<auth-service-url>"
```

### 4. Install Dependencies

For each service in the root folder:

```bash
go mod tidy
```

---

## Running the Services

Open a terminal for each service and run:

```bash
go run main.go
```

---

## Authentication Overview

### User Login Flow

- `POST /auth/register` - Register new user
- `POST /oauth/token` - Login and get JWT token

Use this token to authenticate user actions like placing trades.

### Machine Login Flow

Call the API

- `POST /oauth/token` - Authenticate service and get machine token

Use this for internal communication (e.g., trade - data).

---

## Auth Service Endpoints

| Endpoint            | Method | Description                   |
| ------------------- | ------ | ----------------------------- |
| `/health`           | GET    | Health Check                  |
| `/auth/register`    | POST   | Register New User             |
| `/oauth/token`      | POST   | Login and get JWT token       |
| `/oauth/authorize`  | POST   | Authorize the token given     |
| `/oauth/introspect` | POST   | Introspect the token given    |
| `/oauth/revoke`     | POST   | revoke the token given        |
| `/auth/me`          | GET    | Retrieve current user details |

---

## Data Service Endpoints

| Endpoint       | Method | Description                        |
| -------------- | ------ | ---------------------------------- |
| `/data/latest` | GET    | Returns the most recent price      |
| `/data/lowest` | GET    | Returns the lowest price in 24 hrs |

> Requires token with trade-service audience

---

## Trade Service Endpoints

| Endpoint       | Method | Description       |
| -------------- | ------ | ----------------- |
| `/trade/place` | POST   | Place a new trade |

> Requires a token with web-service audience

Trades cannot be placed below 50% of the lowest price in the last 24 hours.

---

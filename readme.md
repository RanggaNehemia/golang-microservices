# Golang Microservices: Auth, Data, Trade (OAuth2-compliant)

A scalable microservices architecture built in Go.\
Utilizes OAuth2-compliant authentication and inter-service communication using JWTs and PostgreSQL.

---

## Microservices Overview

| Service         | Port | Description |
|------------------|------|-------------|
| `auth-service`   | 8080 | Handles user and machine authentication with token issuance |
| `data-service`   | 8081 | Periodically generates and stores random pricing data |
| `trade-service`  | 8082 | Allows users to place validated trades based on price data |
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
DATABASE_URL="host=localhost user=<user> password=<password> dbname=<auth_database_name> port=5432 sslmode=disable TimeZone=UTC"
SECRET_KEY=<secret_key>
```

#### `data-service/.env`

```env
PORT=8081
DATABASE_URL="host=localhost user=<user> password=<password> dbname=<data_database_name> port=5432 sslmode=disable TimeZone=UTC"
SECRET_KEY=<secret_key>
```

#### `trade-service/.env`

```env
PORT=8082
DATABASE_URL="host=localhost user=<user> password=<password> dbname=<trade_database_name> port=5432 sslmode=disable TimeZone=UTC"
SECRET_KEY=<secret_key>
DATA_SERVICE_URL="<data_service_url:port>"
TRADE_SERVICE_TOKEN="<trade service token>"

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
- `POST /auth/login` - Login and get JWT token

Use this token to authenticate user actions like placing trades.

### Machine Login Flow

Add the service into Client table in Authentication Database

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO clients (name, secret, created_at, updated_at)
VALUES (
  '<service name>',
  crypt('<secret key>', gen_salt('bf')),
  NOW(),
  NOW()
);
```

Call the API

- `POST /auth/client/token` - Authenticate service and get machine token

Use this for internal communication (e.g., trade - data).

---

## Auth Service Endpoints

| Endpoint            | Method | Description                                 |
|---------------------|--------|---------------------------------------------|
| `/auth/register`    | POST   | Register New User                           |
| `/auth/login`       | POST   | Login and get JWT token                     |
| `/auth/client/token`| POST   | Authenticate service and get machine token  |
| `/auth/me`       	  | GET    | Retrieve current user details               |

---

## Data Service Endpoints

| Endpoint            | Method | Description                        |
|---------------------|--------|------------------------------------|
| `/data/latest`      | GET    | Returns the most recent price      |
| `/data/lowest`      | GET    | Returns the lowest price in 24 hrs |

> Requires a **machine token**

---

## Trade Service Endpoints

| Endpoint          | Method | Description                                  |
|-------------------|--------|----------------------------------------------|
| `/trades`         | POST   | Place a new trade                            |

> Requires a **user token**

Trades cannot be placed below 50% of the lowest price in the last 24 hours.

---
# Golang High Concurrency

## Prerequisites
- Visual Studio Code (for easier readibility code, add Go extension as well via marketplace)
- PostgreSQL
- Docker Desktop

## How To Run
- create database and tables using schema.sql (make sure choosing right active database - ubersnap)
- docker compose up --build

## How To Test (via Postman)
- POST /api/coupons , body json {"name":"namehere","amount":(amount here)}
- POST /api/coupons/claim , body json {"user_id":"user id here","coupon_name":"coupon name here"}
- GET  /api/coupons/{name} , empty body

## Scenario
- User can only claim same coupon ONCE
- Each coupon has their LIMITED amount

## Architecture Notes : Database Design and Locking Strategy
### Using Optimistic Locking
- Uses version column
- Detects conflicts safely
- Scales horizontally

## Services
- App: http://localhost:8080
- PostgreSQL: localhost:5432

Endpoint:
- POST /api/coupons
- POST /api/coupons/claim
- GET  /api/coupons/{name}

- Language   : Golang
- Framework  : Echo
- Database   : Postgres
- Dockerized and Kafka (on progress)

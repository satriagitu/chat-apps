# Chat Application

## Description

Chat Application is a chat platform that allows users to send messages, communicate within conversations, and manage notifications and files. The application is built using Go with Gin as the HTTP framework, GORM for ORM, and RabbitMQ for background task queues.

## Fitur

- **Pengguna**: Manage users with endpoints for creating and retrieving user information.
- **Notifikasi**: Send and retrieve notifications for users.
- **Pesan**: Send and retrieve messages within conversations.
- **File**: Upload and retrieve files.
- **Percakapan**: Create and retrieve conversations.
- **Antrian Tugas**: Process broadcast notifications using RabbitMQ.

## Prasyarat
Install:
- Go
- PostgreSQL
- RabbitMQ
- Docker

## Create Database
```sql
CREATE DATABASE chat;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE conversation_participants (
    conversation_id INT REFERENCES conversations(id),
    user_id INT REFERENCES users(id),
    PRIMARY KEY (conversation_id, user_id)
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    conversation_id INT REFERENCES conversations(id),
    sender_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    file_url TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'queued',
    queued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

```

## Clone Repository
```bash
git clone https://github.com/username/chat-apps.git
cd chat-apps
```
## Run Unit Test
```bash
go test ./...
```
## Build Aplikasi
```bash
docker compose up --build
```

## Test
Lakukan test endpoint sesuai dokumentasi OpenAPI / Swagger.

## How to test (Unit Test)
```bash
-- test in spesific files
go test -v internal/controller/article_test.go 

-- test in spesific files and specific functions
go test -v internal/service/article_test.go -run ^TestArtik
elService_GetArticleList_Error$

-- percentage test -- 
-- create converage test all file
go test -coverprofile=coverage.out ./...

-- create covereage spesific file
go test -coverprofile=coverage.out ./internal/controller

-- view percentage coverage test
go tool cover -func=coverage.out

-- view percentage coverage test - spesific file
go tool cover -func=coverage.out | grep "internal/repository/file.go"

-- view percentage coverage test in html
go tool cover -html=coverage.out

--- Generate Mock with Mockary ---
mockery --dir=internal/repository --output=internal/repository/mocks --all
mockery --dir=internal/service --output=internal/service/mocks --all
mockery --name=ArticleRepository --dir=internal/repository --output=internal/repository/mocks
mockery --name=ArticleService --dir=internal/service --output=internal/service/mocks

```
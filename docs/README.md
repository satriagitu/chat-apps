# Chat Application

## Deskripsi

Chat Application adalah aplikasi chat yang memungkinkan pengguna untuk mengirim pesan, berkomunikasi dalam percakapan, dan mengelola notifikasi serta file. Aplikasi ini menggunakan Go dengan Gin sebagai framework HTTP, GORM untuk ORM, dan RabbitMQ untuk antrian tugas latar belakang.

## Fitur

- **Pengguna**: Mengelola pengguna dengan endpoint untuk membuat dan mendapatkan informasi pengguna.
- **Notifikasi**: Mengirim dan mengambil notifikasi untuk pengguna.
- **Pesan**: Mengirim dan mengambil pesan dalam percakapan.
- **File**: Mengunggah dan mendapatkan file.
- **Percakapan**: Membuat dan mendapatkan percakapan.
- **Antrian Tugas**: Memproses notifikasi siaran menggunakan RabbitMQ.

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
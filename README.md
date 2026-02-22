# 🛒 POS APP — Point of Sale REST API

A robust, production-ready **Point of Sale (POS)** backend API built with **Go**, designed for managing products, inventory, orders, reservations, and users in a modern retail or restaurant environment.

---

## 📖 Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Environment Configuration](#environment-configuration)
  - [Running the Application](#running-the-application)
- [API Modules](#api-modules)
- [Contributing](#contributing)

---

## Overview

**POS APP** is a scalable RESTful API service for point-of-sale operations. It provides a complete backend solution covering authentication, product & category management, inventory tracking, order processing, table reservations, user management, and real-time notifications — all built on a clean layered architecture.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | [Go](https://golang.org/) |
| Web Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| Database | [PostgreSQL](https://www.postgresql.org/) |
| Logger | [Zap](https://github.com/uber-go/zap) + [Lumberjack](https://github.com/natefinch/lumberjack) |
| Config | [Viper](https://github.com/spf13/viper) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Email Service | SMTP (configurable) |
| UUID | [google/uuid](https://github.com/google/uuid) |
| Decimal | [shopspring/decimal](https://github.com/shopspring/decimal) |

---

## Project Structure

```
project-POS-APP/
├── cmd/                  # Application entry point & server setup
├── internal/
│   ├── adaptor/          # HTTP handlers (controllers)
│   ├── data/             # Database migrations, seeders & repositories
│   │   └── repository/   # Data access layer
│   ├── dto/              # Data Transfer Objects (request & response)
│   ├── usecase/          # Business logic layer
│   └── wire/             # Dependency injection wiring
├── pkg/
│   ├── database/         # Database connection & pool setup
│   └── utils/            # Utilities (config reader, logger, helpers)
├── logs/                 # Application log files (auto-generated)
├── .env.example          # Environment variable template
├── go.mod
└── main.go
```

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) `>= 1.21`
- [PostgreSQL](https://www.postgresql.org/) `>= 14`
- An SMTP email account (for notification features)

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/bayuf/project-POS-APP-golang-string-team.git
   cd project-POS-APP-golang-string-team
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Copy and configure the environment file:**

   ```bash
   cp .env.example .env
   ```

### Environment Configuration

Edit `.env` with your actual values:

```env
# App
APP_NAME=POS-APP
PORT=8080
DEBUG=true
GIN_MODE=debug        # debug | release
LIMIT=3
PATH_LOGGING=./logs/app.log

# Database Connection
DATABASE_NAME=pos_db
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=yourpassword
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_SSL_MODE=disable

# Database Pool
DATABASE_MAX_CONN=10
DATABASE_MAX_IDLE_CONN=5
DATABASE_MAX_OPEN_CONN=10

# Database Migration & Seeder
DATABASE_MIGRATE=true     # Run auto-migration on startup
DATABASE_MIGRATE_SEEDER=false  # Seed initial data

# Email (SMTP)
EMAIL_FROM=your@email.com
EMAIL_PASSWORD=yourpassword
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
```

> **Note:** Set `DATABASE_MIGRATE=true` on **first run** to auto-create tables. Set `DATABASE_MIGRATE_SEEDER=true` to populate initial data. It is recommended to set both to `false` after the first run.

### Running the Application

```bash
go run main.go
```

The server will start on the port specified in your `.env` (default: `8080`).

For production, build the binary first:

```bash
go build -o pos-app .
./pos-app
```

---

## API Modules

| Module | Description |
|---|---|
| **Auth** | User registration, login, JWT-based authentication |
| **Users** | User profile management, role-based access |
| **Products** | Product CRUD, pricing, and details |
| **Categories** | Product category management |
| **Inventory** | Stock tracking and adjustments |
| **Orders** | Order creation, processing, and history |
| **Reservations** | Table or slot reservation management |
| **Notifications** | In-app notification system |
| **Email** | Transactional email delivery via SMTP |

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Commit your changes: `git commit -m "feat: add your feature"`
4. Push to your branch: `git push origin feature/your-feature-name`
5. Open a Pull Request

---

> Built with ❤️ by the **String Team**

# 🍽️ Restro API

A production-style REST API for a restaurant / food-ordering platform, built with **Go**, **Gin**, **GORM** and **MySQL**. It covers the complete customer journey — register, browse the menu, build a cart, place an order, and pay through **Razorpay** — with role-based access control for admins.

[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-1.12-00B386)](https://gin-gonic.com)
[![GORM](https://img.shields.io/badge/GORM-1.31-red)](https://gorm.io)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)](https://www.mysql.com)

---

## ✨ Features

- **JWT authentication** — stateless login with `Bearer` tokens and password hashing via bcrypt.
- **Role-based authorization** — separate `user` and `admin` capabilities enforced by middleware.
- **Menu management** — public product catalogue, admin-only create/update/delete.
- **Shopping cart** — one active cart per user with price snapshots at the time items are added.
- **Orders** — orders are created from the cart; users see only their own, admins see and update all.
- **Razorpay payments** — server-side order creation and HMAC signature verification.
- **Clean layered architecture** — controllers → services → repositories, wired through interfaces for easy testing.
- **Auto migrations** — schema is created/updated on startup via GORM `AutoMigrate`.
- **Unit tests** — table-driven tests for services, controllers and middleware.

---

## 🧱 Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP framework | Gin |
| ORM | GORM (MySQL driver) |
| Database | MySQL 8 |
| Auth | JWT (`golang-jwt/v5`) + bcrypt (`golang.org/x/crypto`) |
| Payments | Razorpay |
| Config | `.env` via `joho/godotenv` |

---

## 📁 Project Structure

```
Restro/
├── cmd/
│   └── main.go                 # Entry point: env, config, DB, routes, server
├── common/
│   ├── common.go               # Shared API response helpers
│   └── database/
│       ├── dbconnection.go     # MySQL connection
│       └── migrations.go       # GORM AutoMigrate
├── config/
│   ├── config.go               # Config loader + validation
│   ├── database.go             # DB config & DSN
│   ├── server.go               # Port, mode, env, JWT secret
│   ├── razorpay.go             # Payment gateway credentials
│   └── log.go                  # Logging config
└── internal/
    ├── controller/             # HTTP handlers (request/response)
    ├── middleware/             # JWT auth, admin guard
    ├── models/                 # GORM models & request/response DTOs
    ├── repository/             # Database access layer
    ├── routes/                 # Route registration & dependency wiring
    └── services/               # Business logic (incl. Razorpay client)
```

**Request flow:** `Route → Middleware → Controller → Service → Repository → MySQL`

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- A running **MySQL 8** instance with a database created
- A [Razorpay](https://razorpay.com/) account (optional — only needed for payment endpoints)

### 1. Clone the repository

```bash
git clone https://github.com/mannanpatel/Restro.git
cd Restro
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
# Server
SERVER_PORT=8085
SERVER_MODE=debug           # debug | release | test
SERVER_ENV=development      # development | staging | production
JWT_SECRET=replace-with-a-long-random-secret

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=restro

# Logging
LOG_LEVEL=debug             # debug | info | warn | error
LOG_FORMAT=text             # text | json
EXTENDED_LOGGING=false

# Razorpay (optional — required only for payment endpoints)
RAZORPAY_KEY_ID=
RAZORPAY_KEY_SECRET=
```

> ⚠️ `.env` is git-ignored. Never commit real credentials.

### 4. Run the server

```bash
go run ./cmd
```

The API starts on `http://localhost:8085`. Verify it:

```bash
curl http://localhost:8085/health
# {"status":"Restro API is healthy!"}
```

### 5. Run the tests

```bash
go test ./...
```

---

## ⚙️ Environment Variables

| Variable | Required | Description |
|---|---|---|
| `SERVER_PORT` | ✅ | Port the HTTP server listens on |
| `SERVER_MODE` | ✅ | Gin mode: `debug`, `release`, or `test` |
| `SERVER_ENV` | ✅ | Deployment environment name |
| `JWT_SECRET` | ✅ | Secret used to sign and verify JWTs |
| `DB_HOST` | ✅ | MySQL host |
| `DB_PORT` | ✅ | MySQL port |
| `DB_USER` | ✅ | MySQL user |
| `DB_PASSWORD` | ✅ | MySQL password |
| `DB_NAME` | ✅ | MySQL database name |
| `LOG_LEVEL` | ✅ | Logging verbosity |
| `LOG_FORMAT` | ✅ | Log output format |
| `EXTENDED_LOGGING` | ❌ | Set to `true` for verbose logs |
| `RAZORPAY_KEY_ID` | ❌ | Razorpay key ID (needed for payments) |
| `RAZORPAY_KEY_SECRET` | ❌ | Razorpay key secret (needed for payments) |

Startup fails fast with a clear error if any required variable is missing.

---

## 📡 API Reference

Base URL: `http://localhost:8085`

Protected endpoints require an `Authorization` header:

```
Authorization: Bearer <token>
```

### Health

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/health` | Public | Service health check |

### Authentication

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Create a new account |
| `POST` | `/auth/login` | Public | Authenticate and receive a JWT |

### Products (Menu)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/products` | Public | List all products |
| `GET` | `/products/:id` | Public | Get a single product |
| `POST` | `/admin/products` | Admin | Create a product |
| `PATCH` | `/admin/products/:id` | Admin | Update a product |
| `DELETE` | `/admin/products/:id` | Admin | Delete a product |

### Users

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/users/me` | User | Get the current profile |
| `PATCH` | `/users/me` | User | Update the current profile |
| `DELETE` | `/users/me` | User | Delete the current account |
| `GET` | `/admin/users` | Admin | List all users |
| `PUT` | `/admin/users/:id/role` | Admin | Change a user's role (`admin` / `user`) |

### Cart

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/cart` | User | Get the current cart |
| `POST` | `/cart/items` | User | Add an item to the cart |
| `PATCH` | `/cart/items/:itemID` | User | Update an item's quantity |
| `DELETE` | `/cart/items/:itemID` | User | Remove an item |
| `DELETE` | `/cart` | User | Empty the cart |

> `POST /cart/add` is kept as an alias of `POST /cart/items` for backward compatibility.

### Orders

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/orders` | User | Create an order from the current cart |
| `GET` | `/orders` | User | List the current user's orders |
| `GET` | `/orders/:id` | User | Get one of the current user's orders |
| `GET` | `/admin/orders` | Admin | List all orders |
| `PATCH` | `/admin/orders/:id/status` | Admin | Update an order's status |

Valid order statuses: `pending`, `confirmed`, `preparing`, `ready`, `completed`, `cancelled`.

### Payments

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/orders/:id/payment` | User | Create a Razorpay order for a saved order |
| `POST` | `/payments/verify` | User | Verify the Razorpay payment signature |

---

## 🔐 Example Usage

**Register**

```bash
curl -X POST http://localhost:8085/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Manan",
    "last_name": "Patel",
    "email": "manan@example.com",
    "password": "Str0ngPassw0rd",
    "address": "Ahmedabad",
    "age": 24
  }'
```

**Login**

```bash
curl -X POST http://localhost:8085/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"manan@example.com","password":"Str0ngPassw0rd"}'
```

```json
{
  "user_id": "1",
  "email": "manan@example.com",
  "role": "user",
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Add to cart**

```bash
curl -X POST http://localhost:8085/cart/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 2}'
```

**Place an order**

```bash
curl -X POST http://localhost:8085/orders \
  -H "Authorization: Bearer $TOKEN"
```

---

## 💳 Payment Flow

1. The client builds a cart and calls `POST /orders` to persist an order.
2. The client calls `POST /orders/:id/payment`. The server looks up the saved order, derives the amount from it (never from client input), creates a Razorpay order and returns the checkout details:

   ```json
   {
     "key_id": "rzp_test_xxx",
     "razorpay_order_id": "order_xxx",
     "amount": 45000,
     "currency": "INR",
     "order_id": 12
   }
   ```

3. The client completes checkout with the Razorpay SDK.
4. The client posts the gateway response to `POST /payments/verify`:

   ```json
   {
     "razorpay_payment_id": "pay_xxx",
     "razorpay_order_id": "order_xxx",
     "razorpay_signature": "9ef4dffb..."
   }
   ```

5. The server validates the HMAC-SHA256 signature with `RAZORPAY_KEY_SECRET` and marks the order as paid only if it matches.

If the Razorpay credentials are not set, the server still starts, but payment endpoints respond with `503 Service Unavailable`.

---

## 🗄️ Database Migrations

The API runs GORM `AutoMigrate` on startup after connecting to MySQL, creating or updating the schema for **users, products, carts, cart items, orders, order items and payments**.

`AutoMigrate` does **not** drop unused columns or tables. Review schema changes before deploying to production, and use explicit, reviewed migrations for destructive changes or data transformations.

---

## 🧪 Testing

```bash
go test ./...                       # all tests
go test ./internal/services/...     # a single package
go test -cover ./...                # with coverage
```

Services and repositories are defined behind interfaces, so handlers and business logic are tested with lightweight fakes — no database required.

---

## 🤝 Contributing

1. Fork the repository and create a branch: `git checkout -b feat/your-feature`
2. Follow the existing layered structure (controller → service → repository).
3. Run `go fmt ./...`, `go vet ./...` and `go test ./...` before pushing.
4. Open a pull request describing the change and how it was tested.

---

## 📄 License

This project is currently unlicensed. Add a `LICENSE` file to define how others may use it.
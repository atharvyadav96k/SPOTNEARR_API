# Spotnearr API — Architecture & Folder Structure

## Overview

Spotnearr is a microservices backend built in Go. The monorepo uses a **Go workspace** (`go.work`) so all services and shared packages can be developed and compiled together without publishing to a registry.

### Services

| Service | Port | Database | Module |
|---|---|---|---|
| `user-service` | 8080 | `spotnearr_user` | `app/cmd` |
| `vendor-service` | 8081 | `spotnearr_vendor` | `vendor-svc` |
| `search-service` | 8082 | `spotnearr_search` | `search-svc` |

### Cross-Service Communication Rule

Each service owns and manages its own database exclusively. Services never read from or write to another service's database directly. Synchronous requests use HTTP; asynchronous events are published over **RabbitMQ** (durable topic exchange `spotnearr.events`). Data changes are never propagated via shared tables.

---

## Repository Root

```
SPOTNEARR_API/
├── go.work                   # Go workspace — ties all modules together
├── go.work.sum
├── docker-compose.yml        # Production compose: postgres, redis, all 3 services
├── docker-compose.test.yml   # Test compose: ephemeral postgres, per-service profiles
├── Dockerfile.app            # User service image
├── Dockerfile.vendor         # Vendor service image
├── Dockerfile.search         # Search service image
├── Dockerfile.vendor.test    # Vendor integration test image (compile at build time)
├── Dockerfile.search.test    # Search integration test image (compile at build time)
├── ARCHITECTURE.md           # This file
├── API_DOCS.md               # HTTP endpoint reference
├── readme.md
├── index.html                # Landing page
│
├── app/                      # User service
├── services/
│   ├── vendor/               # Vendor service
│   └── search/               # Search service
├── database/                 # Shared GORM models + migrations (imported by services)
├── pkg/                      # Shared utilities (imported by all services)
├── docker/                   # Docker init scripts
└── dummy/                    # Local development seed scripts (Node.js)
```

---

## Go Workspace Modules

```
go.work
├── ./app/cmd          → (no explicit module name, entry point for user-service)
├── ./services/vendor  → github.com/atharvyadav96k/spotnearr/vendor-svc
├── ./services/search  → github.com/atharvyadav96k/spotnearr/search-svc
├── ./database         → (shared DB models, imported by services)
└── ./pkg              → github.com/atharvyadav96k/spotnearr/pkg
```

Each service has its own `go.mod` and uses `replace` directives to point at the local `pkg` module.

---

## `pkg/` — Shared Utilities

```
pkg/
├── go.mod
├── dtos/
│   ├── auth.go          # Auth request/response types
│   ├── business.go      # Business DTOs
│   ├── category.go      # Category DTOs
│   ├── inventory.go     # Inventory DTOs
│   ├── product.go       # Product DTOs
│   ├── review.go        # Review DTOs
│   └── search.go        # Search DTOs
├── httputil/
│   └── response.go      # Standard JSON response helpers
├── jwtutil/
│   ├── jwt.go           # JWT sign/parse
│   └── models.go        # Claims, UserRole type
├── middleware/
│   ├── auth.go          # JWT auth middleware
│   ├── business.go      # Business-only middleware
│   ├── captcha.go       # Cloudflare Turnstile verification
│   ├── cors.go          # CORS headers
│   ├── ratelimit.go     # Per-IP rate limiter middleware
│   └── session.go       # Session helpers
├── mq/
│   ├── conn.go          # Connect() — dials RabbitMQ, declares durable topic exchange
│   ├── publisher.go     # Publisher.Publish(ctx, Topic, payload) — JSON encodes + routes
│   ├── subscriber.go    # Subscriber.Subscribe(ctx, queue, Topic, Handler) — ack/nack
│   └── topics.go        # Topic constants (user.business.follow, vendor.product.sync, …)
├── ratelimit/
│   └── ratelimit.go     # Rate limiter core (Redis-backed)
├── tokenizer/
│   ├── split_token.go   # SplitAlphaNumeric — splits on letter/digit boundaries
│   ├── singular.go      # Strips plural suffixes
│   ├── stopwords.go     # Filters common stopwords
│   ├── price_filter.go  # Strips price-like tokens
│   ├── unit_filter.go   # Strips unit-like tokens (kg, ml, etc.)
│   └── token_parse.go   # Pipeline: split → stop → singular → filter
└── validate/
    └── validate.go      # Password validation (uppercase + lowercase + digit + ≥8 chars)
```

---

## `database/` — Shared GORM Models

```
database/
├── go.mod
├── postgres.go          # Generic postgres connection helper (GORM silent logging)
├── vendordb/
│   ├── models.go        # Business, Store, Category, Product, InventoryProduct,
│   │                    # BusinessAccess, ProductToken, BusinessAccount,
│   │                    # SearchSyncOutbox, OutboxPayload
│   ├── migrate.go       # AutoMigrate for vendor models
│   ├── repo.go          # Interfaces incl. IInventoryProduct (WriteUpsertOutboxesForProduct,
│   │                    # WriteDeleteOutboxesForProduct, GetByID, AddProduct, …)
│   └── postgres/        # Postgres-specific query implementations
├── search/
│   ├── models.go        # SearchEntry (id, product_id, name, price, price_unit,
│   │                    # search_tokens, category_ids, lat, long, geo_hash, available)
│   │                    # TokenCategoryFreq (token, category_id, count)
│   ├── migrate.go       # AutoMigrate + legacy category_ids jsonb conversion
│   ├── repo.go          # ISearchRepository interface
│   └── postgres/        # JSONB + geo bounding-box search queries, upsert, soft-delete
└── user/
    ├── models.go        # User model
    ├── migrate.go       # AutoMigrate for user models
    ├── repo.go
    └── postgres/
```

---

## `services/vendor/` — Vendor Service

Handles business registration, authentication, product catalog management, inventory, and categories. Propagates inventory changes to the search service via a transactional outbox.

```
services/vendor/
├── go.mod
├── go.sum
├── cmd/
│   └── main.go              # Entry point: calls applayer.Init(), starts HTTP server
├── applayer/
│   ├── init.go              # Wires DB, Redis, RabbitMQ publisher/subscriber,
│   │                        # services, handlers; starts outbox flusher goroutine
│   └── routes.go            # Registers all routes on a gorilla/mux router;
│                            #   public /products/{id}/detail registered before auth subrouter
├── config/
│   └── config.go            # Reads env vars: DATABASE_URL, CACHE_URL, JWT_SECRET,
│                            #   RABBITMQ_URL, USER_SERVICE_URL, PORT
├── connections/
│   ├── database/            # GORM postgres connection + auto-migrate
│   └── cache/               # Redis connection + rate limiter factory
├── handlers/
│   ├── base_handler.go      # Shared handler helpers
│   ├── auth_handler.go      # POST /api/v1/auth/register, /login, /refresh
│   ├── business_handler.go  # GET/PATCH /api/v1/businesses/...
│   ├── category_handler.go  # GET/POST /api/v1/categories/
│   ├── inventory_handler.go # CRUD /api/v1/inventory/...
│   └── product_handler.go   # CRUD /api/v1/products/... + public ProductDetail handler;
│                            #   ProductAdd/Update/Delete each call notify() after write
├── internal_handlers/
│   └── internal.go          # /internal/inventory-products/{id}
│                            # /internal/users/{id}/access
│                            # (network-isolated, no auth middleware)
├── services/
│   ├── base_service.go      # Response code constants, shared service helpers
│   ├── auth_service.go      # Register, login, refresh token
│   ├── business_service.go  # Business profile, update; calls user-service HTTP
│   ├── category_service.go  # Category add, list
│   ├── inventory_service.go # Inventory CRUD, add/remove products
│   └── product_service.go   # Product add, list, get, update, delete;
│                            #   Update writes upsert outboxes, Delete writes delete outboxes
├── repository/
│   ├── iaccess_repo.go
│   ├── ibusinesses_repo.go
│   ├── icategory_repo.go
│   ├── iinventory_product_repo.go
│   ├── iproduct_repo.go
│   ├── iproduct_token_repo.go
│   ├── istore_repo.go
│   └── implementation/      # Concrete GORM implementations of each interface
├── models/
│   └── ...                  # Local model types extending vendordb models
├── events/
│   └── publisher.go         # Outbox flusher: polls search_sync_outboxes with SKIP LOCKED,
│                            #   publishes each payload to RabbitMQ topic vendor.product.sync
├── utils/
│   ├── validation.go
│   └── request/
└── main.tf                  # Terraform (infrastructure as code)
```

### Route Summary

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/health` | — | Health check |
| POST | `/api/v1/auth/register` | — | Create business account |
| POST | `/api/v1/auth/login` | — | Login, returns JWT |
| POST | `/api/v1/auth/refresh` | — | Refresh access token |
| POST | `/api/v1/businesses/register` | JWT | Register business profile |
| GET | `/api/v1/businesses/{bizId}/profile` | JWT | Get business profile |
| PATCH | `/api/v1/businesses/` | JWT + BizOnly | Update business |
| GET | `/api/v1/categories/` | JWT | List categories |
| POST | `/api/v1/categories/` | JWT + BizOnly | Add category |
| POST | `/api/v1/products/` | JWT + BizOnly | Add product |
| GET | `/api/v1/products/` | JWT + BizOnly | List products |
| GET | `/api/v1/products/{invProductId}/detail` | — | Get inventory product detail (public) |
| GET | `/api/v1/products/{id}` | JWT + BizOnly | Get product |
| PATCH | `/api/v1/products/{id}` | JWT + BizOnly | Update product |
| DELETE | `/api/v1/products/{id}` | JWT + BizOnly | Delete product |
| POST | `/api/v1/inventory/` | JWT + BizOnly | Create inventory |
| GET | `/api/v1/inventory/` | JWT + BizOnly | List inventories |
| PATCH | `/api/v1/inventory/{invId}` | JWT + BizOnly | Update inventory |
| GET | `/api/v1/inventory/{invId}/products` | JWT + BizOnly | List inventory products |
| POST | `/api/v1/inventory/{invId}/products` | JWT + BizOnly | Add product to inventory |
| PATCH | `/api/v1/inventory/{invId}/products/{id}` | JWT + BizOnly | Update inventory product |
| DELETE | `/api/v1/inventory/{invId}/products/{id}` | JWT + BizOnly | Remove product from inventory |
| GET | `/internal/inventory-products/{id}` | None | Internal: get inventory product |
| GET | `/internal/users/{id}/access` | None | Internal: get user business access |

### Outbox Pattern (vendor → search)

When inventory products are added, updated, or removed, the vendor service writes a `SearchSyncOutbox` row in the **same database transaction** as the triggering write. A background `Flusher` goroutine polls unprocessed outbox rows every 30 seconds (and immediately on-demand after each write) using `SKIP LOCKED` for multi-instance safety. It publishes each payload as JSON to the RabbitMQ topic **`vendor.product.sync`** and marks the row as processed. The search service subscribes to this topic and applies upsert/delete events to its index. This ensures eventual consistency without direct database coupling or HTTP coupling between services.

---

## `services/search/` — Search Service

Maintains a denormalized search index. Accepts sync events from the vendor service and serves full-text + geo search queries to clients.

```
services/search/
├── go.mod
├── go.sum
├── cmd/
│   └── main.go              # Entry point
├── applayer/
│   ├── init.go              # Wires DB, Redis, RabbitMQ subscriber, freq flusher,
│   │                        # sync consumer, search service
│   └── routes.go            # GET /health, GET /api/v1/search
├── config/
│   └── config.go            # Reads env vars: DATABASE_URL, CACHE_URL, CACHE_PASSWORD,
│                            #   RABBITMQ_URL, PORT
├── connections/
│   └── database/            # GORM postgres connection + auto-migrate
├── cache/
│   └── cache.go             # SearchCache: GetTokenFreqs / SetTokenFreqs with 30 s Redis TTL
│                            #   key format: search:freqs:{sorted-tokens}
├── handlers/
│   └── search_handler.go    # GET /api/v1/search
├── services/
│   ├── search_service.go    # Search query logic; uses SearchCache for freq lookup
│   └── scoring.go           # Scoring (token coverage + category freq + distance tier),
│                            #   dedup by ProductID, cap at 100 results; returns SearchResult
├── repository/
│   └── search_repo.go       # Local repo interface (wraps database/search/postgres)
├── models/
│   └── search_entry.go      # Local mirror of SearchEntry (name, price, price_unit added)
├── freq/
│   └── flusher.go           # Accumulates token-category freq deltas in Redis (HINCRBY),
│                            #   flushes to DB every 30 s
├── sync/
│   ├── consumer.go          # RunConsumer: subscribes to vendor.product.sync via RabbitMQ
│   ├── subscriber.go        # Applier.Apply: routes upsert/delete payloads
│   └── poller.go            # applyUpsert / applyDelete: writes to search_entries
├── replica/                 # (read replica support)
└── main.tf                  # Terraform
```

### Route Summary

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/health` | — | Health check |
| GET | `/api/v1/search` | — | Search products (`?q=`, `?lat=`, `?long=`, `?range=`) |

### Sync Flow

The search service subscribes to RabbitMQ topic **`vendor.product.sync`** (queue `search.product.sync`). Each message carries an `OutboxPayload` with `event_type: "upsert"` or `event_type: "delete"`. On upsert, the search index row is created or replaced (`ON CONFLICT (id) DO UPDATE`). On delete, the row is soft-deleted. Messages are acked on success and nacked+requeued on error.

---

## `app/` — User Service

```
app/
├── cmd/
│   └── main.go              # Entry point for user-service
├── docker-compose.yml       # (local dev compose, separate from root)
├── main.tf                  # Terraform
└── start_server.sh
```

---

## `docker/` — Docker Init Scripts

```
docker/
├── init-db.sql    # Creates spotnearr_vendor and spotnearr_search databases
│                  # Runs on first postgres container start
└── postgres/      # Additional postgres config if needed
```

---

## `dummy/` — Local Development Seed Scripts

Node.js scripts for seeding data against a running local stack.

```
dummy/
├── dummy.js          # CLI: register business, seed categories, products
├── seeder.js         # Bulk seeder
├── load-test.js      # Load testing script
├── products.csv      # Sample product data for seeding
├── package.json
└── node_modules/
```

Usage:
```bash
node dummy.js register BIZ_EMAIL=... BIZ_PASSWORD=... BIZ_NAME=... BIZ_PHONE=...
node dummy.js categories BIZ_EMAIL=... BIZ_PASSWORD=...
```

---

## Infrastructure

### Production (`docker-compose.yml`)

```
postgres  (port 5432) ──┬── spotnearr_user   → user-service   (port 8080)
                        ├── spotnearr_vendor → vendor-service  (port 8081)
                        └── spotnearr_search → search-service  (port 8082)
redis     (port 6379) ──┬── user-service    (rate limiting)
                        ├── vendor-service  (rate limiting, caching)
                        └── search-service  (freq-delta accumulation, token-freq cache)
rabbitmq  (port 5672) ──┬── vendor-service  (publisher: vendor.product.sync,
                        │                    consumer: user.business.follow/unfollow)
                        └── search-service  (consumer: vendor.product.sync)
```

All services share one postgres instance, each with its own database. All three services wait on RabbitMQ's healthcheck before starting (`rabbitmq-diagnostics ping`).

### Test (`docker-compose.test.yml`)

Ephemeral postgres (tmpfs — no disk persistence, fresh DB every run). Two Docker Compose profiles:

- `--profile vendor`: postgres-test + redis-test + vendor-test container
- `--profile search`: postgres-test + search-test container

Test images compile the test binary at Docker build time (`go test -tags integration -c`). A compilation error fails the image build, blocking deployment before any test runs.

### Terraform (`main.tf`)

Each service directory contains a `main.tf` for cloud infrastructure provisioning.

---

## Environment Variables

### vendor-service

| Variable | Description |
|---|---|
| `DATABASE_URL` | Postgres connection string for `spotnearr_vendor` |
| `CACHE_URL` | Redis URL |
| `CACHE_PASSWORD` | Redis password (empty for local) |
| `JWT_SECRET` | Secret for signing JWTs |
| `RABBITMQ_URL` | RabbitMQ AMQP URL (e.g. `amqp://admin:admin123@rabbitmq:5672/`) |
| `USER_SERVICE_URL` | Base URL of user-service (for business access checks) |
| `PORT` | HTTP listen port (default 8080) |

### search-service

| Variable | Description |
|---|---|
| `DATABASE_URL` | Postgres connection string for `spotnearr_search` |
| `CACHE_URL` | Redis URL |
| `CACHE_PASSWORD` | Redis password (empty for local) |
| `RABBITMQ_URL` | RabbitMQ AMQP URL |
| `PORT` | HTTP listen port (default 8080) |

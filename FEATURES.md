# Feature Tracker

Last updated: 2026-06-13

Status legend: `✅ Done` · `🚧 In Progress` · `⬜ Not Started`

---

## User Service — `:8080`

Path: `app/cmd/`

### Auth
| Feature | Status | Route |
|---|---|---|
| Register (with Turnstile captcha) | ✅ Done | `POST /api/v1/auth/users/register` |
| Login (with Turnstile captcha) | ✅ Done | `POST /api/v1/auth/users/login` |
| Refresh access token | ✅ Done | `POST /api/v1/auth/refresh` |
| Password reset request (email) | ✅ Done | `POST /api/v1/auth/reset-request` |
| Password reset (session token) | ✅ Done | `POST /api/v1/auth/reset-password` |
| Get auth status | ✅ Done | `GET /api/v1/auth/` |
| Logout from all devices | ✅ Done | `POST /api/v1/auth/logout-all-devices` |

### Users
| Feature | Status | Route |
|---|---|---|
| Get user profile | ✅ Done | `GET /api/v1/users/{userId}/profile` |
| Update own profile | ✅ Done | `PATCH /api/v1/users/profile` |
| Ban user (admin) | ✅ Done | `POST /api/v1/users/{userId}/ban` |

### Claims (Product Reservations)
| Feature | Status | Route |
|---|---|---|
| List user's active claims | ✅ Done | `GET /api/v1/claims` |
| Reserve a product (claim) | ✅ Done | `POST /api/v1/claims/{invProductId}` |
| Cancel a claim | ✅ Done | `DELETE /api/v1/claims/{claimId}` |

### Reviews
| Feature | Status | Route |
|---|---|---|
| Review an offer — create | ✅ Done | `POST /api/v1/review/offers` |
| Review an offer — read | ✅ Done | `GET /api/v1/review/offers/{offerId}` |
| Review an offer — update | ✅ Done | `PUT /api/v1/review/offers` |
| Review an offer — delete | ✅ Done | `DELETE /api/v1/review/offers` |
| Review a spotlight — create | ✅ Done | `POST /api/v1/review/spotlights` |
| Review a spotlight — read | ✅ Done | `GET /api/v1/review/spotlights/{spotlightId}` |
| Review a spotlight — update | ✅ Done | `PUT /api/v1/review/spotlights` |
| Review a spotlight — delete | ✅ Done | `DELETE /api/v1/review/spotlights` |
| Review a business — create | ✅ Done | `POST /api/v1/review/businesses` |
| Review a business — read | ✅ Done | `GET /api/v1/review/businesses` |
| Review a business — update | ✅ Done | `PUT /api/v1/review/businesses` |
| Review a business — delete | ✅ Done | `DELETE /api/v1/review/businesses` |
| Review a product — create | ✅ Done | `POST /api/v1/review/products` |
| Review a product — read | ✅ Done | `GET /api/v1/review/products` |
| Review a product — update | ✅ Done | `PUT /api/v1/review/products` |
| Review a product — delete | ✅ Done | `DELETE /api/v1/review/products` |

### Social & Feed
| Feature | Status | Notes |
|---|---|---|
| Follow a business | ✅ Done | `POST/DELETE /api/v1/social/businesses/{bizId}/follow` |
| Like a product | ✅ Done | `POST/DELETE /api/v1/social/products/{invProductId}/like` |
| Save a product | ✅ Done | `POST/DELETE /api/v1/social/products/{invProductId}/save` |
| Like a spotlight | ✅ Done | `POST/DELETE /api/v1/social/spotlights/{spotlightId}/like` |
| Save a spotlight | ✅ Done | `POST/DELETE /api/v1/social/spotlights/{spotlightId}/save` |
| List followed businesses | ✅ Done | `GET /api/v1/social/businesses/followed` |
| List liked products | ✅ Done | `GET /api/v1/social/products/liked` |
| List saved products | ✅ Done | `GET /api/v1/social/products/saved` |
| List liked spotlights | ✅ Done | `GET /api/v1/social/spotlights/liked` |
| List saved spotlights | ✅ Done | `GET /api/v1/social/spotlights/saved` |
| Follow another user | ⬜ Not Started | Removed from current scope |
| Spotlight feed (hyperlocal, 10 km) | ✅ Done | `GET /api/v1/feed/spotlights?lat=&long=&range=` — joins spotlights→stores by biz, default 10 km |
| Real-time notifications | ⬜ Not Started | Live offers, nearby updates, comments |

### Internal
| Feature | Status | Route |
|---|---|---|
| Invalidate user refresh token | ✅ Done | `POST /internal/users/{userId}/invalidate-refresh` |
| Get claims by product IDs | ✅ Done | `GET /internal/claims?product_ids=` — used by vendor service |
| Get single claim | ✅ Done | `GET /internal/claims/{claimId}` — used by vendor service |
| Update claim status | ✅ Done | `PATCH /internal/claims/{claimId}/status` — used by vendor service |

---

## Vendor Service — `:8081`

Path: `services/vendor/`

### Auth
| Feature | Status | Route |
|---|---|---|
| Register business account | ✅ Done | `POST /api/v1/auth/register` |
| Login | ✅ Done | `POST /api/v1/auth/login` |
| Refresh access token | ✅ Done | `POST /api/v1/auth/refresh` |

### Business Profile
| Feature | Status | Route |
|---|---|---|
| Register business profile | ✅ Done | `POST /api/v1/businesses/register` |
| Get business profile | ✅ Done | `GET /api/v1/businesses/{bizId}/profile` |
| Update business (name / desc) | ✅ Done | `PATCH /api/v1/businesses/` |

### Products
| Feature | Status | Route |
|---|---|---|
| Create product | ✅ Done | `POST /api/v1/products/` |
| List all products | ✅ Done | `GET /api/v1/products/` |
| Get single product | ✅ Done | `GET /api/v1/products/{productId}` |
| Get inventory product detail (public) | ✅ Done | `GET /api/v1/products/{invProductId}/detail` — no auth |
| Update product | ✅ Done | `PATCH /api/v1/products/{productId}` |
| Soft-delete product | ✅ Done | `DELETE /api/v1/products/{productId}` |

### Inventory (Stores)
| Feature | Status | Route |
|---|---|---|
| Create store location | ✅ Done | `POST /api/v1/inventory/` |
| List all stores | ✅ Done | `GET /api/v1/inventory/` |
| Update store | ✅ Done | `PATCH /api/v1/inventory/{invId}` |
| List products in a store | ✅ Done | `GET /api/v1/inventory/{invId}/products` |
| Add product to store | ✅ Done | `POST /api/v1/inventory/{invId}/products` |
| Update product stock/availability | ✅ Done | `PATCH /api/v1/inventory/{invId}/products/{invProductId}` |
| Remove product from store | ✅ Done | `DELETE /api/v1/inventory/{invId}/products/{invProductId}` |

### Categories
| Feature | Status | Route |
|---|---|---|
| List all categories | ✅ Done | `GET /api/v1/categories/` |
| Create category | ✅ Done | `POST /api/v1/categories/` |

### Offers & Coupons
| Feature | Status | Route |
|---|---|---|
| Create offer / coupon | ✅ Done | `POST /api/v1/offers/` |
| List offers | ✅ Done | `GET /api/v1/offers/` |
| Get single offer | ✅ Done | `GET /api/v1/offers/{offerId}` |
| Update offer | ✅ Done | `PATCH /api/v1/offers/{offerId}` |
| Delete offer | ✅ Done | `DELETE /api/v1/offers/{offerId}` |

### Spotlight Marketing
| Feature | Status | Route |
|---|---|---|
| Post spotlight (image/video) | ✅ Done | `POST /api/v1/spotlights/` — types: product, offer, general |
| List spotlights | ✅ Done | `GET /api/v1/spotlights/` |
| Get single spotlight | ✅ Done | `GET /api/v1/spotlights/{spotlightId}` |
| Delete spotlight | ✅ Done | `DELETE /api/v1/spotlights/{spotlightId}` |

### Claim Management
| Feature | Status | Route |
|---|---|---|
| View incoming claims | ✅ Done | `GET /api/v1/claims` — scoped to vendor's inventory products |
| Accept a claim | ✅ Done | `PATCH /api/v1/claims/{claimId}/accept` |
| Reject a claim | ✅ Done | `PATCH /api/v1/claims/{claimId}/reject` |

### Internal
| Feature | Status | Route |
|---|---|---|
| Get inventory product by ID | ✅ Done | `GET /internal/inventory-products/{id}` |
| Get user access by ID | ✅ Done | `GET /internal/users/{id}/access` |

---

## Search Service — `:8082`

Path: `services/search/`

### Search
| Feature | Status | Route |
|---|---|---|
| Proximity product search | ✅ Done | `GET /api/v1/search?q=&lat=&long=&range=` |
| Result scoring / ranking | ✅ Done | `services/search/services/scoring.go` |
| Search results include name + price | ✅ Done | Stored in `search_entries`; returned as `SearchResult` objects |
| Result cap at 100 | ✅ Done | Applied in `rankProducts` after dedup |
| Search index sync via RabbitMQ | ✅ Done | Topic `vendor.product.sync` → `services/search/sync/consumer.go` |
| Query frequency flusher | ✅ Done | `services/search/freq/flusher.go` — flushes Redis deltas to DB every 30 s |
| Token-category freq caching | ✅ Done | `services/search/cache/cache.go` — Redis cache with 30 s TTL, shared across users |
| Category filter | ✅ Done | `?category_ids=1,2,3` — jsonb `&&` overlap on `search_entries.category_ids` |
| Price range filter | ✅ Done | `?min_price=10&max_price=100` — applied as SQL WHERE on `search_entries.price` |

---

## Cross-Cutting Infrastructure

### Gateway (`gateway/`)
| Component | Status | Notes |
|---|---|---|
| Nginx reverse proxy | ✅ Done | `gateway/nginx.conf` — single entry point on port `80` |
| Path-prefix routing | ✅ Done | `/user/` → user-service, `/vendor/` → vendor-service, `/search/` → search-service; prefix stripped before forwarding |
| Service port isolation | ✅ Done | Go services have no host `ports:` binding — only nginx port `80` is externally reachable |

### Messaging (`pkg/mq`)
| Component | Status | Notes |
|---|---|---|
| RabbitMQ container | ✅ Done | `docker-compose.yml` — port `5672`, management UI `15672`; all services wait on healthcheck |
| Connection wrapper | ✅ Done | `pkg/mq/conn.go` — `Connect()`, durable topic exchange `spotnearr.events`; logs on connect |
| Predefined topics | ✅ Done | `pkg/mq/topics.go` — add new event types here |
| Publisher | ✅ Done | `pkg/mq/publisher.go` — `Publish(ctx, Topic, payload)`; logs "ready" on init |
| Subscriber | ✅ Done | `pkg/mq/subscriber.go` — `Subscribe(ctx, queue, Topic, Handler)` with ack/nack; logs "ready" on init |
| Follow event (user → vendor) | ✅ Done | Topic `user.business.follow` → vendor increments `follower_count` |
| Unfollow event (user → vendor) | ✅ Done | Topic `user.business.unfollow` → vendor decrements `follower_count` |
| Product sync (vendor → search) | ✅ Done | Topic `vendor.product.sync` → search upserts/deletes index entry |

### Observability
| Component | Status | Notes |
|---|---|---|
| GORM query logging | ✅ Done | Disabled (`logger.Silent`) — `database/postgres.go` |
| RabbitMQ startup logs | ✅ Done | Logs connection, publisher ready, subscriber ready on startup |

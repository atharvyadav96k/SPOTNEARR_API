# Feature Tracker

Last updated: 2026-06-12

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
| Follow a business | 🚧 In Progress | Routes: `POST/DELETE /api/v1/social/businesses/{bizId}/follow` |
| Like a product | 🚧 In Progress | Routes: `POST/DELETE /api/v1/social/products/{invProductId}/like` |
| Save a product | 🚧 In Progress | Routes: `POST/DELETE /api/v1/social/products/{invProductId}/save` |
| Like a spotlight | 🚧 In Progress | Routes: `POST/DELETE /api/v1/social/spotlights/{spotlightId}/like` |
| Save a spotlight | 🚧 In Progress | Routes: `POST/DELETE /api/v1/social/spotlights/{spotlightId}/save` |
| Follow another user | ⬜ Not Started | Removed from current scope |
| Spotlight feed (hyperlocal, 10 km) | ⬜ Not Started | |
| Real-time notifications | ⬜ Not Started | Live offers, nearby updates, comments |

### Internal
| Feature | Status | Route |
|---|---|---|
| Invalidate user refresh token | ✅ Done | `POST /internal/users/{userId}/invalidate-refresh` |

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
| Feature | Status | Notes |
|---|---|---|
| Create offer / coupon | ⬜ Not Started | Flat discount, %, min order threshold |
| List offers | ⬜ Not Started | |
| Update offer | ⬜ Not Started | |
| Delete offer | ⬜ Not Started | |

### Spotlight Marketing
| Feature | Status | Notes |
|---|---|---|
| Post spotlight (image/video) | ⬜ Not Started | Types: Product, Offer, General |
| List spotlights | ⬜ Not Started | |
| Delete spotlight | ⬜ Not Started | |

### Claim Management
| Feature | Status | Notes |
|---|---|---|
| View incoming claims | ⬜ Not Started | |
| Accept / reject a claim | ⬜ Not Started | |

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
| Search index sync (outbox poller) | ✅ Done | `services/search/sync/poller.go` |
| Query frequency flusher | ✅ Done | `services/search/freq/flusher.go` |
| Category filter | ⬜ Not Started | Filter results by category IDs |
| Price range filter | ⬜ Not Started | |

### Internal
| Feature | Status | Route |
|---|---|---|
| Manual sync trigger | ✅ Done | `POST /internal/sync` |

# Spotnearr API Reference

## Services & Base URLs

All three services sit behind an **Nginx gateway** (port `80`). Individual service ports are internal-only and not reachable from outside Docker. Clients always talk to the gateway.

| Service | Public Path Prefix | Internal Port | Purpose |
|---|---|---|---|
| **User Service** | `/user` | `8080` | Customer accounts, auth, social, reviews, claims |
| **Vendor Service** | `/vendor` | `8081` | Business owner dashboard — products, inventory, offers, spotlights |
| **Search Service** | `/search` | `8082` | Product discovery (public, no auth) |

**Public URL format:** `http(s)://<host>/<prefix>/<path>`

The gateway strips the prefix before forwarding, so the Go service sees the original path unchanged:

| Client sends | Gateway forwards to |
|---|---|
| `POST /user/api/v1/auth/users/login` | `POST user-service:8080/api/v1/auth/users/login` |
| `GET /vendor/api/v1/products/` | `GET vendor-service:8080/api/v1/products/` |
| `GET /search/api/v1/search?q=apples` | `GET search-service:8080/api/v1/search?q=apples` |

All route paths in this document are the **service-side paths** (after prefix stripping). Prepend the prefix above for the actual public URL. CORS is enabled on every service.

---

## Standard Response Envelope

Every endpoint returns this shape:

```json
{ "message": "string", "data": <object | array | null> }
```

HTTP status is set on the response. `data` is `null` (or omitted) on error responses.

### Common Status Codes

| Code | Meaning |
|---|---|
| `200` | OK |
| `201` | Created |
| `400` | Validation error / bad request body |
| `401` | Missing, expired, or invalid JWT |
| `403` | Valid JWT but insufficient role (e.g. non-business token on BusinessOnly route) |
| `404` | Resource not found or not owned by caller |
| `409` | Conflict — duplicate (email, coupon code, slug, already-claimed, etc.) |
| `500` | Internal server error |
| `503` | Upstream service (vendor / user) unreachable |

---

## Auth Header

All protected endpoints require:

```
Authorization: Bearer <access_token>
```

**BusinessOnly** endpoints additionally require the token to carry a `business_id` claim (present after logging into the Vendor Service).

**Captcha** endpoints require:

```
X-Captcha-Token: <cloudflare_turnstile_token>
```

---

## JWT Claims (decoded payload)

```json
{
  "user_id":    42,
  "business_id": 1,
  "token_type": "access",
  "role":       "user",
  "exp":        1700000000
}
```

- `business_id` is `null` for pure user accounts
- `role`: `"user"` | `"owner"`
- Access token lifetime: ~5 minutes — refresh before expiry

---

---

# User Service — `:8080`

End-user (customer) accounts. Register / login requires Cloudflare Turnstile.

---

## Auth

### `POST /api/v1/auth/users/register`

**Headers:** `X-Captcha-Token`

**Request body:**
```json
{
  "name":     "John Doe",
  "email":    "john@example.com",
  "phone":    "+919876543210",
  "password": "Abcd1234"
}
```

| Field | Required | Rules |
|---|---|---|
| `name` | Yes | Non-empty string |
| `email` | Yes | Valid email format |
| `phone` | Yes | E.164 format |
| `password` | Yes | Min 8 chars, must include uppercase, lowercase, and digit |

**Response `200`:**
```json
{
  "message": "User registered successfully",
  "data": {
    "access_token":            "<jwt>",
    "access_token_expires_at":  1749259500,
    "refresh_token":           "<jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```

---

### `POST /api/v1/auth/users/login`

**Headers:** `X-Captcha-Token`

**Request body:**
```json
{ "email": "john@example.com", "password": "Abcd1234" }
```

**Response `200`:** Same token shape as register.

---

### `POST /api/v1/auth/refresh`

Exchange a valid refresh token for a new access token. No captcha required.

**Request body:**
```json
{
  "access_token":  "<expired_or_active_jwt>",
  "refresh_token": "<refresh_jwt>"
}
```

**Response `200`:**
```json
{
  "message": "New access token",
  "data": {
    "access_token":            "<new_jwt>",
    "access_token_expires_at":  1749259500,
    "refresh_token":           "<same_refresh_jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```

---

### `GET /api/v1/auth/`

Auth: **JWT**. Validates the caller's token. Returns `200` if valid, `401` otherwise.

**Response `200`:** `{ "message": "", "data": null }`

---

### `POST /api/v1/auth/logout-all-devices`

Auth: **JWT**. Invalidates the caller's refresh token so no device can refresh.

**Response `200`:** `{ "message": "Successfully logged out from all devices" }`

---

### `POST /api/v1/auth/reset-request`

**Headers:** `X-Captcha-Token`

**Request body:**
```json
{ "email": "john@example.com" }
```

Generates a one-time session link (logged server-side; email delivery not yet wired).

**Response `200`:** `{ "message": "", "data": "session=<token>" }`

---

### `POST /api/v1/auth/reset-password?session=<token>`

**Headers:** `X-Captcha-Token`. Requires the session query param from the reset-request step.

**Request body:**
```json
{ "password": "NewPass123" }
```

Password rules: min 8 chars, uppercase + lowercase + digit.

**Response `200`:** `{ "message": "Password updated successfully" }`. Invalidates all refresh tokens.

---

## Users

### `GET /api/v1/users/{userId}/profile`

Auth: **JWT**. Returns the caller's own profile (the path `userId` is present but the user is resolved from the JWT).

**Response `200`:**
```json
{
  "message": "User profile",
  "data": {
    "id":             1,
    "fullName":       "John Doe",
    "email":          "john@example.com",
    "phone":          "+919876543210",
    "isVerifiedEmail": false,
    "isVerifiedPhone": false,
    "isActive":       true,
    "createdAt":      "2024-01-01T00:00:00Z",
    "updatedAt":      "2024-01-01T00:00:00Z"
  }
}
```

---

### `PATCH /api/v1/users/profile`

Auth: **JWT**. Update the caller's own display name.

**Request body:**
```json
{ "fullName": "Jane Doe" }
```

| Field | Required | Rules |
|---|---|---|
| `fullName` | Yes | Non-empty string |

**Response `200`:** `{ "message": "Profile updated" }`

---

### `POST /api/v1/users/{userId}/ban`

Auth: **JWT**. Deactivates the target account and invalidates their refresh token. (Admin intent — no role guard in current code, the JWT must be valid.)

**Response `200`:** `{ "message": "User banned" }`

---

## Claims (Product Reservations)

### `GET /api/v1/claims`

Auth: **JWT**. Returns all active claims for the authenticated user.

**Response `200`:**
```json
{
  "message": "Claims fetched successfully",
  "data": [
    {
      "id":                  1,
      "userId":              42,
      "inventoryProductId":  7,
      "status":              "pending",
      "product": {
        "id": 7, "name": "Organic Apples", "price": 120.00,
        "available": true, "storeName": "Main Branch"
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

`product` is a snapshot taken at claim time — it does not change if the vendor later updates the product.

`status`: `"pending"` | `"accepted"` | `"rejected"`

---

### `POST /api/v1/claims/{invProductId}`

Auth: **JWT**. Reserve an inventory product.

- Returns `404` if the product doesn't exist in the vendor's inventory.
- Returns `400` if the product is marked `available: false`.
- Returns `409` if the user already has an active claim on this product.

**Response `201`:**
```json
{
  "message": "Claim submitted successfully",
  "data": {
    "id":                  1,
    "userId":              42,
    "inventoryProductId":  7,
    "status":              "pending",
    "product":             { ... snapshot ... },
    "createdAt":           "2024-01-01T00:00:00Z",
    "updatedAt":           "2024-01-01T00:00:00Z"
  }
}
```

---

### `DELETE /api/v1/claims/{claimId}`

Auth: **JWT**. Cancel a pending claim. Returns `404` if claim doesn't belong to the caller.

**Response `200`:** `{ "message": "Claim removed successfully" }`

---

## Reviews

All review endpoints require **JWT**. GET endpoints are publicly accessible (same JWT required in current implementation).

Review target types: `"business"` | `"product"` | `"offer"` | `"spotlight"`

### Review response object

```json
{
  "id":         1,
  "userId":     42,
  "targetType": "offer",
  "targetId":   3,
  "stars":      5,
  "comment":    "Great deal!",
  "createdAt":  "2024-01-01T00:00:00Z",
  "updatedAt":  "2024-01-01T00:00:00Z"
}
```

---

### Offer Reviews

| Method | Route | Description |
|---|---|---|
| `GET` | `/api/v1/review/offers/{offerId}` | List all reviews for an offer |
| `POST` | `/api/v1/review/offers` | Add a review |
| `PUT` | `/api/v1/review/offers` | Update your review |
| `DELETE` | `/api/v1/review/offers` | Delete your review |

**POST / PUT body:**
```json
{ "targetId": 3, "stars": 4, "comment": "Really good offer" }
```

| Field | Required | Rules |
|---|---|---|
| `targetId` | Yes | ID of the offer being reviewed |
| `stars` | Yes | Integer 1–5 |
| `comment` | No | Free text |

**DELETE body:**
```json
{ "targetId": 3, "stars": 1, "comment": "" }
```
Only `targetId` is used for deletion; `stars` field must still pass validation (1–5).

**GET Response `200`:**
```json
{ "data": [ { ...review... }, ... ] }
```

**POST Response `200`:** `{ "data": { ...created review... } }`

---

### Spotlight Reviews

| Method | Route |
|---|---|
| `GET` | `/api/v1/review/spotlights/{spotlightId}` |
| `POST` | `/api/v1/review/spotlights` |
| `PUT` | `/api/v1/review/spotlights` |
| `DELETE` | `/api/v1/review/spotlights` |

Same body shape as offer reviews. `targetId` = spotlight ID.

---

### Business Reviews

| Method | Route |
|---|---|
| `GET` | `/api/v1/review/businesses?id={bizId}` |
| `POST` | `/api/v1/review/businesses` |
| `PUT` | `/api/v1/review/businesses` |
| `DELETE` | `/api/v1/review/businesses` |

`targetId` = business ID.

---

### Product Reviews

| Method | Route |
|---|---|
| `GET` | `/api/v1/review/products?id={invProductId}` |
| `POST` | `/api/v1/review/products` |
| `PUT` | `/api/v1/review/products` |
| `DELETE` | `/api/v1/review/products` |

`targetId` = inventory product ID.

---

## Social

All social endpoints require **JWT**.

### Follow / Unfollow a Business

| Method | Route | Description |
|---|---|---|
| `POST` | `/api/v1/social/businesses/{bizId}/follow` | Follow a business |
| `DELETE` | `/api/v1/social/businesses/{bizId}/follow` | Unfollow a business |
| `GET` | `/api/v1/social/businesses/followed` | List followed business IDs |

**POST Response `200`:** `{ "message": "Followed successfully" }`
**DELETE Response `200`:** `{ "message": "Unfollowed successfully" }`
**GET Response `200`:**
```json
{ "message": "Followed businesses", "data": [1, 4, 7] }
```
Returns an array of business IDs. Returns `[]` (empty array) if none.

Following publishes a `user.business.follow` event to RabbitMQ so the vendor service increments `followerCount`.

---

### Like / Save a Product

| Method | Route | Description |
|---|---|---|
| `POST` | `/api/v1/social/products/{invProductId}/like` | Like an inventory product |
| `DELETE` | `/api/v1/social/products/{invProductId}/like` | Unlike |
| `GET` | `/api/v1/social/products/liked` | List liked inventory product IDs |
| `POST` | `/api/v1/social/products/{invProductId}/save` | Save a product |
| `DELETE` | `/api/v1/social/products/{invProductId}/save` | Unsave |
| `GET` | `/api/v1/social/products/saved` | List saved inventory product IDs |

**GET Response `200`:**
```json
{ "message": "Liked products", "data": [7, 12, 33] }
```

---

### Like / Save a Spotlight

| Method | Route | Description |
|---|---|---|
| `POST` | `/api/v1/social/spotlights/{spotlightId}/like` | Like a spotlight |
| `DELETE` | `/api/v1/social/spotlights/{spotlightId}/like` | Unlike |
| `GET` | `/api/v1/social/spotlights/liked` | List liked spotlight IDs |
| `POST` | `/api/v1/social/spotlights/{spotlightId}/save` | Save a spotlight |
| `DELETE` | `/api/v1/social/spotlights/{spotlightId}/save` | Unsave |
| `GET` | `/api/v1/social/spotlights/saved` | List saved spotlight IDs |

**GET Response `200`:**
```json
{ "message": "Liked spotlights", "data": [2, 5] }
```

---

---

# Vendor Service — `:8081`

Business owner dashboard. Vendor has its own auth flow — independent of the user service.

---

## Auth

### `POST /api/v1/auth/register`

Create a vendor (business owner) account.

**Request body:**
```json
{
  "name":     "Owner Name",
  "email":    "owner@example.com",
  "phone":    "+919876543210",
  "password": "Abcd1234",
  "desc":     "optional business description"
}
```

| Field | Required | Rules |
|---|---|---|
| `name` | Yes | Non-empty |
| `email` | Yes | Valid email |
| `phone` | Yes | E.164 format |
| `password` | Yes | Min 8 chars, uppercase + lowercase + digit |
| `desc` | No | — |

**Response `200`:**
```json
{
  "data": {
    "access_token":            "<jwt>",
    "access_token_expires_at":  1749259500,
    "refresh_token":           "<jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```

---

### `POST /api/v1/auth/login`

**Request body:**
```json
{ "email": "owner@example.com", "password": "Abcd1234" }
```

**Response `200`:** Same token shape as register.

---

### `POST /api/v1/auth/refresh`

**Request body:**
```json
{ "access_token": "<jwt>", "refresh_token": "<jwt>" }
```

**Response `200`:** New `access_token` + `access_token_expires_at`.

---

## Business Profile

### `POST /api/v1/businesses/register`

Auth: **JWT**. Register a business profile linked to the authenticated account.

**Request body:**
```json
{
  "name":  "My Store",
  "email": "store@example.com",
  "phone": "+91...",
  "desc":  "optional"
}
```

| Field | Required |
|---|---|
| `name` | Yes |
| `email` | Yes |
| `phone` | Yes |
| `desc` | No |

**Response `200`:**
```json
{
  "data": {
    "id":               1,
    "businessName":     "My Store",
    "email":            "store@example.com",
    "phone":            "+91...",
    "description":      "...",
    "isActive":         true,
    "verifiedBusiness": false,
    "followerCount":    0,
    "createdAt":        "2024-01-01T00:00:00Z",
    "updatedAt":        "2024-01-01T00:00:00Z"
  }
}
```

---

### `GET /api/v1/businesses/{bizId}/profile`

Auth: **JWT**. Get any business profile by ID.

**Response `200`:** Business object in `data` (same shape as above).

---

### `PATCH /api/v1/businesses/`

Auth: **JWT + BusinessOnly**. Update the authenticated business's name or description.

**Request body:**
```json
{ "name": "Updated Name", "desc": "Updated description" }
```

| Field | Required |
|---|---|
| `name` | Yes |
| `desc` | No |

**Response `200`:** Updated business object in `data`.

---

## Products

All product endpoints require **JWT + BusinessOnly** except `GET .../detail`.

### `POST /api/v1/products/`

Create a product.

**Request body:**
```json
{
  "name": "Organic Apples",
  "price":    { "value": 120.00, "unit": "rs" },
  "quantity": { "value": 1.0,   "unit": "kg" },
  "desc":        "optional description",
  "categoryIds": [1, 3],
  "storeIds":    [1]
}
```

| Field | Required | Rules |
|---|---|---|
| `name` | Yes | Non-empty |
| `price.value` | Yes | `> 0` |
| `price.unit` | Yes | Non-empty string (e.g. `"rs"`, `"usd"`) |
| `quantity.unit` | Yes | Non-empty string (e.g. `"kg"`, `"pcs"`) |
| `categoryIds` | Yes | At least one category ID |
| `storeIds` | No | Auto-links to each listed store |

**Response `201`:**
```json
{
  "data": {
    "id":           5,
    "name":         "Organic Apples",
    "price":        120.00,
    "priceUnit":    "rs",
    "quantity":     1.0,
    "quantityUnit": "kg",
    "desc":         "...",
    "businessId":   1,
    "categories":   [{ "id": 1, "name": "Fruits", "slug": "fruits" }],
    "createdAt":    "2024-01-01T00:00:00Z",
    "updatedAt":    "2024-01-01T00:00:00Z"
  }
}
```

---

### `GET /api/v1/products/`

List all products for the authenticated business.

**Response `200`:** `data` is an array of product objects.

---

### `GET /api/v1/products/{productId}`

Get a single product owned by the authenticated business. Returns `404` if not owned.

**Response `200`:** Single product object in `data`.

---

### `GET /api/v1/products/{invProductId}/detail`

**Public — no auth required.** Returns full inventory product detail including store and product info.

**Response `200`:**
```json
{
  "data": {
    "id":        7,
    "storeId":   1,
    "productId": 5,
    "count":     10,
    "available": true,
    "store": {
      "id": 1, "name": "Main Branch", "streetAddress": "123 Market St",
      "lat": 12.9716, "long": 77.5946
    },
    "product": {
      "id": 5, "name": "Organic Apples", "price": 120.00, "priceUnit": "rs",
      "categories": [{ "id": 1, "name": "Fruits", "slug": "fruits" }]
    },
    "createdAt": "...", "updatedAt": "..."
  }
}
```

---

### `PATCH /api/v1/products/{productId}`

Update a product. `price` and `quantity` are required even if unchanged.

**Request body:**
```json
{
  "name":     "Updated Name",
  "price":    { "value": 130.00, "unit": "rs" },
  "quantity": { "value": 1.0,   "unit": "kg" },
  "desc":     "optional"
}
```

**Response `200`:** Updated product object in `data`.

---

### `DELETE /api/v1/products/{productId}`

Soft-delete a product. Propagates a delete event to the search index via RabbitMQ.

**Response `200`:** `{ "message": "..." }`

---

## Inventory (Stores)

All inventory endpoints require **JWT + BusinessOnly**.

### `POST /api/v1/inventory/`

Create a store location.

**Request body:**
```json
{
  "name":          "Main Branch",
  "streetAddress": "123 Market St",
  "lat":           12.9716,
  "long":          77.5946
}
```

| Field | Required |
|---|---|
| `name` | Yes |
| `streetAddress` | Yes |
| `lat` | Yes |
| `long` | Yes |

`geoHash` is computed automatically.

**Response `200`:**
```json
{
  "data": {
    "id":            1,
    "name":          "Main Branch",
    "streetAddress": "123 Market St",
    "businessId":    1,
    "lat":           12.9716,
    "long":          77.5946,
    "geoHash":       "tdr1u",
    "createdAt":     "2024-01-01T00:00:00Z",
    "updatedAt":     "2024-01-01T00:00:00Z"
  }
}
```

---

### `GET /api/v1/inventory/`

List all stores for the authenticated business.

**Response `200`:** `data` is an array of store objects.

---

### `PATCH /api/v1/inventory/{invId}`

Update a store. `lat` and `long` are required (used to recompute `geoHash`).

**Request body:**
```json
{
  "name":    "Updated Name",
  "address": "New Street",
  "lat":     12.9716,
  "long":    77.5946
}
```

**Response `200`:** Updated store object in `data`.

---

### `GET /api/v1/inventory/{invId}/products`

List all inventory products (stock) in a store.

**Response `200`:**
```json
{
  "data": [
    {
      "id":        7,
      "storeId":   1,
      "productId": 5,
      "count":     10,
      "available": true,
      "createdAt": "...",
      "updatedAt": "..."
    }
  ]
}
```

---

### `POST /api/v1/inventory/{invId}/products`

Add a product to a store.

**Request body:**
```json
{ "productID": 5, "count": 10, "available": true }
```

| Field | Required | Notes |
|---|---|---|
| `productID` | Yes | Must belong to the same business |
| `count` | No | `null` = unlimited |
| `available` | No | Defaults to `true` |

**Response `200`:** Inventory product object in `data`.

---

### `PATCH /api/v1/inventory/{invId}/products/{invProductId}`

Update stock count or availability.

**Request body:**
```json
{ "count": 20, "available": true }
```

Both fields are optional — omit to leave unchanged.

**Response `200`:** Updated inventory product in `data`.

---

### `DELETE /api/v1/inventory/{invId}/products/{invProductId}`

Remove a product from a store. Publishes a delete event to the search index.

**Response `200`:** `{ "message": "..." }`

---

## Categories

### `GET /api/v1/categories/`

Auth: **JWT**. List all categories.

**Response `200`:**
```json
{
  "data": [
    { "id": 1, "name": "Fruits", "slug": "fruits" }
  ]
}
```

---

### `POST /api/v1/categories/`

Auth: **JWT + BusinessOnly**. Create a category.

**Request body:**
```json
{ "name": "Dairy", "slug": "dairy" }
```

Both `name` and `slug` must be globally unique.

**Response `201`:**
```json
{ "data": { "id": 3, "name": "Dairy", "slug": "dairy" } }
```

---

## Offers & Coupons

All offer endpoints require **JWT + BusinessOnly**.

### Offer object

```json
{
  "id":            1,
  "businessId":    1,
  "title":         "Weekend Sale",
  "description":   "20% off all fruits",
  "discountType":  "percent",
  "discountValue": 20.00,
  "minOrderValue": 100.00,
  "code":          "SAVE20",
  "maxUsage":      500,
  "usedCount":     12,
  "active":        true,
  "expiresAt":     "2024-12-31T23:59:59Z",
  "createdAt":     "2024-01-01T00:00:00Z",
  "updatedAt":     "2024-01-01T00:00:00Z"
}
```

---

### `POST /api/v1/offers/`

Create an offer or coupon.

**Request body:**
```json
{
  "title":         "Weekend Sale",
  "description":   "20% off all fruits",
  "discountType":  "percent",
  "discountValue": 20.00,
  "minOrderValue": 100.00,
  "code":          "SAVE20",
  "maxUsage":      500,
  "expiresAt":     "2024-12-31T23:59:59Z"
}
```

| Field | Required | Rules |
|---|---|---|
| `title` | Yes | Non-empty |
| `discountType` | Yes | `"flat"` or `"percent"` |
| `discountValue` | Yes | `> 0`; max `100` when type is `"percent"` |
| `description` | No | — |
| `minOrderValue` | No | Minimum cart value to apply offer |
| `code` | No | Globally unique coupon code |
| `maxUsage` | No | `null` = unlimited |
| `expiresAt` | No | ISO 8601 datetime |

**Response `201`:** Offer object in `data`.

---

### `GET /api/v1/offers/`

List all offers for the authenticated business.

**Response `200`:** `data` is an array of offer objects.

---

### `GET /api/v1/offers/{offerId}`

Get a single offer. Returns `404` if not owned by the caller's business.

**Response `200`:** Offer object in `data`.

---

### `PATCH /api/v1/offers/{offerId}`

Update an offer. All fields are optional.

**Request body:**
```json
{
  "title":         "New Title",
  "description":   "Updated desc",
  "discountType":  "flat",
  "discountValue": 50.00,
  "minOrderValue": 200.00,
  "code":          "FLAT50",
  "maxUsage":      100,
  "active":        false,
  "expiresAt":     "2025-06-30T23:59:59Z"
}
```

`active: false` deactivates the offer without deleting it.

**Response `200`:** Updated offer object in `data`.

---

### `DELETE /api/v1/offers/{offerId}`

Soft-delete an offer.

**Response `200`:** `{ "message": "offer deleted" }`

---

## Spotlight Marketing

All spotlight endpoints require **JWT + BusinessOnly** except the public feed.

### Spotlight object

```json
{
  "id":          1,
  "businessId":  1,
  "type":        "product",
  "title":       "New Arrival!",
  "caption":     "Fresh organic apples just arrived",
  "mediaUrl":    "https://cdn.example.com/img/apples.jpg",
  "mediaType":   "image",
  "productId":   5,
  "offerId":     null,
  "likeCount":   0,
  "createdAt":   "2024-01-01T00:00:00Z",
  "updatedAt":   "2024-01-01T00:00:00Z"
}
```

---

### `POST /api/v1/spotlights/`

Post a spotlight.

**Request body:**
```json
{
  "type":      "product",
  "title":     "New Arrival!",
  "caption":   "Fresh organic apples just arrived",
  "mediaUrl":  "https://cdn.example.com/img/apples.jpg",
  "mediaType": "image",
  "productId": 5,
  "offerId":   null
}
```

| Field | Required | Rules |
|---|---|---|
| `type` | Yes | `"product"` \| `"offer"` \| `"general"` |
| `title` | Yes | Non-empty |
| `mediaUrl` | Yes | Non-empty URL |
| `mediaType` | Yes | `"image"` \| `"video"` |
| `caption` | No | — |
| `productId` | No | Required when `type` is `"product"` (convention) |
| `offerId` | No | Required when `type` is `"offer"` (convention) |

**Response `201`:** Spotlight object in `data`.

---

### `GET /api/v1/spotlights/`

List all spotlights for the authenticated business, newest first.

**Response `200`:** `data` is an array of spotlight objects.

---

### `GET /api/v1/spotlights/{spotlightId}`

Get a single spotlight. Returns `404` if not owned.

**Response `200`:** Spotlight object in `data`.

---

### `DELETE /api/v1/spotlights/{spotlightId}`

Soft-delete a spotlight.

**Response `200`:** `{ "message": "spotlight deleted" }`

---

## Claim Management (Vendor-side)

All claim management endpoints require **JWT + BusinessOnly**.

Claims are owned by the User Service. The Vendor Service coordinates via internal HTTP — vendors cannot see claims for other vendors' products.

### `GET /api/v1/claims`

List all incoming claims for the authenticated business's inventory products.

**Response `200`:**
```json
{
  "message": "claims",
  "data": [
    {
      "id":                  1,
      "userId":              42,
      "inventoryProductId":  7,
      "status":              "pending",
      "product":             { ... snapshot ... },
      "createdAt":           "2024-01-01T00:00:00Z",
      "updatedAt":           "2024-01-01T00:00:00Z"
    }
  ]
}
```

Returns `[]` if the business has no inventory products.

---

### `PATCH /api/v1/claims/{claimId}/accept`

Accept a claim. Returns `403` if the claimed product doesn't belong to this business.

**No request body required.**

**Response `200`:** `{ "message": "claim accepted" }`

---

### `PATCH /api/v1/claims/{claimId}/reject`

Reject a claim. Same ownership validation as accept.

**No request body required.**

**Response `200`:** `{ "message": "claim rejected" }`

---

## Spotlight Feed (Public)

### `GET /api/v1/feed/spotlights`

**No auth required.** Returns spotlights from businesses that have at least one store within `range` km of the given coordinates. Results are ordered newest first, capped at 50.

**Query params:**

| Param | Required | Default | Description |
|---|---|---|---|
| `lat` | Yes | — | Latitude (decimal degrees) |
| `long` | Yes | — | Longitude (decimal degrees) |
| `range` | No | `10.0` | Radius in km |

**Example:** `GET /api/v1/feed/spotlights?lat=12.97&long=77.59&range=5`

**Response `200`:**
```json
{
  "message": "spotlight feed",
  "data": [
    {
      "id":         1,
      "businessId": 1,
      "type":       "product",
      "title":      "New Arrival!",
      "caption":    "Fresh organic apples",
      "mediaUrl":   "https://cdn.example.com/img/apples.jpg",
      "mediaType":  "image",
      "productId":  5,
      "offerId":    null,
      "likeCount":  14,
      "business": {
        "id": 1, "businessName": "Green Market",
        "isActive": true, "verifiedBusiness": false
      },
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

---

# Search Service — `:8082`

Public, no auth required.

---

## `GET /api/v1/search`

Tokenizes the query, boosts results by category affinity and distance tier, returns up to 100 ranked products.

**Query params:**

| Param | Required | Default | Description |
|---|---|---|---|
| `q` | Yes | — | Search query string |
| `lat` | No | — | Latitude — enables geo-ranking |
| `long` | No | — | Longitude |
| `range` | No | `10.0` | Geo-filter radius in km (only applies when `lat`+`long` provided) |
| `category_ids` | No | — | Comma-separated category IDs — filters to products in ANY of these categories |
| `min_price` | No | — | Minimum price filter (inclusive) |
| `max_price` | No | — | Maximum price filter (inclusive) |

**Validation:**
- `q` must be non-empty after trimming
- `min_price` must be ≥ 0
- `max_price` must be ≥ `min_price` when both are provided

**Example:**
```
GET /api/v1/search?q=apples&lat=12.97&long=77.59&range=5&category_ids=1,3&min_price=50&max_price=200
```

**Response `200`:**
```json
{
  "message": "search results",
  "data": [
    {
      "inv_product_id": 7,
      "product_id":     5,
      "name":           "Organic Apples",
      "price":          120.00,
      "price_unit":     "rs",
      "lat":            12.9716,
      "long":           77.5946,
      "distance_km":    1.2
    }
  ]
}
```

`distance_km` is omitted when no `lat`/`long` is provided.

For full product detail (store address, availability, categories) call:
```
GET /vendor/api/v1/products/{inv_product_id}/detail
```

**Response `400`** — missing `q`:
```json
{ "message": "q is required" }
```

**Response `504`** — search timed out (query took > 5 seconds):
```json
{ "message": "search timed out" }
```

---

---

# Internal Endpoints

These routes are **not exposed through the gateway** and are **not reachable from outside Docker**. Services call each other directly by container name over the shared `spotnearr-network`. No auth middleware is applied; Docker network isolation is the security boundary.

## User Service Internal — `user-service:8080/internal/`

### `POST /internal/users/{userId}/invalidate-refresh`

Called by the Vendor Service after a business account action that should invalidate the user's session. Deletes the user's refresh token from cache.

**Response `200`.**

---

### `GET /internal/claims?product_ids=1,2,3`

Called by the Vendor Service to fetch all claims for a set of inventory product IDs.

**Query:** `product_ids` — comma-separated inventory product IDs.

**Response `200`:** Array of claim objects.

---

### `GET /internal/claims/{claimId}`

Fetch a single claim by ID.

**Response `200`:** Claim object. **`404`** if not found.

---

### `PATCH /internal/claims/{claimId}/status`

Update a claim's status.

**Request body:**
```json
{ "status": "accepted" }
```

`status`: `"accepted"` | `"rejected"`

**Response `200`.**

---

## Vendor Service Internal — `vendor-service:8080/internal/`

### `GET /internal/inventory-products/{id}`

Called by the User Service during claim creation to validate product existence and availability.

**Response `200`:** Inventory product detail object. **`404`** if not found.

---

### `GET /internal/users/{id}/access`

Called by the User Service on login to check if the user has a business account.

**Response `200`:**
```json
{ "business_id": 1, "role": "owner" }
```

**`404`** if the user has no business access.

---

### `POST /internal/businesses/{bizId}/follow`
### `POST /internal/businesses/{bizId}/unfollow`

Called by the User Service when a user follows/unfollows a business (via RabbitMQ consumer). Increments or decrements `followerCount` on the business record.

**Response `200`.**

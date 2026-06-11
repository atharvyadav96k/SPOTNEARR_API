# Spotnearr API Reference

## Services & Base URLs

| Service | Default Port | Purpose |
|---|---|---|
| **User Service** | `8080` | Customer accounts, auth, reviews, claims |
| **Vendor Service** | `8081` | Business owner dashboard — products, inventory, categories |
| **Search Service** | `8082` | Product search (public, no auth) |

All endpoints under `/api/v1/...` unless noted. CORS is enabled on all services.

---

## Response Shape

Every response follows this envelope:
```json
{ "message": "string", "data": <object|array|null> }
```
HTTP status is always set on the response. `data` is omitted when null.

**Common error codes**

| Code | Meaning |
|---|---|
| `400` | Validation failure / missing required field |
| `401` | Missing, invalid, or expired JWT |
| `403` | Valid JWT but insufficient role |
| `404` | Resource not found or not owned by caller |
| `409` | Duplicate (email, slug, etc.) |
| `500` | Server error |

---

## Auth Header

Protected endpoints require:
```
Authorization: Bearer <access_token>
```

**BusinessOnly** endpoints additionally require the token to carry a `business_id` (obtained after logging into the vendor service).

---

# Vendor Service — `:8081`

Business owners register and manage products here. The vendor service has its own auth — separate from the user service.

## Auth

### `POST /api/v1/auth/register`
Create a business account.

**Body:**
```json
{
  "name": "Biz Name",
  "email": "owner@example.com",
  "phone": "+919876543210",
  "password": "Abcd1234",
  "desc": "optional"
}
```
- `password` — must have uppercase, lowercase, digit, min 8 chars
- `phone` — E.164 format

**Response `200`:**
```json
{
  "data": {
    "access_token": "<jwt>",
    "access_token_expires_at": 1749259500,
    "refresh_token": "<jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```
- Access token expires in ~5 minutes
- Refresh token expires in ~30 days

---

### `POST /api/v1/auth/login`
**Body:**
```json
{ "email": "owner@example.com", "password": "Abcd1234" }
```
**Response `200`:** same shape as register.

---

### `POST /api/v1/auth/refresh`
Exchange a refresh token for a new access token.

**Body:**
```json
{ "refresh_token": "<jwt>" }
```
**Response `200`:**
```json
{
  "data": {
    "access_token": "<new_jwt>",
    "access_token_expires_at": 1749259500
  }
}
```

---

## Categories
Auth: **JWT**

### `GET /api/v1/categories/`
List all categories.

**Response `200`:**
```json
{
  "data": [
    { "id": 1, "name": "Fruits", "slug": "fruits" }
  ]
}
```

### `POST /api/v1/categories/`
Auth: **JWT + BusinessOnly**

**Body:**
```json
{ "name": "Dairy", "slug": "dairy" }
```
Both `name` and `slug` must be globally unique.

**Response `201`:**
```json
{ "data": { "id": 3, "name": "Dairy", "slug": "dairy" } }
```

---

## Products
Auth: **JWT + BusinessOnly** on all product endpoints.

### `POST /api/v1/products/`
Create a product.

**Body:**
```json
{
  "name": "Organic Apples",
  "price":    { "value": 120.00, "unit": "rs" },
  "quantity": { "value": 1.0,   "unit": "kg" },
  "desc": "optional",
  "categoryIds": [1, 3],
  "storeIds": [1]
}
```
- `name`, `price`, `quantity`, `categoryIds` — required
- `storeIds` — optional; auto-links product to the first store if omitted

**Response `201`:**
```json
{
  "data": {
    "id": 5,
    "name": "Organic Apples",
    "price": 120.00, "priceUnit": "rs",
    "quantity": 1.0, "quantityUnit": "kg",
    "desc": "...",
    "businessId": 1,
    "categories": [{ "id": 1, "name": "Fruits", "slug": "fruits" }],
    "createdAt": "...", "updatedAt": "..."
  }
}
```

### `GET /api/v1/products/`
List all products for the authenticated business.

**Response `200`:** `data` is an array of product objects (same shape as above).

### `GET /api/v1/products/{productId}`
Get a single product. Returns `404` if it doesn't belong to the caller's business.

### `PATCH /api/v1/products/{productId}`
Update a product. `price` and `quantity` are required even if unchanged.

**Body:**
```json
{
  "name": "Updated Name",
  "price":    { "value": 130.00, "unit": "rs" },
  "quantity": { "value": 1.0,   "unit": "kg" },
  "desc": "optional"
}
```
**Response `200`:** updated product object in `data`.

### `DELETE /api/v1/products/{productId}`
Soft-delete a product. **Response `200`.**

---

## Inventory (Stores)
Auth: **JWT + BusinessOnly** on all inventory endpoints.

### `POST /api/v1/inventory/`
Create a store location.

**Body:**
```json
{
  "name": "Main Branch",
  "streetAddress": "123 Market St",
  "lat": 12.9716,
  "long": 77.5946
}
```
`geoHash` is computed automatically.

**Response `200`:**
```json
{
  "data": {
    "id": 1,
    "name": "Main Branch",
    "streetAddress": "123 Market St",
    "businessId": 1,
    "lat": 12.9716, "long": 77.5946,
    "geoHash": "tdr1u",
    "createdAt": "...", "updatedAt": "..."
  }
}
```

### `GET /api/v1/inventory/`
List all stores for the authenticated business. **Response `200`:** array in `data`.

### `PATCH /api/v1/inventory/{invId}`
Update a store. `lat` and `long` are required (used to recompute `geoHash`).

**Body:**
```json
{
  "name": "Updated Name",
  "address": "New Address",
  "lat": 12.9716,
  "long": 77.5946
}
```
**Response `200`.**

### `GET /api/v1/inventory/{invId}/products`
List all products linked to a store.

**Response `200`:**
```json
{
  "data": [
    { "id": 7, "storeId": 1, "productId": 5, "count": 10, "available": true, "createdAt": "...", "updatedAt": "..." }
  ]
}
```

### `POST /api/v1/inventory/{invId}/products`
Add a product to a store.

**Body:**
```json
{ "productID": 5, "count": 10, "available": true }
```
**Response `200`.**

### `PATCH /api/v1/inventory/{invId}/products/{invProductId}`
Update stock/availability for a product in a store.

**Body:**
```json
{ "count": 20, "available": true }
```
Both fields are optional (omit to leave unchanged).
**Response `200`.**

### `DELETE /api/v1/inventory/{invId}/products/{invProductId}`
Remove a product from a store. **Response `200`.**

---

## Business Profile
Auth: **JWT**

### `POST /api/v1/businesses/register`
Register a business profile after creating an account (call after first login). Returns the business object.

**Body:**
```json
{ "name": "Store Name", "email": "biz@example.com", "phone": "+91...", "desc": "optional" }
```

### `GET /api/v1/businesses/{bizId}/profile`
Get a business profile by ID. **Response `200`:** business object in `data`.

### `PATCH /api/v1/businesses/`
Auth: **JWT + BusinessOnly**. Update name or description.

**Body:** `{ "name": "...", "desc": "..." }` — both optional.

---

# Search Service — `:8082`

No auth required. Public.

## `GET /api/v1/search`

**Query params:**

| Param | Required | Description |
|---|---|---|
| `q` | Yes | Search query string |
| `lat` | No | Latitude |
| `long` | No | Longitude |
| `range` | No | Radius in km (default: `10.0`) |

**Example:** `GET /api/v1/search?q=apples&lat=12.97&long=77.59&range=5`

**Response `200`:**
```json
{
  "data": [
    {
      "id": 7,
      "productName": "Organic Apples",
      "price": 120.00, "priceUnit": "rs",
      "available": true,
      "storeName": "Main Branch",
      "streetAddress": "123 Market St",
      "lat": 12.9716, "long": 77.5946,
      "categories": [{ "id": 1 }]
    }
  ]
}
```

---

# User Service — `:8080`

End-user accounts (customers). Requires Cloudflare Turnstile captcha on register/login.

## Auth

### `POST /api/v1/auth/users/register`
**Headers:** `X-Captcha-Token: <turnstile_token>`

**Body:**
```json
{ "name": "John Doe", "email": "john@example.com", "phone": "+91...", "password": "Abcd1234" }
```
**Response `200`:** `{ "data": { "access_token", "access_token_expires_at", "refresh_token", "refresh_token_expires_at" } }`

### `POST /api/v1/auth/users/login`
**Headers:** `X-Captcha-Token: <turnstile_token>`

**Body:** `{ "email": "...", "password": "..." }`
**Response `200`:** same token shape.

### `POST /api/v1/auth/refresh`
**Body:** `{ "access_token": "...", "refresh_token": "..." }`
**Response `200`:** new `access_token` + updated `access_token_expires_at`.

### `POST /api/v1/auth/reset-request`
**Headers:** `X-Captcha-Token`

**Body:** `{ "email": "..." }` — sends reset email. **Response `200`.**

### `POST /api/v1/auth/reset-password?session=<token>`
**Headers:** `X-Captcha-Token`

**Body:** `{ "password": "NewPass123" }` **Response `200`.**

---

## Users

### `GET /api/v1/users/{userId}/profile`
Auth: **JWT**. Returns the caller's profile (path param is ignored; user is derived from token).

**Response `200`:**
```json
{
  "data": {
    "id": 1, "fullName": "John Doe",
    "email": "john@example.com", "phone": "+91...",
    "isVerifiedEmail": false, "isVerifiedPhone": false,
    "isActive": true, "createdAt": "...", "updatedAt": "..."
  }
}
```

---

## JWT Claims (decoded)

```json
{
  "user_id": 42,
  "business_id": 1,
  "token_type": "access",
  "role": "admin",
  "exp": 1700000000
}
```
- `business_id` is `null` until a business is registered
- `role`: `"user"` | `"admin"` | `"sub-admin"`
- Access token: ~5 min lifetime — refresh before expiry using `refresh_token`

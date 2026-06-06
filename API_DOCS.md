# SPOTNEARR API — Complete Reference

**Base URL:** `http://<host>:<PORT>/api/v1`  
**Default port:** `8080`  
**Framework:** Go + Gorilla Mux  
**Auth:** JWT (HS256), access token in `Authorization: Bearer <token>` header

---

## Global Conventions

### Response Envelope
Every response (success or error) uses this shape:
```json
{
  "message": "string",
  "data": <object|array|null>
}
```
HTTP status code is set directly on the response. `data` is omitted when null.

### Auth Levels
| Level | Header Required |
|---|---|
| **None** | No headers needed |
| **JWT** | `Authorization: Bearer <access_token>` |
| **BusinessOnly** | JWT + token must have `business_id` and role ≠ `"user"` |
| **Captcha** | `X-Captcha-Token: <turnstile_token>` (Cloudflare Turnstile) |

### Rate Limits (per user, per endpoint)
| Profile | Limit |
|---|---|
| Strict | 5 req / min |
| Normal | 10–30 req / min |
| Relaxed | 60 req / min |

---

## 1. Health

### GET `/health`
**Auth:** None  
**Rate limit:** None  
Returns `200 OK` when the server is running.

**Response:**
```
HTTP 200
(empty body)
```

---

## 2. Authentication

### POST `/auth/users/register`
**Auth:** Captcha  
**Rate limit:** None  
Register a new user account.

**Headers:**
```
X-Captcha-Token: <turnstile_token>
Content-Type: application/json
```

**Request body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+919876543210",
  "password": "SecurePass123!"
}
```
- `email` — must be a valid email address
- `phone` — must be a valid phone number
- `password` — validated (min length + complexity enforced server-side)
- `name` — optional but recommended

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "access_token": "<jwt>",
    "access_token_expires_at": 1749259500,
    "refresh_token": "<jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```

---

### POST `/auth/users/login`
**Auth:** Captcha  
**Rate limit:** None  

**Headers:**
```
X-Captcha-Token: <turnstile_token>
Content-Type: application/json
```

**Request body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "access_token": "<jwt>",
    "access_token_expires_at": 1749259500,
    "refresh_token": "<jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```

---

### POST `/auth/reset-request`
**Auth:** Captcha  
**Rate limit:** None  
Send a password-reset session notification to the user's email.

**Headers:**
```
X-Captcha-Token: <turnstile_token>
Content-Type: application/json
```

**Request body:**
```json
{
  "email": "john@example.com"
}
```

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### POST `/auth/reset-password?session=<sessionId>`
**Auth:** Captcha + valid session query param  
**Rate limit:** None  
Reset the user's password. The `session` query parameter is the token received via the reset notification.

**Headers:**
```
X-Captcha-Token: <turnstile_token>
Content-Type: application/json
```

**Query params:**
| Param | Type | Required | Description |
|---|---|---|---|
| `session` | string | Yes | Session token from reset notification |

**Request body:**
```json
{
  "password": "NewSecurePass456!"
}
```

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### POST `/auth/refresh`
**Auth:** None  
**Rate limit:** None  
Exchange a valid refresh token for a new access/refresh token pair.

**Request body:**
```json
{
  "access_token": "<expired_or_current_access_token>",
  "refresh_token": "<valid_refresh_token>"
}
```

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "access_token": "<new_jwt>",
    "access_token_expires_at": 1749259500,
    "refresh_token": "<jwt>",
    "refresh_token_expires_at": 1751847900
  }
}
```
> Note: `refresh_token` is unchanged (same token passed in). `refresh_token_expires_at` reflects its original expiry from its JWT claims.

---

### GET `/auth/`
**Auth:** JWT  
**Rate limit:** None  
Verify that the current access token is valid. Returns `200` with no data if valid.

**Response `200`:**
```
(empty body)
```

---

### POST `/auth/logout-all-devices`
**Auth:** JWT  
**Rate limit:** None  
Invalidate all refresh tokens for the authenticated user (forces re-login on all devices).

**Response `200`:**
```json
{
  "message": "..."
}
```

---

## 3. Users

### GET `/users/{userId}/profile`
**Auth:** JWT  
**Rate limit:** Normal (10 req/min)  
Get profile of any user by ID. The server uses the userId from the JWT, not the path param.

**Path params:**
| Param | Type | Description |
|---|---|---|
| `userId` | uint | Target user ID |

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 1,
    "fullName": "John Doe",
    "email": "john@example.com",
    "phone": "+919876543210",
    "isVerifiedEmail": false,
    "isVerifiedPhone": false,
    "isActive": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```
> Note: `passwordHash` and `pushToken` are never returned.

---

## 4. Businesses

### POST `/businesses/register`
**Auth:** JWT  
**Rate limit:** Strict (5 req/min)  
Register a new business linked to the authenticated user.

**Request body:**
```json
{
  "name": "Baaner Store",
  "email": "store@baaner.com",
  "phone": "+919876543210",
  "desc": "A description of the business"
}
```

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 1,
    "businessName": "Baaner Store",
    "email": "store@baaner.com",
    "phone": "+919876543210",
    "description": "A description of the business",
    "isActive": true,
    "verifiedBusiness": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```
> After registering a business, the user must log in again (or refresh) to get a token with `business_id` embedded — business-only endpoints require a token that includes `business_id`.

---

### GET `/businesses/{bizId}/profile`
**Auth:** JWT  
**Rate limit:** Normal (30 req/min)

**Path params:**
| Param | Type | Description |
|---|---|---|
| `bizId` | uint | Business ID |

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 1,
    "businessName": "Baaner Store",
    "email": "store@baaner.com",
    "phone": "+919876543210",
    "description": "...",
    "isActive": true,
    "verifiedBusiness": false,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### PATCH `/businesses/`
**Auth:** JWT + BusinessOnly  
**Rate limit:** Normal (30 req/min)  
Update the authenticated user's business profile.

**Request body:**
```json
{
  "name": "Updated Business Name",
  "desc": "Updated description"
}
```
Both fields are optional (omit to leave unchanged).

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### DELETE `/businesses/`
**Auth:** JWT + BusinessOnly  
**Rate limit:** Strict (5 req/min)  
Delete the authenticated business. (Stub — returns `200 OK` currently.)

---

### GET `/businesses/inventories`
**Auth:** JWT + BusinessOnly  
**Rate limit:** Normal (30 req/min)  
List all stores/inventories belonging to the authenticated business.

**Response `200`:**
```json
{
  "message": "...",
  "data": [
    {
      "id": 1,
      "name": "Main Branch",
      "streetAddress": "123 Market St",
      "businessId": 1,
      "lat": 12.9716,
      "long": 77.5946,
      "geoHash": "tdr1u",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 5. Inventory (Stores)

All inventory endpoints require **JWT + BusinessOnly**.

### POST `/inventory/`
**Rate limit:** Strict (5 req/min)  
Create a new store/inventory location.

**Request body:**
```json
{
  "name": "Main Branch",
  "streetAddress": "123 Market St, Bengaluru",
  "lat": 12.9716,
  "long": 77.5946
}
```
All fields are required. `geoHash` is computed automatically.

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 1,
    "name": "Main Branch",
    "streetAddress": "123 Market St, Bengaluru",
    "businessId": 1,
    "lat": 12.9716,
    "long": 77.5946,
    "geoHash": "tdr1u",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### PATCH `/inventory/{invId}`
**Rate limit:** Normal (30 req/min)  
Update store details. `lat` and `long` are **required** even if unchanged (used to recompute geoHash).

**Path params:**
| Param | Type | Description |
|---|---|---|
| `invId` | uint | Inventory/store ID |

**Request body:**
```json
{
  "name": "Updated Branch Name",
  "address": "456 New Street",
  "lat": 12.9716,
  "long": 77.5946
}
```

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### DELETE `/inventory/{invId}`
**Rate limit:** Strict (5 req/min)  
Delete a store. (Stub — returns `200 OK` currently.)

---

### GET `/inventory/{invId}/products`
**Rate limit:** Normal (30 req/min)  
List all products in a specific store.

**Path params:**
| Param | Type | Description |
|---|---|---|
| `invId` | uint | Inventory/store ID |

**Response `200`:**
```json
{
  "message": "...",
  "data": [
    {
      "id": 1,
      "storeId": 1,
      "productId": 5,
      "count": 10,
      "available": true,
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### POST `/inventory/{invId}/products`
**Rate limit:** Normal (30 req/min)  
Add a product to a store inventory. Uses **query parameters** (not JSON body).

**Path params:**
| Param | Type | Description |
|---|---|---|
| `invId` | uint | Inventory/store ID |

**Query params:**
| Param | Type | Required | Description |
|---|---|---|---|
| `productID` | uint | Yes | ID of product to add |
| `count` | int | No | Stock count |
| `available` | bool | No | `"true"` or `"false"` (default: false) |

**Example:** `POST /inventory/1/products?productID=5&count=10&available=true`

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### DELETE `/inventory/{invId}/products`
**Rate limit:** Strict (5 req/min)  
Remove a product from a store inventory. (Stub — returns `200 OK` currently.)

---

## 6. Products

All product endpoints require **JWT + BusinessOnly**.

### POST `/products/`
**Rate limit:** Strict (5 req/min)  
Create a new product for the authenticated business.

**Request body:**
```json
{
  "name": "Organic Apples",
  "price": {
    "value": 120.00,
    "unit": "INR"
  },
  "quantity": {
    "value": 1.0,
    "unit": "kg"
  },
  "desc": "Fresh organic apples from Himachal",
  "categoryIds": [1, 3],
  "storeIds": [1, 2]
}
```
- `name`, `price`, `quantity`, `categoryIds` — required
- `storeIds` — optional; links product to specific stores immediately
- `desc` — optional

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 5,
    "name": "Organic Apples",
    "price": 120.00,
    "priceUnit": "INR",
    "quantity": 1.0,
    "quantityUnit": "kg",
    "desc": "Fresh organic apples from Himachal",
    "businessId": 1,
    "categories": [
      { "id": 1, "name": "Fruits", "slug": "fruits" }
    ],
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### GET `/products/{productId}`
**Rate limit:** Relaxed (60 req/min)  
Get a specific product by ID. Only returns products belonging to the authenticated business.

**Path params:**
| Param | Type | Description |
|---|---|---|
| `productId` | uint | Product ID |

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 5,
    "name": "Organic Apples",
    "price": 120.00,
    "priceUnit": "INR",
    "quantity": 1.0,
    "quantityUnit": "kg",
    "desc": "...",
    "businessId": 1,
    "categories": [
      { "id": 1, "name": "Fruits", "slug": "fruits" }
    ],
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### PATCH `/products/{productId}`
**Rate limit:** Normal (30 req/min)  
Update a product. `price` and `quantity` are required even if unchanged.

**Path params:**
| Param | Type | Description |
|---|---|---|
| `productId` | uint | Product ID |

**Request body:**
```json
{
  "name": "Updated Apples",
  "price": {
    "value": 130.00,
    "unit": "INR"
  },
  "quantity": {
    "value": 1.0,
    "unit": "kg"
  },
  "desc": "Updated description"
}
```

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### DELETE `/products/{productId}`
**Rate limit:** Strict (5 req/min)  
Soft-delete a product.

**Path params:**
| Param | Type | Description |
|---|---|---|
| `productId` | uint | Product ID |

**Response `200`:**
```json
{
  "message": "..."
}
```

---

### GET `/products/nearby`
**Rate limit:** Relaxed (60 req/min)  
Get nearby products. (Stub — returns `200 OK` currently. Use `/search` for live geo search.)

---

## 7. Search

### GET `/search`
**Auth:** None  
**Rate limit:** None  
Full-text product search with optional geo filtering.

**Query params:**
| Param | Type | Required | Description |
|---|---|---|---|
| `search` | string | Yes | Search term (tokenized against product names/descriptions) |
| `lat` | float64 | No | Latitude for geo filtering |
| `long` | float64 | No | Longitude for geo filtering |
| `range` | float64 | No | Search radius in km (default: `5.0`) |

**Example:** `GET /search?search=apples&lat=12.9716&long=77.5946&range=10`

**Response `200`:**
```json
{
  "message": "...",
  "data": [
    {
      "id": 5,
      "name": "Organic Apples",
      "price": 120.00,
      "priceUnit": "INR",
      "quantity": 1.0,
      "quantityUnit": "kg",
      "desc": "...",
      "businessId": 1,
      "categories": [...],
      "createdAt": "...",
      "updatedAt": "..."
    }
  ]
}
```

---

## 8. Claims

All claim endpoints require **JWT**.

### GET `/claims`
**Rate limit:** Normal (10 req/min)  
Get all claims made by the authenticated user.

**Response `200`:**
```json
{
  "message": "...",
  "data": [
    {
      "id": 1,
      "userId": 42,
      "inventoryProductId": 7,
      "status": "pending",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

### POST `/claims/{invProductId}`
**Rate limit:** Normal (10 req/min)  
Claim a product from a store (by its inventory-product join ID).

**Path params:**
| Param | Type | Description |
|---|---|---|
| `invProductId` | uint | ID from the inventory_products table |

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 1,
    "userId": 42,
    "inventoryProductId": 7,
    "status": "pending",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

---

### DELETE `/claims/{claimId}`
**Rate limit:** Normal (10 req/min)  
Remove/cancel a claim owned by the authenticated user.

**Path params:**
| Param | Type | Description |
|---|---|---|
| `claimId` | uint | Claim ID |

**Response `200`:**
```json
{
  "message": "..."
}
```

---

## 9. Reviews

All review endpoints require **JWT**.

### Review body (POST / PUT)
```json
{
  "targetId": 1,
  "stars": 4,
  "comment": "Great product!"
}
```
- `targetId` — ID of the entity being reviewed
- `stars` — integer 0–255 (validated by business logic)
- `comment` — optional text

### Delete body (DELETE)
```json
{
  "targetId": 1,
  "stars": 0,
  "comment": ""
}
```
Only `targetId` is used for deletion; `stars`/`comment` are ignored.

---

### Business Reviews

| Method | Path | Description | Rate |
|---|---|---|---|
| `GET` | `/review/businesses?id={businessId}` | Get all reviews for a business | Relaxed |
| `POST` | `/review/businesses` | Add a review for a business | Normal |
| `PUT` | `/review/businesses` | Update your review for a business | Normal |
| `DELETE` | `/review/businesses` | Delete your review for a business | Normal |

**GET query params:**
| Param | Type | Required |
|---|---|---|
| `id` | uint | Yes |

**GET Response `200`:**
```json
{
  "message": "...",
  "data": [
    {
      "id": 1,
      "userId": 42,
      "targetType": "business",
      "targetId": 1,
      "stars": 4,
      "comment": "Great service!",
      "createdAt": "...",
      "updatedAt": "..."
    }
  ]
}
```

---

### Product Reviews

| Method | Path | Description | Rate |
|---|---|---|---|
| `GET` | `/review/products?id={productId}` | Get all reviews for a product | Relaxed |
| `POST` | `/review/products` | Add a review for a product | Normal |
| `PUT` | `/review/products` | Update your review for a product | Normal |
| `DELETE` | `/review/products` | Delete your review for a product | Normal |

**GET query params:**
| Param | Type | Required |
|---|---|---|
| `id` | uint | Yes |

---

### Offer Reviews

| Method | Path | Description | Rate |
|---|---|---|---|
| `GET` | `/review/offers/{offerId}` | Get all reviews for an offer | Relaxed |
| `POST` | `/review/offers` | Add a review for an offer | Normal |
| `PUT` | `/review/offers` | Update your review for an offer | Normal |
| `DELETE` | `/review/offers` | Delete your review for an offer | Normal |

---

### Spotlight Reviews

| Method | Path | Description | Rate |
|---|---|---|---|
| `GET` | `/review/spotlights/{spotlightId}` | Get all reviews for a spotlight | Relaxed |
| `POST` | `/review/spotlights` | Add a review for a spotlight | Normal |
| `PUT` | `/review/spotlights` | Update your review for a spotlight | Normal |
| `DELETE` | `/review/spotlights` | Delete your review for a spotlight | Normal |

---

## 10. Categories

### GET `/categories/`
**Auth:** JWT  
**Rate limit:** Normal (30 req/min)  
List all product categories.

**Response `200`:**
```json
{
  "message": "...",
  "data": [
    {
      "id": 1,
      "name": "Fruits",
      "slug": "fruits"
    },
    {
      "id": 2,
      "name": "Vegetables",
      "slug": "vegetables"
    }
  ]
}
```

---

### POST `/categories/`
**Auth:** JWT + BusinessOnly  
**Rate limit:** Strict (5 req/min)  
Create a new product category.

**Request body:**
```json
{
  "name": "Dairy",
  "slug": "dairy"
}
```
Both `name` and `slug` must be unique.

**Response `200`:**
```json
{
  "message": "...",
  "data": {
    "id": 3,
    "name": "Dairy",
    "slug": "dairy"
  }
}
```

---

## 11. Offers

All offer endpoints require **JWT**. Business write operations also require **BusinessOnly**.

> **Note:** Offer creation/update/delete are stubs that return `200 OK` with no data. Only `GET /offers/nearby` routing is active.

| Method | Path | Auth | Rate | Status |
|---|---|---|---|---|
| `GET` | `/offers/nearby` | JWT | Relaxed | Stub |
| `POST` | `/offers` | JWT + BusinessOnly | Strict | Stub |
| `PUT` | `/offers/{offerID}` | JWT + BusinessOnly | Normal | Stub |
| `DELETE` | `/offers/{offerID}` | JWT + BusinessOnly | Strict | Stub |

---

## Data Models Reference

### User
```json
{
  "id": 1,
  "fullName": "John Doe",
  "email": "john@example.com",
  "phone": "+919876543210",
  "isVerifiedEmail": false,
  "isVerifiedPhone": false,
  "isActive": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### TokenResponse
```json
{
  "access_token": "<jwt>",
  "access_token_expires_at": 1749259500,
  "refresh_token": "<jwt>",
  "refresh_token_expires_at": 1751847900
}
```
- **`access_token_expires_at`** — Unix timestamp (seconds) when the access token expires (~5 minutes from issuance)
- **`refresh_token_expires_at`** — Unix timestamp (seconds) when the refresh token expires (~30 days from issuance)
- Use `POST /auth/refresh` to exchange a refresh token for a new access token before expiry

### JWT Claims (decoded)
```json
{
  "user_id": 42,
  "business_id": 1,
  "token_type": "access",
  "role": "admin",
  "exp": 1700000000,
  "iat": 1700000000,
  "nbf": 1700000000
}
```
- `business_id` is `null` for regular users until they register a business
- `role`: `"user"` | `"admin"` | `"sub-admin"`

### Business
```json
{
  "id": 1,
  "businessName": "Baaner Store",
  "email": "store@baaner.com",
  "phone": "+919876543210",
  "description": "...",
  "isActive": true,
  "verifiedBusiness": false,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Store (Inventory)
```json
{
  "id": 1,
  "name": "Main Branch",
  "streetAddress": "123 Market St",
  "businessId": 1,
  "lat": 12.9716,
  "long": 77.5946,
  "geoHash": "tdr1unz6g",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Product
```json
{
  "id": 5,
  "name": "Organic Apples",
  "price": 120.00,
  "priceUnit": "INR",
  "quantity": 1.0,
  "quantityUnit": "kg",
  "desc": "Fresh organic apples",
  "businessId": 1,
  "categories": [
    { "id": 1, "name": "Fruits", "slug": "fruits" }
  ],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### InventoryProduct (Store-Product link)
```json
{
  "id": 7,
  "storeId": 1,
  "productId": 5,
  "count": 10,
  "available": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Claim
```json
{
  "id": 1,
  "userId": 42,
  "inventoryProductId": 7,
  "status": "pending",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```
- `status`: `"pending"` | `"accepted"` | `"rejected"`

### Review
```json
{
  "id": 1,
  "userId": 42,
  "targetType": "business",
  "targetId": 1,
  "stars": 4,
  "comment": "Great service!",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```
- `targetType`: `"business"` | `"product"` | `"offer"` | `"spotlight"`

### Category
```json
{
  "id": 1,
  "name": "Fruits",
  "slug": "fruits"
}
```

---

## Error Responses

| HTTP Status | Meaning |
|---|---|
| `400 Bad Request` | Missing required field, invalid format, or validation failure |
| `401 Unauthorized` | Missing or invalid JWT, or expired access token |
| `403 Forbidden` | Valid JWT but insufficient role (e.g. non-business user hitting BusinessOnly endpoint) |
| `404 Not Found` | Resource does not exist or does not belong to caller |
| `429 Too Many Requests` | Rate limit exceeded |
| `500 Internal Server Error` | Unexpected server error |

**Error body:**
```json
{
  "message": "descriptive error message"
}
```

---

## Authentication Flow (Step by Step)

```
1. Register       POST /auth/users/register  (requires Captcha)
                  → receives { access_token, access_token_expires_at,
                               refresh_token, refresh_token_expires_at }

2. Login          POST /auth/users/login      (requires Captcha)
                  → receives { access_token, access_token_expires_at,
                               refresh_token, refresh_token_expires_at }

3. API calls      Authorization: Bearer <access_token>
                  (access token valid for 5 min; check access_token_expires_at)

4. Refresh        POST /auth/refresh          (no auth required)
                  body: { access_token, refresh_token }
                  → receives new { access_token, access_token_expires_at,
                                   refresh_token, refresh_token_expires_at }
                  (refresh token valid for 30 days; check refresh_token_expires_at)

5. Register biz   POST /businesses/register   (requires JWT)
                  → creates business linked to user
                  → must refresh/re-login to get token with business_id

6. Biz endpoints  Authorization: Bearer <access_token_with_business_id>
                  Token must have business_id != null AND role != "user"
```

---

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `JWT_SECRET` | Yes | — | Secret key for JWT signing |
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `CAPTCHA_URL` | Yes | — | Cloudflare Turnstile verification URL |
| `CAPTCHA_SECRET_KEY` | Yes | — | Cloudflare Turnstile secret key |
| `CACHE_URL` | No | `redis://localhost:6379` | Redis connection URL |
| `CACHE_PASSWORD` | No | `` | Redis password |
| `PORT` | No | `8080` | Server port |

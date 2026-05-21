# Spotnearr API Reference

All endpoints return JSON in the following envelope:

```json
{
  "success": true,
  "message": "human-readable description",
  "data": { ... },
  "error": null
}
```

On failure, `success` is `false`, `data` is `null`, and `error` contains details.
Validation failures return **HTTP 422** with a structured error:

```json
{
  "success": false,
  "message": "Validation Error",
  "error": [
    { "field": "email_or_phone", "message": "at least one of email or phone is required" }
  ]
}
```

---

## Authentication

### POST `/auth/user/register`
Register a new user.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `full_name` | string | yes | |
| `email` | string | no* | *at least one of email or phone required |
| `phone` | string | no* | |
| `password` | string | yes | stored as hash |
| `role` | string | yes | `customer` \| `business` \| `admin` |

**Response 201** → `UserResponse`

---

### POST `/auth/user/login`
Authenticate an existing user.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `email` | string | no* | *at least one required |
| `phone` | string | no* | |
| `password` | string | yes | |

**Response 200** → `UserResponse`

---

### POST `/auth/business/register`
Register a new business account.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `owner_id` | uuid | yes | |
| `name` | string | yes | |
| `category_id` | uuid | yes | |
| `email` | string | no | |
| `phone` | string | no | |
| `website` | string | no | |
| `description` | string | no | |

**Response 201** → `BusinessResponse`

---

### POST `/auth/business/login`
Authenticate a business owner. Looks up the user by email and verifies the account has `role = business` and is verified.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `email` | string | yes |

**Response 200** → `UserResponse`

---

## User Profile

### GET `/users/profile/profile-info`
Get the authenticated user's profile.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `UserResponse`

---

### DELETE `/users/profile/profile-image/delete-profile-image`
Remove the user's avatar.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

## Business Profile

### GET `/business/profile/profile-info`
Get a business profile by owner.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `owner_id` | uuid | yes |

**Response 200** → `BusinessResponse`

---

### POST `/business/profile/profile-image/set-profile-image-in-db`
Attach a logo URL to a business after upload. `logo_url` must be non-empty.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `id` | uuid | yes | business id |
| `logo_url` | string | yes | URL of uploaded image |

**Response 200** → `null`

---

### DELETE `/business/profile/profile-image/delete-profile-image`
Remove a business logo or cover image.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

## Business Categories

### POST `/business/categories/add_categories`
Create a new business category.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `name` | string | yes |
| `icon_url` | string | no |

**Response 201** → `BusinessCategoryResponse`

---

### GET `/business/categories/get_categories`
List all business categories.

**Body** — none required

**Response 200** → `[]BusinessCategoryResponse`

---

## Business Locations

### POST `/business/location/create_location`
Add a branch location for a business. Validates all required address and hours fields.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `business_id` | uuid | yes | |
| `branch_name` | string | yes | |
| `address_line1` | string | yes | |
| `address_line2` | string | no | |
| `city` | string | yes | |
| `state` | string | yes | |
| `pin_code` | string | yes | |
| `latitude` | float64 | yes | |
| `longitude` | float64 | yes | |
| `is_main` | bool | no | default false |
| `opening_time` | string | yes | e.g. `"09:00"` |
| `closing_time` | string | yes | e.g. `"21:00"` |
| `working_days` | []string | yes | min 1 day |

**Response 201** → `BusinessLocationResponse`

---

### GET `/business/location/get_location`
Get a specific location by ID.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `BusinessLocationResponse`

---

### PUT `/business/location/update_location`
Update a branch location.

**Body** — same fields as create; `id` required.

**Response 200** → `BusinessLocationResponse`

---

### DELETE `/business/location/delete_location`
Delete a branch location.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

## Products

### POST `/business/products/new_product`
Create a product. Validates category, name, and unit.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `category_id` | uuid | yes | |
| `name` | string | yes | |
| `description` | string | no | |
| `unit` | string | yes | e.g. `"kg"`, `"piece"` |
| `is_available` | bool | no | default false |
| `is_active` | bool | no | default false |
| `tags` | []string | no | |

**Response 201** → `ProductResponse`

---

### GET `/business/products/get_products`
Get a product by ID.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `ProductResponse`

---

### GET `/business/products/get_business_products`
List all products for a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]ProductResponse`

---

### PUT `/business/products/update_products`
Update a product.

**Body** — same fields as create; `id` required.

**Response 200** → `ProductResponse`

---

### DELETE `/business/products/delete_products`
Delete a product.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

## Product Inventory

### POST `/business/inventory/create`
Create inventory records for one or more products at a location.
Validates each item: `product_id`, `location_id`, and `price > 0` required.

**Body** — JSON array of inventory objects:

```json
[
  {
    "product_id": "uuid",
    "location_id": "uuid",
    "price": 1500,
    "discounted_price": 1200,
    "stock": 50,
    "is_available": true
  }
]
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `product_id` | uuid | yes | |
| `location_id` | uuid | yes | |
| `price` | int64 | yes | in smallest unit (paise/cents); must be > 0 |
| `discounted_price` | int64 | no | 0 means no discount |
| `stock` | int | no | |
| `is_available` | bool | no | |

**Response 201** → `null`

---

### GET `/business/inventory/get`
Get a specific inventory record.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `ProductInventoryResponse`

---

### GET `/business/inventory/get_by_business`
List all inventory records for a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]ProductInventoryResponse`

---

### PUT `/business/inventory/update`
Update an inventory record.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `price` | int64 | no | must be > 0 if provided |
| `discounted_price` | int64 | no | |
| `stock` | int | no | |
| `is_available` | bool | no | |

**Response 200** → `ProductInventoryResponse`

---

### DELETE `/business/inventory/delete`
Delete an inventory record.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

## Spotlights (Business)

### POST `/business/spotlights/create_spotlight`
Create a spotlight post. Validates `business_id`, `type`, `title`, `media_url`, `media_type`.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `business_id` | uuid | yes | |
| `type` | string | yes | `product` \| `offer` \| `general` |
| `status` | string | no | `draft` \| `published` \| `expired` |
| `title` | string | yes | |
| `description` | string | no | |
| `media_url` | string | yes | |
| `media_type` | string | yes | `image` \| `video` |
| `thumbnail_url` | string | no | |
| `product_id` | uuid | no | required when type = `product` |
| `offer_id` | uuid | no | required when type = `offer` |
| `expires_at` | datetime | no | ISO 8601 |

**Response 201** → `SpotlightResponse`

---

### GET `/business/spotlights/get_spotlights`
Get a spotlight by ID.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `SpotlightResponse`

---

### PUT `/business/spotlights/update_spotlight`
Update a spotlight.

**Body** — same fields as create; `id` required.

**Response 200** → `SpotlightResponse`

---

### DELETE `/business/spotlights/delete_spotlight`
Delete a spotlight.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

## Plans & Subscriptions

### POST `/business/plans/add_plans`
Create a subscription plan (admin only).

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | |
| `tier` | string | yes | enum: plan tier |
| `price` | float64 | yes | |
| `description` | string | no | |
| `is_active` | bool | no | |

**Response 201** → `PlanResponse`

---

### GET `/business/plans/get_business_plans`
List all available plans.

**Body** — none required

**Response 200** → `[]PlanResponse`

---

### GET `/business/plans/get_active_plan`
Get the active subscription for a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `BusinessSubscriptionResponse`

---

### POST `/business/plans/set_business_active_plan`
Assign a plan to a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `business_id` | uuid | yes |
| `plan_id` | uuid | yes |

**Response 200** → `BusinessSubscriptionResponse`

---

## Transactions

### POST `/business/transactions/create_order`
Create a payment order.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `business_id` | uuid | yes |
| `plan_id` | uuid | yes |
| `amount` | int64 | yes |

**Response 201** → order object

---

### POST `/business/transactions/complete_order`
Mark a payment order as complete after payment gateway callback.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `order_id` | string | yes |
| `payment_id` | string | yes |
| `signature` | string | yes |

**Response 200** → `null`

---

## User — Social: Follow

### POST `/users/follows/follow-user`
Follow another user. Validates `follower_id` and `following_id`.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `follower_id` | uuid | yes |
| `following_id` | uuid | yes |

**Response 201** → `null`

---

### POST `/users/follows/unfollow-user`
Unfollow a user.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `follower_id` | uuid | yes |
| `following_id` | uuid | yes |

**Response 201** → `null`

---

### GET `/users/follows/get-followers`
Get all followers of a user.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]UserFollowResponse`

```json
{ "follower_id": "uuid", "following_id": "uuid", "created_at": "datetime" }
```

---

### GET `/users/follows/get-following-users`
Get all users a user is following.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]UserFollowResponse`

---

### POST `/users/follows/follow-business`
Follow a business. Validates `user_id` and `business_id`.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 201** → `null`

---

### POST `/users/follows/unfollow-business`
Unfollow a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 200** → `null`

---

### GET `/users/follows/get-following-businesses`
Get all businesses a user follows.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]FollowBusinessResponse`

```json
{ "business_id": "uuid", "created_at": "datetime" }
```

---

## User — Social: Product Interactions

### POST `/users/products/save`
Save a product to the user's list. Validates `user_id` and `product_id`.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `product_id` | uuid | yes |

**Response 201** → `null`

---

### POST `/users/products/unsave`
Remove a product from saved list.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `product_id` | uuid | yes |

**Response 201** → `null`

---

### GET `/users/products/get-saved`
Get all products saved by a user.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |

**Response 200** → `[]SavedProductResponse`

```json
{ "product_id": "uuid", "created_at": "datetime" }
```

---

### POST `/users/likes/product-like`
Like a product. Validates `user_id` and `product_id`.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `product_id` | uuid | yes |

**Response 201** → `null`

---

### POST `/users/likes/product-unlike`
Remove a product like.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `product_id` | uuid | yes |

**Response 200** → `null`

---

### GET `/users/likes/get-liked-products`
Get all products liked by a user.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]LikedProductResponse`

```json
{ "product_id": "uuid", "created_at": "datetime" }
```

---

### POST `/users/likes/business-like`
Like a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 201** → `null`

---

### POST `/users/likes/business-dislike`
Remove a business like.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 200** → `null`

---

### GET `/users/products/near_by_products`
Get active products sold by businesses within 5 km of the given coordinates.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `latitude` | float64 | yes |
| `longitude` | float64 | yes |

**Response 200** → `[]ProductResponse`

---

### GET `/users/products/search_proudcts`
Full-text search across product names and descriptions (case-insensitive).

**Body**
| Field | Type | Required |
|-------|------|----------|
| `query` | string | yes |

**Response 200** → `[]ProductResponse`

---

## User — Spotlights

### GET `/users/spotlights/followed_business_spotlights`
Get spotlight feed from businesses the user follows.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]SpotlightResponse`

---

### GET `/users/spotlights/get_spotlight_by_location`
Get published, non-expired spotlights from businesses within 10 km of the given coordinates.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `latitude` | float64 | yes |
| `longitude` | float64 | yes |

**Response 200** → `[]SpotlightResponse`

---

### GET `/users/spotlights/get_liked_spotlights`
Get spotlights the user has liked/saved.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]SpotlightResponse`

---

### GET `/users/spotlights/get-saved`
Get spotlights saved by the user.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]SpotlightResponse`

---

## Spotlight Comments

### POST `/users/spotlights/comments/create`
Post a comment on a spotlight. Validates `spotlight_id`, `user_id`, and `content`.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `spotlight_id` | uuid | yes | |
| `user_id` | uuid | yes | |
| `parent_id` | uuid | no | for threaded replies |
| `content` | string | yes | non-empty |

**Response 201** → `SpotlightCommentResponse`

```json
{
  "id": "uuid",
  "spotlight_id": "uuid",
  "user_id": "uuid",
  "parent_id": "uuid or null",
  "content": "string",
  "created_at": "datetime"
}
```

---

### PUT `/users/spotlights/comments/update`
Edit a comment. Validates `spotlight_id`, `user_id`, and `content`.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `spotlight_id` | uuid | yes |
| `user_id` | uuid | yes |
| `content` | string | yes |

**Response 200** → `SpotlightCommentResponse`

---

### DELETE `/users/spotlights/comments/delete`
Soft-delete a comment.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `null`

---

### GET `/users/spotlights/comments/get_by_spotlight`
Get all comments for a spotlight.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |

**Response 200** → `[]SpotlightCommentResponse`

---

## User — Claims

### POST `/users/claims/claim-product`
Reserve a product from a specific inventory entry. Stock is decremented immediately. Returns HTTP 409 if out of stock.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `user_id` | uuid | yes | |
| `business_id` | uuid | yes | |
| `inventory_id` | uuid | yes | ties the claim to a specific location + price |
| `quantity` | int | yes | must be > 0 |
| `note` | string | no | optional message to the business |

**Response 201** → `ProductClaimResponse`

---

### GET `/users/claims/get-my-claims`
List all claims made by a user.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |

**Response 200** → `[]ProductClaimResponse`

---

### PUT `/users/claims/cancel-claim`
Cancel a pending claim. Restores the reserved stock.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `user_id` | uuid | yes |

**Response 200** → `null`

> Returns 403 if `user_id` doesn't match the claim owner. Returns 400 if claim is not `pending`.

---

### PUT `/users/claims/mark-received`
Mark an accepted claim as completed after the user picks up the product.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `user_id` | uuid | yes |

**Response 200** → `null`

> Returns 400 if the claim is not in `accepted` status.

---

## User — Liked Businesses

### GET `/users/likes/get-liked-businesses`
Get all businesses a user has liked.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |

**Response 200** → `[]LikedBusinessResponse`

```json
{ "business_id": "uuid", "created_at": "datetime" }
```

---

## User — Offers

### GET `/users/offers/get-offers-by-business`
Get active, non-expired offers for a specific business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `business_id` | uuid | yes |

**Response 200** → `[]OfferResponse`

---

## User — Reviews

### POST `/users/reviews/add-review`
Submit a review for a business. One review per user per business (unique constraint).

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `business_id` | uuid | yes | |
| `user_id` | uuid | yes | |
| `rating` | int | yes | 1–5 |
| `comment` | string | no | |

**Response 201** → `ReviewResponse`

---

### GET `/users/reviews/get-reviews`
Get all reviews for a business, newest first.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `business_id` | uuid | yes |

**Response 200** → `[]ReviewResponse`

---

## User — Notifications

### GET `/users/notifications/get-notifications`
Get all notifications for a user, newest first.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `user_id` | uuid | yes |

**Response 200** → `[]NotificationResponse`

---

### PUT `/users/notifications/mark-read`
Mark notifications as read. If `id` is provided, marks that specific notification; otherwise marks all unread notifications for the user.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `user_id` | uuid | yes | |
| `id` | uuid | no | omit to mark all as read |

**Response 200** → `null`

---

## Business — Claims

### GET `/business/claims/get-incoming-claims`
List all product claims received by a business.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `business_id` | uuid | yes |

**Response 200** → `[]ProductClaimResponse`

---

### PUT `/business/claims/accept-claim`
Accept a pending claim. No stock change (already reserved on creation).

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 200** → `null`

> Returns 403 if `business_id` doesn't match the claim. Returns 400 if claim is not `pending`.

---

### PUT `/business/claims/reject-claim`
Reject a pending claim. Restores the reserved stock.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 200** → `null`

---

## Business — Offers

### POST `/business/offers/create_offer`
Create a discount offer for the business.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `business_id` | uuid | yes | |
| `title` | string | yes | |
| `description` | string | no | |
| `discount_type` | string | yes | `percentage` \| `flat` |
| `discount_value` | int64 | yes | must be > 0 |
| `min_order_value` | int64 | no | default 0 |
| `coupon_code` | string | no | |
| `banner_url` | string | no | |
| `is_active` | bool | no | default true |
| `starts_at` | datetime | yes | ISO 8601 |
| `expires_at` | datetime | yes | ISO 8601 |

**Response 201** → `OfferResponse`

---

### GET `/business/offers/get_business_offers`
List all offers for a business (active and inactive).

**Body**
| Field | Type | Required |
|-------|------|----------|
| `business_id` | uuid | yes |

**Response 200** → `[]OfferResponse`

---

### PUT `/business/offers/update_offer`
Update an existing offer. Validates all required fields. `id` and `business_id` are used to scope the update.

**Body** — same fields as create; `id` required.

**Response 200** → `OfferResponse`

---

### DELETE `/business/offers/delete_offer`
Delete an offer.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `id` | uuid | yes |
| `business_id` | uuid | yes |

**Response 200** → `null`

---

## Bridge Service

### POST `/bridge-service/image-upload/upload-signed-url`
Get a GCS signed URL to upload an image directly from the client.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `file_name` | string | yes | |
| `content_type` | string | yes | e.g. `image/jpeg` |

**Response 200** → `{ "url": "string", "object_name": "string" }`

---

### POST `/bridge-service/image-upload/download-signed-url`
Get a signed URL to read a private GCS object.

**Body**
| Field | Type | Required |
|-------|------|----------|
| `object_name` | string | yes |

**Response 200** → `{ "url": "string" }`

---

### POST `/bridge-service/image-upload/on-image-upload`
GCS event trigger — called automatically when an image finishes uploading. Records the media entry in the database.

**Body** — GCS Pub/Sub event (internal, not client-facing)

---

### POST `/bridge-service/push_notification`
Send a push notification to a user device.

**Body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `user_id` | uuid | yes | |
| `title` | string | yes | |
| `body` | string | no | |
| `ref_id` | uuid | no | spotlight / offer / product id |
| `ref_type` | string | no | `spotlight` \| `offer` \| `product` |

**Response 200** → `null`

---

## Response Shapes Reference

### `UserResponse`
```json
{
  "id": "uuid",
  "full_name": "string",
  "email": "string or null",
  "phone": "string or null",
  "role": "customer|business|admin",
  "avatar_url": "string or null",
  "is_verified": false
}
```

### `BusinessResponse`
```json
{
  "id": "uuid",
  "owner_id": "uuid",
  "name": "string",
  "description": "string or null",
  "category_id": "uuid",
  "email": "string or null",
  "phone": "string or null",
  "website": "string or null",
  "logo_url": "string or null",
  "cover_url": "string or null",
  "status": "pending|active|suspended",
  "is_verified": false,
  "rating": 0.0,
  "total_reviews": 0,
  "created_at": "datetime"
}
```

### `ProductResponse`
```json
{
  "id": "uuid",
  "category_id": "uuid",
  "name": "string",
  "description": "string or null",
  "unit": "string",
  "is_available": true,
  "is_active": true,
  "tags": ["string"],
  "created_at": "datetime"
}
```

### `SpotlightResponse`
```json
{
  "id": "uuid",
  "business_id": "uuid",
  "type": "product|offer|general",
  "status": "draft|published|expired",
  "title": "string",
  "description": "string or null",
  "media_url": "string",
  "media_type": "image|video",
  "thumbnail_url": "string or null",
  "product_id": "uuid or null",
  "offer_id": "uuid or null",
  "view_count": 0,
  "like_count": 0,
  "share_count": 0,
  "comment_count": 0,
  "expires_at": "datetime or null",
  "created_at": "datetime"
}
```

### `OfferResponse`
```json
{
  "id": "uuid",
  "business_id": "uuid",
  "title": "string",
  "description": "string or null",
  "discount_type": "percentage|flat",
  "discount_value": 0,
  "min_order_value": 0,
  "coupon_code": "string or null",
  "banner_url": "string or null",
  "is_active": true,
  "starts_at": "datetime",
  "expires_at": "datetime",
  "created_at": "datetime"
}
```

### `ProductClaimResponse`
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "business_id": "uuid",
  "inventory_id": "uuid",
  "quantity": 1,
  "status": "pending|accepted|rejected|completed|cancelled",
  "note": "string or null",
  "claimed_at": "datetime",
  "updated_at": "datetime"
}
```

### `ReviewResponse`
```json
{
  "id": "uuid",
  "business_id": "uuid",
  "user_id": "uuid",
  "rating": 5,
  "comment": "string or null",
  "is_verified": false,
  "created_at": "datetime"
}
```

### `NotificationResponse`
```json
{
  "id": "uuid",
  "type": "spotlight_posted|offer_live|new_product|review",
  "title": "string",
  "body": "string or null",
  "ref_id": "uuid or null",
  "ref_type": "spotlight|offer|product or null",
  "is_read": false,
  "created_at": "datetime"
}
```

---

## Database Migrations

Run these SQL statements in order to set up all required tables.

```sql
-- Offers
CREATE TABLE offers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id    UUID NOT NULL REFERENCES businesses(id),
    title          TEXT NOT NULL,
    description    TEXT,
    discount_type  TEXT NOT NULL,
    discount_value BIGINT NOT NULL,
    min_order_value BIGINT NOT NULL DEFAULT 0,
    coupon_code    TEXT,
    banner_url     TEXT,
    is_active      BOOLEAN NOT NULL DEFAULT true,
    starts_at      TIMESTAMPTZ NOT NULL,
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON offers(business_id);

-- Business Reviews (one review per user per business)
CREATE TABLE business_reviews (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id UUID NOT NULL REFERENCES businesses(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    rating      INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment     TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (business_id, user_id)
);
CREATE INDEX ON business_reviews(business_id);

-- Notifications
CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id),
    type       TEXT NOT NULL,
    title      TEXT NOT NULL,
    body       TEXT,
    ref_id     UUID,
    ref_type   TEXT,
    is_read    BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON notifications(user_id);
CREATE INDEX ON notifications(is_read);

-- Product Claims (add if not already present from previous migration)
CREATE TABLE IF NOT EXISTS product_claims (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id),
    business_id  UUID NOT NULL REFERENCES businesses(id),
    inventory_id UUID NOT NULL REFERENCES product_inventory(id),
    quantity     INT NOT NULL DEFAULT 1,
    status       TEXT NOT NULL DEFAULT 'pending',
    note         TEXT,
    claimed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS ON product_claims(user_id);
CREATE INDEX IF NOT EXISTS ON product_claims(business_id);
CREATE INDEX IF NOT EXISTS ON product_claims(status);
```

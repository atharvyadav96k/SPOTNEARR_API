# Spotnearr Platform


###  API Reference

#### Quick Links

| [Vendor](./API_DOCS.md#vendor-service--8081) | [Products](./API_DOCS.md#products) | [Searching](./API_DOCS.md#search-service--8082) | [User Service](./API_DOCS.md#user-service--8080) |

---


Detailed endpoints are documented in the main file:

 **[Go to API_DOCS.md](./API_DOCS.md)**

---

#### Services

| Service | Port | Purpose |
|---|---|---|
| User | `8080` | Accounts, auth, reviews, claims |
| Vendor | `8081` | Products, inventory, categories |
| Search | `8082` | Public proximity search |

---


#### Core Routes

```
POST /api/v1/auth/login
GET  /api/v1/search?q=&lat=&long=&range=
POST /api/v1/inventory/{invId}/products
```

#### Response Envelope

```json
{ "message": "...", "data": <object|array|null> }
```

#### Auth Header

```
Authorization: Bearer <access_token>
```

## About Spotnearr

Spotnearr is a hyperlocal marketplace platform designed to bridge the gap between physical brick-and-mortar businesses and nearby customers. By leveraging real-time geographic coordinates, the platform allows local merchants to digitally broadcast products and promotional media to consumers within their immediate vicinity.

---

## Core Features

The system serves two distinct user roles — **Businesses** and **Customers** — interconnected through location-aware discovery layers.

### 1. Merchant Capabilities (Business)

- **Multi-Branch Management:** Register multiple physical store locations with exact GPS coordinates, operational hours, and working days.
- **Localized Inventory:** Manage independent product catalogs, standard pricing, and dynamic discounted prices per location.
- **Spotlight Marketing Feed:** Broadcast real-time media updates (images or videos) categorized as *Product*, *Offer*, or *General* announcements.
- **Offer and Coupon Engine:** Targeted promotional campaigns featuring flat discounts, percentage drops, minimum order thresholds, and coupon codes.
- **Claim Management:** View incoming product reservation requests from local buyers and manually accept or reject them.

### 2. Consumer Experience (User)

- **Proximity-Based Discovery:** Browse and search active products sold by verified businesses within a **5 km radius** of current coordinates.
- **Hyperlocal Social Feed:** Access a tailored "Spotlight" stream showing active promotional media from local shops within a **10 km radius**.
- **Product Claims & Reservations:** Instantly reserve items from a local store's inventory — stock is locked immediately to prevent overselling for in-store pickup.
- **Social Engagement:** Follow businesses, follow other users, leave 1–5 star reviews, and save or like items and spotlights.
- **Real-Time Alerts:** Instant notifications for live offers, nearby product updates, and comment interactions.

---

## Technical Architecture

- **Runtime:** Built entirely in *Go 1.22* and deployed as *Google Cloud Functions (2nd Gen)*, where each individual endpoint acts as an independent Cloud Run microservice.
- **Database:** A central *PostgreSQL* instance managing relational integrity, user/business entities, and location-based proximity calculations.
- **Storage & Media:** Integrated with *Google Cloud Storage (GCS)*. Clients request secure, short-lived Signed URLs from a dedicated internal bridge service to upload multimedia directly to GCS buckets.
- **Automated QA:** Self-contained end-to-end integration test suite deployed on GCP, triggered automatically via Cloud Scheduler — simulates 34 complete user journeys from registration to product claiming.


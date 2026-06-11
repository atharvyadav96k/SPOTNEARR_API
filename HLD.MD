
    # Product Search System — High Level Design

A scalable architecture for a product upload and search platform built around three independent services, designed to handle millions of requests.

---

![Alt Text](hld_dig.svg)

---
## Overview

The system allows vendors to upload products and users to search for them with accurate, nearest-match results. Product data flows asynchronously from the vendor service into a search index, keeping uploads fast and search reliable at scale.

---

## Architecture

```
Vendor client                          User client
     |                                      |
     └──────────────┬───────────────────────┘
                    ▼
              API gateway
          Auth · Rate limit · Route
                    |
              Load balancer
          ┌─────────┼──────────┐
          ▼         ▼          ▼
    User service  Vendor    Search
                  service   service
                    |    ──────────▶ (consume)
                    ▼
              Message queue
                (publish)
```

---

## Services

### User service
- Handles user registration and login
- Issues auth tokens used across all services
- Reads and writes to the **user database**

### Vendor service
- Accepts product uploads from vendors
- Parses and validates the upload token
- Stores products in the **product database**
- Publishes a message to the **message queue** after every successful upload

### Search service
- Consumes messages from the queue asynchronously
- Indexes incoming products into the **search index**
- Serves nearest-match search queries to users
- Uses a **cache** to serve repeated or popular queries without hitting the index

---

## Data Flow

### Product upload flow

```
1. Vendor sends product  →  API gateway  →  Vendor service
2. Vendor service saves product  →  Product database
3. Vendor service publishes event  →  Message queue
4. Search service consumes event  →  Search index updated
```

### User search flow

```
1. User sends query  →  API gateway  →  Search service
2. Search service checks Cache
   ├── HIT  →  return result immediately
   └── MISS →  query Search index  →  return result  →  store in Cache
```

---

## Components

| Component | Role |
|---|---|
| API gateway | Single entry point. Handles auth, rate limiting, and routing |
| Load balancer | Distributes traffic across service instances |
| User service | Registration, login, token issuance |
| Vendor service | Product upload, token parsing, storage |
| Search service | Product indexing and nearest-match query |
| Message queue | Async event bus between vendor and search service |
| User database | Stores user profiles |
| Product database | Stores vendor product data |
| Search index | Optimised storage for fast search queries |
| Cache | Stores hot/frequent query results |
| CDN | Serves static assets like product images |
| Monitoring | Logs, alerts, and metrics across all services |

---

## Scalability Design

### Horizontal scaling
Each service is stateless and sits behind its own load balancer. Any service can be scaled independently based on load — if search traffic spikes, only the search service scales up.

### Async product indexing
The vendor service does not write directly to the search index. It publishes to a message queue and returns a fast response to the vendor. The search service indexes in the background. This decouples the two services and prevents upload latency from growing with search index size.

### Caching
Frequently searched queries are stored in the cache layer. This reduces load on the search index and keeps response times low for popular searches.

### Rate limiting at the gateway
The API gateway enforces rate limits before requests reach any service. This protects all downstream services from traffic bursts and abuse.

---

## Key Design Decisions

**Why a message queue between vendor and search?**
Direct synchronous writes to the search index on every upload would couple the two services tightly. If the search index is slow or down, uploads would fail. The queue makes them independent — uploads always succeed and indexing catches up.

**Why a separate search index?**
The product database is optimised for writes and structured queries. Search requires a different access pattern — fuzzy matching, ranking, and nearest-neighbour lookup. A dedicated search index serves this without degrading product database performance.

**Why a cache in front of search?**
User search patterns tend to repeat. A cache with a short TTL significantly reduces the number of index queries, especially during high-traffic periods.

---

## Folder Structure (suggested)

```
/
├── user-service/
├── vendor-service/
├── search-service/
├── api-gateway/
├── docs/
│   └── hld.md
└── README.md
```

---

## Getting Started

Each service is independently deployable. Start with the API gateway, then bring up individual services as needed.

1. Configure auth settings in the API gateway
2. Start the user service and connect to the user database
3. Start the vendor service and connect to the product database and message queue
4. Start the search service and connect to the search index and message queue
5. Deploy the CDN for static assets
6. Set up monitoring across all services

---

## Contributing

Each service owns its own codebase, database, and deployment. Changes to one service should not require changes to another. Communication between services happens only through the API gateway or the message queue.
## About Spotnearr

Spotnearr is a hyperlocal marketplace platform designed to bridge the gap between physical brick-and-mortar businesses and nearby customers. By leveraging real-time geographic coordinates, the platform allows local merchants to digitally broadcast products and promotional media to consumers within their immediate vicinity.

---

## Core Features

The system serves two distinct user roles—**Businesses** and **Customers**—interconnected through location-aware discovery layers.

### 1. Merchant Capabilities (Business)
* **Multi-Branch Management:** Businesses can register multiple physical store locations with exact GPS coordinates (latitude and longitude), operational hours, and working days.
* **Localized Inventory:** Merchants can manage independent product catalogs, standard pricing, and dynamic discounted prices per location.
* **Spotlight Marketing Feed:** Businesses can broadcast real-time media updates (images or videos) categorized as *Product*, *Offer*, or *General* announcements.
* **Offer and Coupon Engine:** Built-in support for targeted promotional campaigns featuring flat discounts, percentage drops, minimum order thresholds, and coupon codes.
* **Claim Management:** Merchants can view incoming product reservation requests from local buyers and manually accept or reject them.

### 2. Consumer Experience (User)
* **Proximity-Based Discovery:** Customers can browse and search active products sold by verified businesses located within a 5 km radius of their current coordinates.
* **Hyperlocal Social Feed:** Users can access a tailored "Spotlight" stream showing active promotional media from local shops within a 10 km radius.
* **Product Claims and Reservations:** Users can instantly reserve items from a local store's inventory. The system immediately locks the stock to prevent overselling, allowing the user to pick up and complete the transaction in-store.
* **Social Engagement:** Customers can follow businesses to curate their feed, follow other users, leave structured 1-to-5 star reviews, and save or like items and spotlights.
* **Real-Time System Alerts:** Users receive instant notifications regarding live offers, nearby product updates, and comment interactions.

---

## Technical Architecture

The application is engineered as a highly modular, decoupled microservices backend optimized for serverless environments:

* **Runtime:** Built entirely in **Go 1.22** and deployed as **Google Cloud Functions (2nd Gen)**, where each individual endpoint acts as an independent Cloud Run microservice.
* **Database:** A central **PostgreSQL** instance managing relational integrity, user/business entities, and location-based proximity calculations.
* **Storage and Media:** Integrated with **Google Cloud Storage (GCS)**. To optimize bandwidth, clients request secure, short-lived *Signed URLs* from a dedicated internal bridge service to upload multimedia directly to GCS buckets.
* **Automated Quality Assurance:** Features a self-contained end-to-end integration test suite deployed on GCP. Triggered automatically via Cloud Scheduler, it simulates 34 complete user journeys—from user registration and inventory initialization to real-time product claiming and feedback loops.
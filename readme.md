# Spotnearr Platform

<table>
  <tr>
    <td width="35%" valign="top" style="border-right: 1px solid #30363d; padding-right: 15px;">
      <h3>📖 API Reference</h3>
      <p>Detailed endpoints are documented in the main file:</p>
      <p>👉 <strong><a href="./API_DOCS.md">Go to API_DOCS.md</a></strong></p>
      <hr>
      <h4>Services</h4>
      <ul>
        <li><strong>User Service</strong> — <code>:8080</code> — Customer accounts, auth, reviews, claims</li>
        <li><strong>Vendor Service</strong> — <code>:8081</code> — Business dashboard, products, inventory</li>
        <li><strong>Search Service</strong> — <code>:8082</code> — Public proximity search, no auth</li>
      </ul>
      <hr>
      <h4>Quick Links</h4>
      <ul>
        <li><a href="./API_DOCS.md#vendor-service----8081">Vendor Auth</a></li>
        <li><a href="./API_DOCS.md#products">Products</a></li>
        <li><a href="./API_DOCS.md#inventory-stores">Inventory &amp; Stores</a></li>
        <li><a href="./API_DOCS.md#search-service----8082">Proximity Search</a></li>
        <li><a href="./API_DOCS.md#user-service----8080">User Service Auth</a></li>
      </ul>
      <hr>
      <h4>Core Routes</h4>
      <p><code>POST /api/v1/auth/login</code></p>
      <p><code>GET  /api/v1/search?q=&amp;lat=&amp;long=&amp;range=</code></p>
      <p><code>POST /api/v1/inventory/{invId}/products</code></p>
      <hr>
      <h4>Response Envelope</h4>
      <pre><code>{ "message": "...", "data": &lt;object|array|null&gt; }</code></pre>
      <h4>Auth Header</h4>
      <pre><code>Authorization: Bearer &lt;access_token&gt;</code></pre>
    </td>

    <td width="65%" valign="top" style="padding-left: 15px;">
      <h2>About Spotnearr</h2>
      <p>Spotnearr is a hyperlocal marketplace platform designed to bridge the gap between physical brick-and-mortar businesses and nearby customers. By leveraging real-time geographic coordinates, the platform allows local merchants to digitally broadcast products and promotional media to consumers within their immediate vicinity.</p>

      <h2>Core Features</h2>
      <p>The system serves two distinct user roles — <strong>Businesses</strong> and <strong>Customers</strong> — interconnected through location-aware discovery layers.</p>

      <h3>1. Merchant Capabilities (Business)</h3>
      <ul>
        <li><strong>Multi-Branch Management:</strong> Register multiple physical store locations with exact GPS coordinates, operational hours, and working days.</li>
        <li><strong>Localized Inventory:</strong> Manage independent product catalogs, standard pricing, and dynamic discounted prices per location.</li>
        <li><strong>Spotlight Marketing Feed:</strong> Broadcast real-time media updates (images or videos) categorized as <em>Product</em>, <em>Offer</em>, or <em>General</em> announcements.</li>
        <li><strong>Offer and Coupon Engine:</strong> Targeted promotional campaigns featuring flat discounts, percentage drops, minimum order thresholds, and coupon codes.</li>
        <li><strong>Claim Management:</strong> View incoming product reservation requests from local buyers and manually accept or reject them.</li>
      </ul>

      <h3>2. Consumer Experience (User)</h3>
      <ul>
        <li><strong>Proximity-Based Discovery:</strong> Browse and search active products sold by verified businesses within a <strong>5 km radius</strong> of current coordinates.</li>
        <li><strong>Hyperlocal Social Feed:</strong> Access a tailored "Spotlight" stream showing active promotional media from local shops within a <strong>10 km radius</strong>.</li>
        <li><strong>Product Claims &amp; Reservations:</strong> Instantly reserve items from a local store's inventory — stock is locked immediately to prevent overselling for in-store pickup.</li>
        <li><strong>Social Engagement:</strong> Follow businesses, follow other users, leave 1–5 star reviews, and save or like items and spotlights.</li>
        <li><strong>Real-Time Alerts:</strong> Instant notifications for live offers, nearby product updates, and comment interactions.</li>
      </ul>

      <h2>Technical Architecture</h2>
      <ul>
        <li><strong>Runtime:</strong> Built entirely in <em>Go 1.22</em> and deployed as <em>Google Cloud Functions (2nd Gen)</em>, where each individual endpoint acts as an independent Cloud Run microservice.</li>
        <li><strong>Database:</strong> A central <em>PostgreSQL</em> instance managing relational integrity, user/business entities, and location-based proximity calculations.</li>
        <li><strong>Storage &amp; Media:</strong> Integrated with <em>Google Cloud Storage (GCS)</em>. Clients request secure, short-lived Signed URLs from a dedicated internal bridge service to upload multimedia directly to GCS buckets.</li>
        <li><strong>Automated QA:</strong> Self-contained end-to-end integration test suite deployed on GCP, triggered automatically via Cloud Scheduler — simulates 34 complete user journeys from registration to product claiming.</li>
      </ul>
    </td>
  </tr>
</table>

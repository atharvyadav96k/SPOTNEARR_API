package applayer

import (
	"net/http"
	"time"

	pkgmid "github.com/atharvyadav96k/spotnearr/pkg/middleware"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/config"
	"github.com/gorilla/mux"
)

func (a *application) NewMux() *mux.Router {
	router := mux.NewRouter()
	router.Use(pkgmid.CORS)
	router.Use(pkgmid.PrometheusMetrics("vendor-svc"))
	router.Handle("/metrics", pkgmid.MetricsHandler())

	auth := pkgmid.Auth(config.C.JWTSecret)
	bizOnly := pkgmid.BusinessOnly(config.C.JWTSecret)
	rl := a.cache.GetRateLimit()

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	a.healthRouter(apiV1)
	a.authRouter(apiV1)
	a.businessRouter(apiV1, auth, bizOnly, rl)
	a.inventoryRouter(apiV1, auth, bizOnly, rl)
	a.productRouter(apiV1, auth, bizOnly, rl)
	a.categoryRouter(apiV1, auth, bizOnly, rl)
	a.offerRouter(apiV1, auth, bizOnly, rl)
	a.spotlightRouter(apiV1, auth, bizOnly, rl)
	a.claimRouter(apiV1, auth, bizOnly, rl)
	a.feedRouter(apiV1, rl)

	// Internal routes — no auth middleware, expected to be network-isolated.
	internal := router.PathPrefix("/internal").Subrouter()
	a.internalRouter(internal)

	return router
}

func (a *application) authRouter(router *mux.Router) {
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", a.authHandler.Register).Methods(http.MethodPost)
	auth.HandleFunc("/login", a.authHandler.Login).Methods(http.MethodPost)
	auth.HandleFunc("/refresh", a.authHandler.Refresh).Methods(http.MethodPost)
}

func (a *application) healthRouter(router *mux.Router) {
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)
}

func (a *application) businessRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 300, time.Minute)
	strictRL := pkgmid.RateLimit(rl, 500, time.Minute)

	biz := router.PathPrefix("/businesses").Subrouter()
	biz.Use(auth)

	biz.Handle("/register",
		strictRL(http.HandlerFunc(a.bizHandler.BusinessRegister)),
	).Methods(http.MethodPost)

	biz.Handle("/{bizId}/profile",
		normalRL(http.HandlerFunc(a.bizHandler.BusinessProfile)),
	).Methods(http.MethodGet)

	bizOnlySub := biz.PathPrefix("").Subrouter()
	bizOnlySub.Use(bizOnly)

	bizOnlySub.Handle("/",
		normalRL(http.HandlerFunc(a.bizHandler.BusinessUpdate)),
	).Methods(http.MethodPatch)
}

func (a *application) inventoryRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 100000, time.Minute)
	strictRL := pkgmid.RateLimit(rl, 100000, time.Minute)

	inv := router.PathPrefix("/inventory").Subrouter()
	inv.Use(auth)
	inv.Use(bizOnly)

	inv.Handle("/",
		normalRL(http.HandlerFunc(a.invHandler.InventoryList)),
	).Methods(http.MethodGet)

	inv.Handle("/",
		strictRL(http.HandlerFunc(a.invHandler.InventoryCreate)),
	).Methods(http.MethodPost)

	inv.Handle("/{invId}",
		normalRL(http.HandlerFunc(a.invHandler.InventoryUpdate)),
	).Methods(http.MethodPatch)

	inv.Handle("/{invId}/products",
		normalRL(http.HandlerFunc(a.invHandler.InventoryGetProducts)),
	).Methods(http.MethodGet)

	inv.Handle("/{invId}/products",
		normalRL(http.HandlerFunc(a.invHandler.InventoryAddProduct)),
	).Methods(http.MethodPost)

	inv.Handle("/{invId}/products/{invProductId}",
		normalRL(http.HandlerFunc(a.invHandler.InventoryUpdateProduct)),
	).Methods(http.MethodPatch)

	inv.Handle("/{invId}/products/{invProductId}",
		strictRL(http.HandlerFunc(a.invHandler.InventoryRemoveProduct)),
	).Methods(http.MethodDelete)
}

func (a *application) productRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 3000, time.Minute)
	strictRL := pkgmid.RateLimit(rl, 5000, time.Minute)
	relaxedRL := pkgmid.RateLimit(rl, 6000, time.Minute)

	products := router.PathPrefix("/products").Subrouter()
	products.Use(auth)

	// Auth only — any registered user (or vendor) can fetch product detail.
	products.Handle("/{invProductId}/detail",
		relaxedRL(http.HandlerFunc(a.productHandler.ProductDetail)),
	).Methods(http.MethodGet)

	// Vendor-only routes — require business token on top of auth.
	bizProducts := products.PathPrefix("").Subrouter()
	bizProducts.Use(bizOnly)

	bizProducts.Handle("/",
		strictRL(http.HandlerFunc(a.productHandler.ProductAdd)),
	).Methods(http.MethodPost)

	bizProducts.Handle("/",
		relaxedRL(http.HandlerFunc(a.productHandler.ProductList)),
	).Methods(http.MethodGet)

	bizProducts.Handle("/{productId}",
		relaxedRL(http.HandlerFunc(a.productHandler.ProductGet)),
	).Methods(http.MethodGet)

	bizProducts.Handle("/{productId}",
		normalRL(http.HandlerFunc(a.productHandler.ProductUpdate)),
	).Methods(http.MethodPatch)

	bizProducts.Handle("/{productId}",
		strictRL(http.HandlerFunc(a.productHandler.ProductDelete)),
	).Methods(http.MethodDelete)
}

func (a *application) categoryRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 3000, time.Minute)
	strictRL := pkgmid.RateLimit(rl, 5000, time.Minute)

	cat := router.PathPrefix("/categories").Subrouter()
	cat.Use(auth)

	cat.Handle("/",
		normalRL(http.HandlerFunc(a.categoryHandler.CategoryList)),
	).Methods(http.MethodGet)

	bizOnlySub := cat.PathPrefix("").Subrouter()
	bizOnlySub.Use(bizOnly)

	bizOnlySub.Handle("/",
		strictRL(http.HandlerFunc(a.categoryHandler.CategoryAdd)),
	).Methods(http.MethodPost)
}

func (a *application) offerRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 3000, time.Minute)
	strictRL := pkgmid.RateLimit(rl, 5000, time.Minute)

	offers := router.PathPrefix("/offers").Subrouter()
	offers.Use(auth)
	offers.Use(bizOnly)

	offers.Handle("/",
		normalRL(http.HandlerFunc(a.offerHandler.OfferList)),
	).Methods(http.MethodGet)

	offers.Handle("/",
		strictRL(http.HandlerFunc(a.offerHandler.OfferCreate)),
	).Methods(http.MethodPost)

	offers.Handle("/{offerId}",
		normalRL(http.HandlerFunc(a.offerHandler.OfferGet)),
	).Methods(http.MethodGet)

	offers.Handle("/{offerId}",
		normalRL(http.HandlerFunc(a.offerHandler.OfferUpdate)),
	).Methods(http.MethodPatch)

	offers.Handle("/{offerId}",
		strictRL(http.HandlerFunc(a.offerHandler.OfferDelete)),
	).Methods(http.MethodDelete)
}

func (a *application) spotlightRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 3000, time.Minute)
	strictRL := pkgmid.RateLimit(rl, 5000, time.Minute)

	spotlights := router.PathPrefix("/spotlights").Subrouter()
	spotlights.Use(auth)
	spotlights.Use(bizOnly)

	spotlights.Handle("/",
		normalRL(http.HandlerFunc(a.spotlightHandler.SpotlightList)),
	).Methods(http.MethodGet)

	spotlights.Handle("/",
		strictRL(http.HandlerFunc(a.spotlightHandler.SpotlightPost)),
	).Methods(http.MethodPost)

	spotlights.Handle("/{spotlightId}",
		normalRL(http.HandlerFunc(a.spotlightHandler.SpotlightGet)),
	).Methods(http.MethodGet)

	spotlights.Handle("/{spotlightId}",
		strictRL(http.HandlerFunc(a.spotlightHandler.SpotlightDelete)),
	).Methods(http.MethodDelete)
}

func (a *application) claimRouter(router *mux.Router, auth, bizOnly mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRL := pkgmid.RateLimit(rl, 3000, time.Minute)

	claims := router.PathPrefix("/claims").Subrouter()
	claims.Use(auth)
	claims.Use(bizOnly)

	claims.Handle("/",
		normalRL(http.HandlerFunc(a.claimHandler.ListClaims)),
	).Methods(http.MethodGet)

	claims.Handle("/{claimId}/accept",
		normalRL(http.HandlerFunc(a.claimHandler.AcceptClaim)),
	).Methods(http.MethodPatch)

	claims.Handle("/{claimId}/reject",
		normalRL(http.HandlerFunc(a.claimHandler.RejectClaim)),
	).Methods(http.MethodPatch)
}

func (a *application) feedRouter(router *mux.Router, rl pkgmid.RateLimiter) {
	relaxedRL := pkgmid.RateLimit(rl, 60000, time.Minute)

	feed := router.PathPrefix("/feed").Subrouter()
	feed.Handle("/spotlights",
		relaxedRL(http.HandlerFunc(a.spotlightHandler.SpotlightFeed)),
	).Methods(http.MethodGet)
}

func (a *application) internalRouter(router *mux.Router) {
	router.HandleFunc("/inventory-products/{id}", a.internalHandler.GetInventoryProduct).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}/access", a.internalHandler.GetUserAccess).Methods(http.MethodGet)
	router.HandleFunc("/businesses/{bizId}/follow", a.internalHandler.FollowBusiness).Methods(http.MethodPost)
	router.HandleFunc("/businesses/{bizId}/unfollow", a.internalHandler.UnfollowBusiness).Methods(http.MethodPost)
}

package applayer

import (
	"net/http"
	"time"

	middleware "github.com/atharvyadav96k/SPOTNEARR_API/middlewares"
	"github.com/gorilla/mux"
)

func (a *application) NewMux() *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.CORS)
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	a.healthRouter(apiV1)
	a.authRouter(apiV1)
	a.userRouter(apiV1)
	a.businessRouter(apiV1)
	a.inventoryRouter(apiV1)
	a.productRouter(apiV1)
	a.claimRouter(apiV1)
	a.offerRouter(apiV1)
	a.reviewRouter(apiV1)
	a.searchRouter(apiV1)
	a.categoryRouter(apiV1)
	// Internal routes — network-isolated, no auth middleware.
	internal := router.PathPrefix("/internal").Subrouter()
	a.internalRouter(internal)
	return router
}

func (a *application) internalRouter(router *mux.Router) {
	router.HandleFunc("/users/{userId}/invalidate-refresh",
		a.internalHandler.InvalidateRefresh).Methods(http.MethodPost)
}

func (a *application) healthRouter(router *mux.Router) {
	router.HandleFunc("/health", a.healthHandler.HealthOK).Methods(http.MethodGet)
}

func (a *application) searchRouter(router *mux.Router) {
	router.HandleFunc("/search", a.searchHandler.Search).Methods(http.MethodGet)
}

func (a *application) categoryRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 3000, time.Minute)
	strictRateLimit := middleware.RateLimit(a.GetCache(), 5000, time.Minute)

	cat := router.PathPrefix("/categories").Subrouter()
	cat.Use(middleware.Auth)
	cat.Handle("/", normalRateLimit(http.HandlerFunc(a.categoryHandler.CategoryList))).Methods(http.MethodGet)

	bizOnly := cat.PathPrefix("").Subrouter()
	bizOnly.Use(middleware.BusinessOnly)
	bizOnly.Handle("/", strictRateLimit(http.HandlerFunc(a.categoryHandler.CategoryAdd))).Methods(http.MethodPost)
}

func (a *application) authRouter(router *mux.Router) {
	authBase := router.PathPrefix("/auth").Subrouter()

	authBase.HandleFunc("/refresh", a.authHandler.RefreshToken).Methods(http.MethodPost)

	captchaRoutes := authBase.PathPrefix("").Subrouter()
	captchaRoutes.Use(middleware.CaptchaValidation)
	captchaRoutes.HandleFunc("/users/register", a.authHandler.Register).Methods(http.MethodPost)
	captchaRoutes.HandleFunc("/users/login", a.authHandler.Login).Methods(http.MethodPost)
	captchaRoutes.HandleFunc("/reset-request", a.authHandler.Session).Methods(http.MethodPost)

	authBase.Handle("/reset-password",
		middleware.SessionValidation(
			middleware.CaptchaValidation(
				http.HandlerFunc(a.authHandler.ResetPassword),
			),
		),
	).Methods(http.MethodPost)

	protectedAuth := authBase.PathPrefix("").Subrouter()
	protectedAuth.Use(middleware.Auth)
	protectedAuth.HandleFunc("/", a.authHandler.Auth).Methods(http.MethodGet)
	protectedAuth.HandleFunc("/logout-all-devices", a.authHandler.LogoutFromAllDevices).Methods(http.MethodPost)
}
func (a *application) userRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 1000, time.Minute)
	strictRateLimit := middleware.RateLimit(a.GetCache(), 5000, time.Minute)

	protectedAuth := router.PathPrefix("/users").Subrouter()
	protectedAuth.Use(middleware.Auth)

	protectedAuth.Handle("/{userId}/profile",
		normalRateLimit(
			http.HandlerFunc(a.userHandler.Profile),
		),
	).Methods(http.MethodGet)

	protectedAuth.Handle("/{userId}/ban",
		strictRateLimit(
			http.HandlerFunc(a.userHandler.BanUser),
		),
	)
}

func (a *application) businessRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 300, time.Minute)
	strictRateLimit := middleware.RateLimit(a.GetCache(), 500, time.Minute)

	biz := router.PathPrefix("/businesses").Subrouter()
	protectedAuth := biz.PathPrefix("").Subrouter()
	protectedAuth.Use(middleware.Auth)

	protectedAuth.Handle("/register",
		strictRateLimit(
			http.HandlerFunc(a.businessHandler.BusinessRegister),
		),
	).Methods(http.MethodPost)

	protectedAuth.Handle("/{bizId}/profile",
		normalRateLimit(
			http.HandlerFunc(a.businessHandler.BusinessProfile),
		),
	).Methods(http.MethodGet)

	bizOnly := protectedAuth.PathPrefix("").Subrouter()
	bizOnly.Use(middleware.BusinessOnly)

	bizOnly.Handle("/inventories",
		normalRateLimit(
			http.HandlerFunc(a.businessHandler.BusinessInventories),
		),
	).Methods(http.MethodGet)

	bizOnly.Handle("/",
		normalRateLimit(
			http.HandlerFunc(a.businessHandler.BusinessUpdate),
		),
	).Methods(http.MethodPatch)

	bizOnly.Handle("/",
		strictRateLimit(
			http.HandlerFunc(a.businessHandler.BusinessDelete),
		),
	).Methods(http.MethodDelete)
}

func (a *application) inventoryRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 100000, time.Minute)
	strictRateLimit := middleware.RateLimit(a.GetCache(), 100000, time.Minute)

	protectedAuth := router.PathPrefix("/inventory").Subrouter()
	protectedAuth.Use(middleware.Auth)
	protectedAuth.Use(middleware.BusinessOnly)

	protectedAuth.Handle("/",
		strictRateLimit(
			http.HandlerFunc(a.inventoryHandler.InventoryCreate),
		),
	).Methods(http.MethodPost)

	protectedAuth.Handle("/{invId}",
		normalRateLimit(
			http.HandlerFunc(a.inventoryHandler.InventoryUpdate),
		),
	).Methods(http.MethodPatch)

	protectedAuth.Handle("/{invId}",
		strictRateLimit(
			http.HandlerFunc(a.inventoryHandler.InventoryDelete),
		),
	).Methods(http.MethodDelete)

	protectedAuth.Handle("/{invId}/products",
		normalRateLimit(
			http.HandlerFunc(a.inventoryHandler.InventoryGetProducts),
		),
	).Methods(http.MethodGet)

	protectedAuth.Handle("/{invId}/products",
		normalRateLimit(
			http.HandlerFunc(a.inventoryHandler.InventoryAddProduct),
		),
	).Methods(http.MethodPost)

	protectedAuth.Handle("/{invId}/products",
		strictRateLimit(
			http.HandlerFunc(a.inventoryHandler.InventoryRemoveProduct),
		),
	).Methods(http.MethodDelete)
}

func (a *application) productRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 3000, time.Minute)
	strictRateLimit := middleware.RateLimit(a.GetCache(), 5000, time.Minute)
	relaxedRateLimit := middleware.RateLimit(a.GetCache(), 60000, time.Minute)

	protectedAuth := router.PathPrefix("/products").Subrouter()
	protectedAuth.Use(middleware.Auth)
	protectedAuth.Use(middleware.BusinessOnly)

	protectedAuth.Handle("/",
		strictRateLimit(
			http.HandlerFunc(a.productHandler.ProductAdd),
		),
	).Methods(http.MethodPost)

	protectedAuth.Handle("/{productId}",
		relaxedRateLimit(
			http.HandlerFunc(a.productHandler.ProductGet),
		),
	).Methods(http.MethodGet)

	protectedAuth.Handle("/{productId}",
		normalRateLimit(
			http.HandlerFunc(a.productHandler.ProductUpdate),
		),
	).Methods(http.MethodPatch)

	protectedAuth.Handle("/{productId}",
		strictRateLimit(
			http.HandlerFunc(a.productHandler.ProductDelete),
		),
	).Methods(http.MethodDelete)

	protectedAuth.Handle("/nearby",
		relaxedRateLimit(
			http.HandlerFunc(a.productHandler.ProductNearBy),
		),
	).Methods(http.MethodGet)
}

func (a *application) claimRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 10000, time.Minute)

	protectedAuth := router.PathPrefix("/claims").Subrouter()
	protectedAuth.Use(middleware.Auth)

	protectedAuth.Handle("",
		normalRateLimit(
			http.HandlerFunc(a.claimHandler.GetUserClaims),
		),
	).Methods(http.MethodGet)

	protectedAuth.Handle("/{invProductId}",
		normalRateLimit(
			http.HandlerFunc(a.claimHandler.ClaimProduct),
		),
	).Methods(http.MethodPost)

	protectedAuth.Handle("/{claimId}",
		normalRateLimit(
			http.HandlerFunc(a.claimHandler.ClaimRemove),
		),
	).Methods(http.MethodDelete)
}

func (a *application) offerRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 30000, time.Minute)
	strictRateLimit := middleware.RateLimit(a.GetCache(), 5000, time.Minute)
	relaxedRateLimit := middleware.RateLimit(a.GetCache(), 60000, time.Minute)

	protectedAuth := router.PathPrefix("/offers").Subrouter()
	protectedAuth.Use(middleware.Auth)

	protectedAuth.Handle("/nearby",
		relaxedRateLimit(
			http.HandlerFunc(a.offerHandler.NearByOffers),
		),
	).Methods(http.MethodGet)

	bizOnly := protectedAuth.PathPrefix("/").Subrouter()
	bizOnly.Use(middleware.BusinessOnly)

	bizOnly.Handle("",
		strictRateLimit(
			http.HandlerFunc(a.offerHandler.OfferAdd),
		),
	).Methods(http.MethodPost)

	bizOnly.Handle("/{offerID}",
		normalRateLimit(
			http.HandlerFunc(a.offerHandler.OfferUpdate),
		),
	).Methods(http.MethodPut)

	bizOnly.Handle("/{offerID}",
		strictRateLimit(
			http.HandlerFunc(a.offerHandler.OfferDelete),
		),
	).Methods(http.MethodDelete)
}

func (a *application) reviewRouter(router *mux.Router) {
	normalRateLimit := middleware.RateLimit(a.GetCache(), 1000, time.Minute)
	relaxedRateLimit := middleware.RateLimit(a.GetCache(), 6000, time.Minute)

	r := router.PathPrefix("/review").Subrouter()
	r.Use(middleware.Auth)

	// Offer reviews
	r.Handle("/offers/{offerId}",
		relaxedRateLimit(http.HandlerFunc(a.reviewHandler.ReviewGetByOffer)),
	).Methods(http.MethodGet)
	r.Handle("/offers",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewOffer)),
	).Methods(http.MethodPost)
	r.Handle("/offers",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewOfferUpdate)),
	).Methods(http.MethodPut)
	r.Handle("/offers",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewOfferDelete)),
	).Methods(http.MethodDelete)

	// Spotlight reviews
	r.Handle("/spotlights/{spotlightId}",
		relaxedRateLimit(http.HandlerFunc(a.reviewHandler.ReviewGetBySpotlight)),
	).Methods(http.MethodGet)
	r.Handle("/spotlights",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewSpotlight)),
	).Methods(http.MethodPost)
	r.Handle("/spotlights",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewSpotlightUpdate)),
	).Methods(http.MethodPut)
	r.Handle("/spotlights",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewSpotlightDelete)),
	).Methods(http.MethodDelete)

	// Business reviews
	r.Handle("/businesses",
		relaxedRateLimit(http.HandlerFunc(a.reviewHandler.ReviewGetByBusiness)),
	).Methods(http.MethodGet)
	r.Handle("/businesses",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewBusiness)),
	).Methods(http.MethodPost)
	r.Handle("/businesses",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewBusinessUpdate)),
	).Methods(http.MethodPut)
	r.Handle("/businesses",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewBusinessDelete)),
	).Methods(http.MethodDelete)

	// Product reviews
	r.Handle("/products",
		relaxedRateLimit(http.HandlerFunc(a.reviewHandler.ReviewGetByProduct)),
	).Methods(http.MethodGet)
	r.Handle("/products",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewProduct)),
	).Methods(http.MethodPost)
	r.Handle("/products",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewProductUpdate)),
	).Methods(http.MethodPut)
	r.Handle("/products",
		normalRateLimit(http.HandlerFunc(a.reviewHandler.ReviewProductDelete)),
	).Methods(http.MethodDelete)
}

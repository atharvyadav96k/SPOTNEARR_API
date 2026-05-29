package applayer

import (
	"net/http"

	middleware "github.com/atharvyadav96k/SPOTNEARR_API/middlewares"
	"github.com/gorilla/mux"
)

func (a *application) NewMux() *mux.Router {
	router := mux.NewRouter()
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
	return router
}

func (a *application) healthRouter(router *mux.Router) {
	router.HandleFunc("/health", a.healthHandler.HealthOK).Methods(http.MethodGet)
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

	protectedAuth := router.PathPrefix("").Subrouter()
	protectedAuth.Use(middleware.Auth)
}

func (a *application) userRouter(router *mux.Router) {
	protectedAuth := router.PathPrefix("/users").Subrouter()
	protectedAuth.Use(middleware.Auth)
	protectedAuth.HandleFunc("/{userId}/profile", a.userHandler.Profile).Methods(http.MethodGet) //✅
	// protectedAuth.HandleFunc("/password", a.userHandler.UpdatePassword).Methods(http.MethodPost)
	protectedAuth.HandleFunc("/{userId}/ban", a.userHandler.BanUser)
}

func (a *application) businessRouter(router *mux.Router) {
	biz := router.PathPrefix("/businesses").Subrouter()
	protectedAuth := biz.PathPrefix("").Subrouter()
	protectedAuth.Use(middleware.Auth)
	protectedAuth.HandleFunc("/register", a.businessHandler.BusinessRegister).Methods(http.MethodPost)      // ✅
	protectedAuth.HandleFunc("/{bizId}/profile", a.businessHandler.BusinessProfile).Methods(http.MethodGet) // ✅

	bizOnly := protectedAuth.PathPrefix("").Subrouter()
	bizOnly.Use(middleware.BusinessOnly)
	bizOnly.HandleFunc("/inventories", a.businessHandler.BusinessInventories).Methods(http.MethodGet) // ✅
	bizOnly.HandleFunc("/", a.businessHandler.BusinessUpdate).Methods(http.MethodPatch)               // ✅
	bizOnly.HandleFunc("/", a.businessHandler.BusinessDelete).Methods(http.MethodDelete)
}

func (a *application) inventoryRouter(router *mux.Router) {
	protectedAuth := router.PathPrefix("/inventory").Subrouter()
	protectedAuth.Use(middleware.Auth)
	protectedAuth.Use(middleware.BusinessOnly)
	protectedAuth.HandleFunc("/", a.inventoryHandler.InventoryCreate).Methods(http.MethodPost)                     // ✅
	protectedAuth.HandleFunc("/{invId}", a.inventoryHandler.InventoryUpdate).Methods(http.MethodPatch)             // ✅
	protectedAuth.HandleFunc("/{invId}", a.inventoryHandler.InventoryDelete).Methods(http.MethodDelete)            //
	protectedAuth.HandleFunc("/{invId}/products", a.inventoryHandler.InventoryGetProducts).Methods(http.MethodGet) // ✅
	protectedAuth.HandleFunc("/{invId}/products", a.inventoryHandler.InventoryAddProduct).Methods(http.MethodPost) // ✅
	protectedAuth.HandleFunc("/{invId}/products", a.inventoryHandler.InventoryRemoveProduct).Methods(http.MethodDelete)
}

func (a *application) productRouter(router *mux.Router) {
	router.PathPrefix("/products")
	router.HandleFunc("/${productId}", a.productHandler.ProductGet).Methods(http.MethodGet)
	router.HandleFunc("/${productId}", a.productHandler.ProductUpdate).Methods(http.MethodPut)
	router.HandleFunc("/${productId}", a.productHandler.ProductDelete).Methods(http.MethodDelete)
	router.HandleFunc("/nearby", a.productHandler.ProductNearBy).Methods(http.MethodGet)
}

func (a *application) claimRouter(router *mux.Router) {
	router.PathPrefix("/claims")
	router.HandleFunc("/${productId}", a.claimHandler.ClaimProduct).Methods(http.MethodPost)
	router.HandleFunc("/${productId}", a.claimHandler.ClaimRemove).Methods(http.MethodDelete)
}

func (a *application) offerRouter(router *mux.Router) {
	router.PathPrefix("/offers")
	router.HandleFunc("/nearby", a.offerHandler.NearByOffers).Methods(http.MethodGet)
	router.HandleFunc("/", a.offerHandler.OfferAdd).Methods(http.MethodPost)
	router.HandleFunc("/", a.offerHandler.OfferUpdate).Methods(http.MethodPut)
	router.HandleFunc("/", a.offerHandler.OfferDelete).Methods(http.MethodDelete)
}

func (a *application) reviewRouter(router *mux.Router) {
	router.PathPrefix("/review")
	router.HandleFunc("/offers/${offerId}", a.reviewHandler.ReviewGetByOffer).Methods(http.MethodGet)
	router.HandleFunc("/offers", a.reviewHandler.ReviewOffer).Methods(http.MethodPost)
	router.HandleFunc("/offers", a.reviewHandler.ReviewOfferUpdate).Methods(http.MethodPut)
	router.HandleFunc("/offers", a.reviewHandler.ReviewOfferDelete).Methods(http.MethodDelete)

	router.HandleFunc("/spotlights/${spotlightId}", a.reviewHandler.ReviewGetBySpotlight).Methods(http.MethodGet)
	router.HandleFunc("/spotlights", a.reviewHandler.ReviewSpotlight).Methods(http.MethodPost)
	router.HandleFunc("/spotlights", a.reviewHandler.ReviewSpotlightUpdate).Methods(http.MethodPut)
	router.HandleFunc("/spotlights", a.reviewHandler.ReviewSpotlightDelete).Methods(http.MethodDelete)

	router.HandleFunc("/businesses", a.reviewHandler.ReviewGetByBusiness).Methods(http.MethodGet)
	router.HandleFunc("/businesses", a.reviewHandler.ReviewBusiness).Methods(http.MethodPost)
	router.HandleFunc("/businesses", a.reviewHandler.ReviewBusinessUpdate).Methods(http.MethodPut)
	router.HandleFunc("/businesses", a.reviewHandler.ReviewBusinessDelete).Methods(http.MethodDelete)

	router.HandleFunc("/products", a.reviewHandler.ReviewGetByProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", a.reviewHandler.ReviewProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", a.reviewHandler.ReviewProductUpdate).Methods(http.MethodPut)
	router.HandleFunc("/products", a.reviewHandler.ReviewProductDelete).Methods(http.MethodDelete)
}

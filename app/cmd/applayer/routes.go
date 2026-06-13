package applayer

import (
	"net/http"
	"time"

	pkgmid "github.com/atharvyadav96k/spotnearr/pkg/middleware"
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/gorilla/mux"
)

func (a *application) NewMux() *mux.Router {
	router := mux.NewRouter()
	router.Use(pkgmid.CORS)

	auth := pkgmid.Auth(config.C.JWTSecret)
	captcha := pkgmid.CaptchaValidation(config.C.CaptchaURL, config.C.CaptchaSecretKey)
	session := pkgmid.SessionValidation(config.C.JWTSecret)
	rl := a.GetCache().GetRateLimit()

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	a.healthRouter(apiV1)
	a.authRouter(apiV1, auth, captcha, session, rl)
	a.userRouter(apiV1, auth, rl)
	a.claimRouter(apiV1, auth, rl)
	a.reviewRouter(apiV1, auth, rl)
	a.socialRouter(apiV1, auth, rl)

	// Internal routes — network-isolated, no auth middleware.
	internal := router.PathPrefix("/internal").Subrouter()
	a.internalRouter(internal)

	return router
}

func (a *application) internalRouter(router *mux.Router) {
	router.HandleFunc("/users/{userId}/invalidate-refresh",
		a.internalHandler.InvalidateRefresh).Methods(http.MethodPost)

	// Claim management — called by the vendor service.
	router.HandleFunc("/claims",
		a.internalHandler.GetClaimsByProductIDs).Methods(http.MethodGet)
	router.HandleFunc("/claims/{claimId}",
		a.internalHandler.GetClaimByID).Methods(http.MethodGet)
	router.HandleFunc("/claims/{claimId}/status",
		a.internalHandler.UpdateClaimStatus).Methods(http.MethodPatch)
}

func (a *application) healthRouter(router *mux.Router) {
	router.HandleFunc("/health", a.healthHandler.HealthOK).Methods(http.MethodGet)
}

func (a *application) authRouter(router *mux.Router, auth, captcha, session mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	authBase := router.PathPrefix("/auth").Subrouter()

	authBase.HandleFunc("/refresh", a.authHandler.RefreshToken).Methods(http.MethodPost)

	captchaRoutes := authBase.PathPrefix("").Subrouter()
	captchaRoutes.Use(captcha)
	captchaRoutes.HandleFunc("/users/register", a.authHandler.Register).Methods(http.MethodPost)
	captchaRoutes.HandleFunc("/users/login", a.authHandler.Login).Methods(http.MethodPost)
	captchaRoutes.HandleFunc("/reset-request", a.authHandler.Session).Methods(http.MethodPost)

	authBase.Handle("/reset-password",
		session(captcha(http.HandlerFunc(a.authHandler.ResetPassword))),
	).Methods(http.MethodPost)

	protectedAuth := authBase.PathPrefix("").Subrouter()
	protectedAuth.Use(auth)
	protectedAuth.HandleFunc("/", a.authHandler.Auth).Methods(http.MethodGet)
	protectedAuth.HandleFunc("/logout-all-devices", a.authHandler.LogoutFromAllDevices).Methods(http.MethodPost)
}

func (a *application) userRouter(router *mux.Router, auth mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRateLimit := pkgmid.RateLimit(rl, 1000, time.Minute)
	strictRateLimit := pkgmid.RateLimit(rl, 5000, time.Minute)

	protectedAuth := router.PathPrefix("/users").Subrouter()
	protectedAuth.Use(auth)

	protectedAuth.Handle("/{userId}/profile",
		normalRateLimit(http.HandlerFunc(a.userHandler.Profile)),
	).Methods(http.MethodGet)

	protectedAuth.Handle("/{userId}/ban",
		strictRateLimit(http.HandlerFunc(a.userHandler.BanUser)),
	)
}

func (a *application) claimRouter(router *mux.Router, auth mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRateLimit := pkgmid.RateLimit(rl, 10000, time.Minute)

	protectedAuth := router.PathPrefix("/claims").Subrouter()
	protectedAuth.Use(auth)

	protectedAuth.Handle("",
		normalRateLimit(http.HandlerFunc(a.claimHandler.GetUserClaims)),
	).Methods(http.MethodGet)

	protectedAuth.Handle("/{invProductId}",
		normalRateLimit(http.HandlerFunc(a.claimHandler.ClaimProduct)),
	).Methods(http.MethodPost)

	protectedAuth.Handle("/{claimId}",
		normalRateLimit(http.HandlerFunc(a.claimHandler.ClaimRemove)),
	).Methods(http.MethodDelete)
}

func (a *application) socialRouter(router *mux.Router, auth mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRateLimit := pkgmid.RateLimit(rl, 3000, time.Minute)

	s := router.PathPrefix("/social").Subrouter()
	s.Use(auth)

	// Follow / unfollow business
	s.Handle("/businesses/{bizId}/follow",
		normalRateLimit(http.HandlerFunc(a.socialHandler.FollowBusiness)),
	).Methods(http.MethodPost)
	s.Handle("/businesses/{bizId}/follow",
		normalRateLimit(http.HandlerFunc(a.socialHandler.UnfollowBusiness)),
	).Methods(http.MethodDelete)

	// Like / save product
	s.Handle("/products/{invProductId}/like",
		normalRateLimit(http.HandlerFunc(a.socialHandler.LikeProduct)),
	).Methods(http.MethodPost)
	s.Handle("/products/{invProductId}/like",
		normalRateLimit(http.HandlerFunc(a.socialHandler.UnlikeProduct)),
	).Methods(http.MethodDelete)
	s.Handle("/products/{invProductId}/save",
		normalRateLimit(http.HandlerFunc(a.socialHandler.SaveProduct)),
	).Methods(http.MethodPost)
	s.Handle("/products/{invProductId}/save",
		normalRateLimit(http.HandlerFunc(a.socialHandler.UnsaveProduct)),
	).Methods(http.MethodDelete)

	// Like / save spotlight
	s.Handle("/spotlights/{spotlightId}/like",
		normalRateLimit(http.HandlerFunc(a.socialHandler.LikeSpotlight)),
	).Methods(http.MethodPost)
	s.Handle("/spotlights/{spotlightId}/like",
		normalRateLimit(http.HandlerFunc(a.socialHandler.UnlikeSpotlight)),
	).Methods(http.MethodDelete)
	s.Handle("/spotlights/{spotlightId}/save",
		normalRateLimit(http.HandlerFunc(a.socialHandler.SaveSpotlight)),
	).Methods(http.MethodPost)
	s.Handle("/spotlights/{spotlightId}/save",
		normalRateLimit(http.HandlerFunc(a.socialHandler.UnsaveSpotlight)),
	).Methods(http.MethodDelete)
}

func (a *application) reviewRouter(router *mux.Router, auth mux.MiddlewareFunc, rl pkgmid.RateLimiter) {
	normalRateLimit := pkgmid.RateLimit(rl, 1000, time.Minute)
	relaxedRateLimit := pkgmid.RateLimit(rl, 6000, time.Minute)

	r := router.PathPrefix("/review").Subrouter()
	r.Use(auth)

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

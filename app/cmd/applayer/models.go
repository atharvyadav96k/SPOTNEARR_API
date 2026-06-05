package applayer

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/app"
	"github.com/atharvyadav96k/SPOTNEARR_API/handlers"
)

type application struct {
	app.App
	healthHandler    *handlers.Health
	authHandler      *handlers.AuthHandler
	businessHandler  *handlers.BusinessHandler
	userHandler      *handlers.UserHandler
	productHandler   *handlers.ProductHandler
	inventoryHandler *handlers.InventoryHandler
	claimHandler     *handlers.ClaimHandler
	spotlightHandler *handlers.SpotlightHandler
	offerHandler     *handlers.OfferHandler
	reviewHandler    *handlers.ReviewHandler
	searchHandler    *handlers.SearchHandler
	categoryHandler  *handlers.CategoryHandler
}

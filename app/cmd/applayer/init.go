package applayer

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/handlers"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

func Init() application {
	a := application{}
	if err := a.InitDb(); err != nil {
		panic(err)
	}
	services := services.Init(a.GetDb())
	a.healthHandler = handlers.NewHealthHandler()
	a.authHandler = handlers.NewAuthHandler(services)
	a.businessHandler = handlers.NewBusinessHandler(services)
	a.userHandler = handlers.NewUserHandler(services)
	a.inventoryHandler = handlers.NewInventoryHandler(services)
	a.productHandler = handlers.NewProductHandler()
	a.claimHandler = handlers.NewClaimHandler()
	a.spotlightHandler = handlers.NewSpotlightHandler()
	a.offerHandler = handlers.NewOfferHandler()
	a.reviewHandler = handlers.NewReviewHandler()
	return a
}

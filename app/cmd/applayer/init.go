package applayer

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/atharvyadav96k/SPOTNEARR_API/handlers"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

func Init() application {
	if err := config.Load(); err != nil {
		panic(err)
	}
	a := application{}
	if err := a.InitCache(); err != nil {
		panic(err)
	}
	if err := a.InitDb(); err != nil {
		panic(err)
	}
	if err := a.InitCaptcha(); err != nil {
		panic(err)
	}
	services := services.Init(a.GetDb(), a.GetCache())
	a.healthHandler = handlers.NewHealthHandler()
	a.authHandler = handlers.NewAuthHandler(services)
	a.businessHandler = handlers.NewBusinessHandler(services)
	a.userHandler = handlers.NewUserHandler(services)
	a.inventoryHandler = handlers.NewInventoryHandler(services)
	a.productHandler = handlers.NewProductHandler(services)
	a.claimHandler = handlers.NewClaimHandler(services)
	a.searchHandler = handlers.NewSearchHandler(services)
	a.categoryHandler = handlers.NewCategoryHandler(services)
	a.spotlightHandler = handlers.NewSpotlightHandler()
	a.offerHandler = handlers.NewOfferHandler()
	a.reviewHandler = handlers.NewReviewHandler()
	return a
}

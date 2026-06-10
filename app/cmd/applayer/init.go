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
	svcs := services.Init(a.GetDb(), a.GetCache())
	a.healthHandler = handlers.NewHealthHandler()
	a.authHandler = handlers.NewAuthHandler(svcs)
	a.userHandler = handlers.NewUserHandler(svcs)
	a.claimHandler = handlers.NewClaimHandler(svcs)
	a.reviewHandler = handlers.NewReviewHandler(svcs)
	a.internalHandler = handlers.NewInternalHandler(svcs)
	return a
}

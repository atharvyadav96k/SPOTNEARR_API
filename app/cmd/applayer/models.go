package applayer

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/app"
	"github.com/atharvyadav96k/SPOTNEARR_API/handlers"
)

type application struct {
	app.App
	healthHandler   *handlers.Health
	authHandler     *handlers.AuthHandler
	userHandler     *handlers.UserHandler
	claimHandler    *handlers.ClaimHandler
	reviewHandler   *handlers.ReviewHandler
	socialHandler   *handlers.SocialHandler
	internalHandler *handlers.InternalHandler
	dealHandler     *handlers.DealHandler
}

package auth

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/atharvyadav96k/spotnearr-gcp/auth/user/register/applayer"
	"golang.org/x/crypto/bcrypt"
)

func Function(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := utils.ParseBody[models.User](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.PasswordHash), bcrypt.DefaultCost)
	hashed := string(hashedPassword)
	user.PasswordHash = &hashed
	created, err := app.RegisterUser(*user)
	if err != nil {
		utils.Conflict(w, err.Error())
		return
	}

	utils.Created(w, "user created successfully", created)
}

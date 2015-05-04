package v1

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/mholt/binding"
	"github.com/tommy351/app-studio-server/controller/common"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/util"
)

type tokenForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (form *tokenForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Email:    "email",
		&form.Password: "password",
	}
}

// TokenCreate handles POST /tokens.
func TokenCreate(res http.ResponseWriter, req *http.Request) {
	form := new(tokenForm)

	if common.BindForm(res, req, form) {
		return
	}

	var user *model.User

	if form.Email == "" {
		common.HandleAPIError(res, &util.APIError{
			Code:    util.RequiredError,
			Message: "Email is required.",
			Field:   "email",
		})
		return
	}

	if !govalidator.IsEmail(form.Email) {
		common.HandleAPIError(res, &util.APIError{
			Code:    util.EmailError,
			Message: "Email is invalid.",
			Field:   "email",
		})
		return
	}

	if user, _ = model.GetUserByEmail(form.Email); user == nil {
		common.HandleAPIError(res, &util.APIError{
			Field:   "email",
			Code:    util.UserNotFoundError,
			Message: "User does not exist.",
		})
		return
	}

	if err := user.Authenticate(form.Password); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	token := &model.Token{UserID: user.ID}

	if err := token.Save(); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.NoCacheHeader(res)
	common.RenderJSON(res, http.StatusCreated, token)
}

// TokenDestroy handles DELETE /tokens/:key.
func TokenDestroy(res http.ResponseWriter, req *http.Request) {
	key := common.GetParam(req, keyParam)

	token, err := model.GetTokenBase64(key)

	if err != nil {
		common.HandleAPIError(res, &util.APIError{
			Code:    util.TokenNotFoundError,
			Message: "Token does not exist.",
			Status:  http.StatusNotFound,
		})

		return
	}

	if err := token.Delete(); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

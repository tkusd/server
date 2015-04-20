package v1

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/util"
)

type TokenForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (form *TokenForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Email:    "email",
		&form.Password: "password",
	}
}

func TokenCreate(res http.ResponseWriter, req *http.Request) {
	form := new(TokenForm)

	if util.BindForm(res, req, form) {
		return
	}

	var user *model.User

	if form.Email == "" {
		util.HandleAPIError(res, &util.APIError{
			Code:    util.RequiredError,
			Message: "Email is required.",
			Field:   "email",
		})
		return
	}

	if !govalidator.IsEmail(form.Email) {
		util.HandleAPIError(res, &util.APIError{
			Code:    util.EmailError,
			Message: "Email is invalid.",
			Field:   "email",
		})
		return
	}

	if user, _ = model.GetUserByEmail(form.Email); user == nil {
		util.HandleAPIError(res, &util.APIError{
			Field:   "email",
			Code:    util.UserNotFoundError,
			Message: "User does not exist.",
		})
		return
	}

	if err := user.Authenticate(form.Password); err != nil {
		util.HandleAPIError(res, err)
		return
	}

	token := &model.Token{UserID: user.ID}

	if err := model.CreateToken(token); err != nil {
		util.HandleAPIError(res, err)
		return
	}

	util.NoCacheHeader(res)
	util.RenderJSON(res, http.StatusCreated, token)
}

func TokenDestroy(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	if key, ok := vars["key"]; ok {
		if token, err := model.GetTokenBase64(key); err != nil {
			util.HandleAPIError(res, &util.APIError{
				Code:    util.TokenNotFoundError,
				Message: "Token does not exist.",
				Status:  http.StatusNotFound,
			})
			return
		} else {
			if err := model.DeleteToken(token); err != nil {
				util.HandleAPIError(res, err)
				return
			}

			res.WriteHeader(http.StatusNoContent)
		}
	}
}

package v1

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/util"
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
func TokenCreate(res http.ResponseWriter, req *http.Request) error {
	form := new(tokenForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	var user *model.User

	if form.Email == "" {
		return &util.APIError{
			Code:    util.RequiredError,
			Message: "Email is required.",
			Field:   "email",
		}
	}

	if !govalidator.IsEmail(form.Email) {
		return &util.APIError{
			Code:    util.EmailError,
			Message: "Email is invalid.",
			Field:   "email",
		}
	}

	if user, _ = model.GetUserByEmail(form.Email); user == nil {
		return &util.APIError{
			Field:   "email",
			Code:    util.UserNotFoundError,
			Message: "User does not exist.",
		}
	}

	if err := user.Authenticate(form.Password); err != nil {
		return err
	}

	token := &model.Token{UserID: user.ID}

	if err := token.Save(); err != nil {
		return err
	}

	common.NoCacheHeader(res)
	common.APIResponse(res, req, http.StatusCreated, token)

	return nil
}

// TokenUpdate handles PUT /tokens/:key.
func TokenUpdate(res http.ResponseWriter, req *http.Request) error {
	token, err := GetToken(res, req)

	if err != nil {
		return err
	}

	if err := token.Save(); err != nil {
		return err
	}

	common.NoCacheHeader(res)
	common.APIResponse(res, req, http.StatusOK, token)

	return nil
}

// TokenDestroy handles DELETE /tokens/:key.
func TokenDestroy(res http.ResponseWriter, req *http.Request) error {
	token, err := GetToken(res, req)

	if err != nil {
		return err
	}

	if err := token.Delete(); err != nil {
		return err
	}

	res.WriteHeader(http.StatusNoContent)
	return nil
}

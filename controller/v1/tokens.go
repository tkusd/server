package v1

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
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
func TokenCreate(c *gin.Context) error {
	form := new(tokenForm)

	if err := common.BindForm(c, form); err != nil {
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

	common.NoCacheHeader(c)
	return common.APIResponse(c, http.StatusCreated, token)
}

// TokenShow handles GET /tokens/:key
func TokenShow(c *gin.Context) error {
	token, err := GetToken(c)

	if err != nil {
		return err
	}

	common.NoCacheHeader(c)
	return common.APIResponse(c, http.StatusOK, token.WithoutSecret())
}

// TokenDestroy handles DELETE /tokens/:key.
func TokenDestroy(c *gin.Context) error {
	token, err := GetToken(c)

	if err != nil {
		return err
	}

	if err := token.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

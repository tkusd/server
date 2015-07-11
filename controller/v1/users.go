package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/util"
)

type userForm struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	OldPassword *string `json:"old_password"`
}

func (form *userForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:        "name",
		&form.Email:       "email",
		&form.Password:    "password",
		&form.OldPassword: "old_password",
	}
}

// UserCreate handles POST /users.
func UserCreate(c *gin.Context) error {
	form := new(userForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	if form.Name == nil {
		return &util.APIError{
			Field:   "name",
			Code:    util.RequiredError,
			Message: "Name is required.",
		}
	}

	if form.Email == nil {
		return &util.APIError{
			Field:   "email",
			Code:    util.RequiredError,
			Message: "Email is required.",
		}
	}

	if form.Password == nil {
		return &util.APIError{
			Field:   "password",
			Code:    util.RequiredError,
			Message: "Password is required.",
		}
	}

	user := &model.User{
		Name:  *form.Name,
		Email: *form.Email,
	}

	if err := user.GeneratePassword(*form.Password); err != nil {
		return err
	}

	user.SetActivated(false)

	if err := user.Save(); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, user)
}

// UserShow handles GET /users/:user_id.
func UserShow(c *gin.Context) error {
	user, err := GetUser(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, user.ID); err == nil {
		return common.APIResponse(c, http.StatusOK, user)
	}

	return common.APIResponse(c, http.StatusOK, user.PublicProfile())
}

// UserUpdate handles PUT /users/:user_id.
func UserUpdate(c *gin.Context) error {
	form := new(userForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	user, err := GetUser(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, user.ID); err != nil {
		return err
	}

	if form.Name != nil {
		user.Name = *form.Name
	}

	if form.Password != nil {
		if form.OldPassword == nil {
			return &util.APIError{
				Field:   "old_password",
				Code:    util.RequiredError,
				Message: "Current password is required.",
			}
		}

		if err := user.Authenticate(*form.OldPassword); err != nil {
			return &util.APIError{
				Field:   "old_password",
				Code:    util.WrongPasswordError,
				Message: "Password is wrong.",
			}
		}

		if err := user.GeneratePassword(*form.Password); err != nil {
			return err
		}
	}

	if form.Email != nil {
		user.Email = *form.Email
	}

	if err := user.Save(); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, user)
}

// UserDestroy handles DELETE /users/:user_id.
func UserDestroy(c *gin.Context) error {
	user, err := GetUser(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, user.ID); err != nil {
		return err
	}

	if err := user.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

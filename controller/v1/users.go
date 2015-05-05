package v1

import (
	"net/http"

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
func UserCreate(res http.ResponseWriter, req *http.Request) {
	form := new(userForm)

	if common.BindForm(res, req, form) {
		return
	}

	if form.Name == nil {
		common.HandleAPIError(res, req, &util.APIError{
			Field:   "name",
			Code:    util.RequiredError,
			Message: "Name is required.",
		})
		return
	}

	if form.Email == nil {
		common.HandleAPIError(res, req, &util.APIError{
			Field:   "email",
			Code:    util.RequiredError,
			Message: "Email is required.",
		})
		return
	}

	if form.Password == nil {
		common.HandleAPIError(res, req, &util.APIError{
			Field:   "password",
			Code:    util.RequiredError,
			Message: "Password is required.",
		})
		return
	}

	user := &model.User{
		Name:  *form.Name,
		Email: *form.Email,
	}

	if err := user.GeneratePassword(*form.Password); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	user.SetActivated(false)

	if err := user.Save(); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	common.APIResponse(res, req, http.StatusCreated, user)
}

// UserShow handles GET /users/:user_id.
func UserShow(res http.ResponseWriter, req *http.Request) {
	user, err := GetUser(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	err = CheckUserPermission(res, req, user.ID)

	if err == nil {
		common.APIResponse(res, req, http.StatusOK, user)
	} else {
		common.APIResponse(res, req, http.StatusOK, user.PublicProfile())
	}
}

// UserUpdate handles PUT /users/:user_id.
func UserUpdate(res http.ResponseWriter, req *http.Request) {
	form := new(userForm)

	if common.BindForm(res, req, form) {
		return
	}

	user, err := GetUser(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckUserPermission(res, req, user.ID); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if form.Name != nil {
		user.Name = *form.Name
	}

	if form.Password != nil {
		if form.OldPassword == nil {
			common.HandleAPIError(res, req, &util.APIError{
				Field:   "old_password",
				Code:    util.RequiredError,
				Message: "Current password is required.",
			})
			return
		}

		if err := user.Authenticate(*form.OldPassword); err != nil {
			common.HandleAPIError(res, req, &util.APIError{
				Field:   "old_password",
				Code:    util.WrongPasswordError,
				Message: "Password is wrong.",
			})
			return
		}

		if err := user.GeneratePassword(*form.Password); err != nil {
			common.HandleAPIError(res, req, err)
			return
		}
	}

	if form.Email != nil {
		user.Email = *form.Email
	}

	if err := user.Save(); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	common.APIResponse(res, req, http.StatusOK, user)
}

// UserDestroy handles DELETE /users/:user_id.
func UserDestroy(res http.ResponseWriter, req *http.Request) {
	user, err := GetUser(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckUserPermission(res, req, user.ID); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := user.Delete(); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

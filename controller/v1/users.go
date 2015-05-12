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
func UserCreate(res http.ResponseWriter, req *http.Request) error {
	form := new(userForm)

	if err := common.BindForm(res, req, form); err != nil {
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

	common.APIResponse(res, req, http.StatusCreated, user)
	return nil
}

// UserShow handles GET /users/:user_id.
func UserShow(res http.ResponseWriter, req *http.Request) error {
	user, err := GetUser(res, req)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, user.ID); err == nil {
		common.APIResponse(res, req, http.StatusOK, user)
	} else {
		common.APIResponse(res, req, http.StatusOK, user.PublicProfile())
	}

	return nil
}

// UserUpdate handles PUT /users/:user_id.
func UserUpdate(res http.ResponseWriter, req *http.Request) error {
	form := new(userForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	user, err := GetUser(res, req)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, user.ID); err != nil {
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

	common.APIResponse(res, req, http.StatusOK, user)
	return nil
}

// UserDestroy handles DELETE /users/:user_id.
func UserDestroy(res http.ResponseWriter, req *http.Request) error {
	user, err := GetUser(res, req)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, user.ID); err != nil {
		return err
	}

	if err := user.Delete(); err != nil {
		return err
	}

	res.WriteHeader(http.StatusNoContent)
	return nil
}

package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/util"
)

type UserForm struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	OldPassword *string `json:"old_password"`
}

func (form *UserForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:        "name",
		&form.Email:       "email",
		&form.Password:    "password",
		&form.OldPassword: "old_password",
	}
}

func GetUser(res http.ResponseWriter, req *http.Request) *model.User {
	vars := mux.Vars(req)

	if id, ok := vars["id"]; ok {
		if user, err := model.GetUser(util.ParseUUID(id)); err == nil {
			return user
		}
	}

	util.HandleAPIError(res, &util.APIError{
		Code:    util.UserNotFoundError,
		Message: "User not found.",
		Status:  http.StatusNotFound,
	})

	return nil
}

func UserCreate(res http.ResponseWriter, req *http.Request) {
	form := new(UserForm)

	if util.BindForm(res, req, form) {
		return
	}

	if form.Name == nil {
		util.HandleAPIError(res, &util.APIError{
			Field:   "name",
			Code:    util.RequiredError,
			Message: "Name is required.",
		})
		return
	}

	if form.Email == nil {
		util.HandleAPIError(res, &util.APIError{
			Field:   "email",
			Code:    util.RequiredError,
			Message: "Email is required.",
		})
		return
	}

	if form.Password == nil {
		util.HandleAPIError(res, &util.APIError{
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
		util.HandleAPIError(res, err)
		return
	}

	user.SetActivated(false)

	if err := model.CreateUser(user); err != nil {
		util.HandleAPIError(res, err)
		return
	}

	util.RenderJSON(res, http.StatusCreated, user)
}

func UserShow(res http.ResponseWriter, req *http.Request) {
	if user := GetUser(res, req); user != nil {
		util.RenderJSON(res, http.StatusOK, user)
	}
}

func UserUpdate(res http.ResponseWriter, req *http.Request) {
	form := new(UserForm)

	if util.BindForm(res, req, form) {
		return
	}

	var user *model.User

	if user = GetUser(res, req); user == nil {
		return
	}

	if form.Name != nil {
		user.Name = *form.Name
	}

	if form.Password != nil {
		if form.OldPassword == nil {
			util.HandleAPIError(res, &util.APIError{
				Field:   "old_password",
				Code:    util.RequiredError,
				Message: "Current password is required.",
			})
			return
		}

		if err := user.Authenticate(*form.OldPassword); err != nil {
			util.HandleAPIError(res, &util.APIError{
				Field:   "old_password",
				Code:    util.WrongPasswordError,
				Message: "Password is wrong.",
			})
			return
		}

		if err := user.GeneratePassword(*form.Password); err != nil {
			util.HandleAPIError(res, err)
			return
		}
	}

	if form.Email != nil {
		user.Email = *form.Email
	}

	if err := model.UpdateUser(user); err != nil {
		util.HandleAPIError(res, err)
		return
	}

	util.RenderJSON(res, http.StatusOK, user)
}

func UserDestroy(res http.ResponseWriter, req *http.Request) {
	var user *model.User

	if user = GetUser(res, req); user == nil {
		return
	}

	if err := model.DeleteUser(user); err != nil {
		util.HandleAPIError(res, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

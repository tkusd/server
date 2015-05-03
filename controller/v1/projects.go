package v1

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/tommy351/app-studio-server/controller/common"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/model/types"
	"github.com/tommy351/app-studio-server/util"
)

// ProjectList handles GET /users/:user_id/projects.
func ProjectList(res http.ResponseWriter, req *http.Request) {
	userID := types.ParseUUID(common.GetParam(req, "user_id"))
	option := &model.ProjectQueryOption{
		UserID: userID,
	}

	if err := CheckUserPermission(res, req, userID); err == nil {
		option.Private = true
	}

	list, err := model.GetProjectList(option)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.RenderJSON(res, http.StatusOK, list)
}

type projectForm struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPrivate   *bool   `json:"is_private"`
}

func (form *projectForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Title:       "title",
		&form.Description: "description",
		&form.IsPrivate:   "is_private",
	}
}

// ProjectCreate handles POST /users/:user_id/projects.
func ProjectCreate(res http.ResponseWriter, req *http.Request) {
	userID := types.ParseUUID(common.GetParam(req, "user_id"))

	if err := CheckUserPermission(res, req, userID); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	form := new(projectForm)

	if common.BindForm(res, req, form) {
		return
	}

	if form.Title == nil {
		common.HandleAPIError(res, &util.APIError{
			Field:   "title",
			Code:    util.RequiredError,
			Message: "Title is required.",
		})
		return
	}

	if form.Description == nil {
		*form.Description = ""
	}

	if form.IsPrivate == nil {
		*form.IsPrivate = false
	}

	project := &model.Project{
		Title:       *form.Title,
		Description: *form.Description,
		IsPrivate:   *form.IsPrivate,
		UserID:      userID,
	}

	if err := project.Save(); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.RenderJSON(res, http.StatusCreated, project)
}

// ProjectShow handles GET /projects/:project_id.
func ProjectShow(res http.ResponseWriter, req *http.Request) {
	project, err := GetProject(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	token, err := GetToken(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if project.IsPrivate && !project.UserID.Equal(token.UserID) {
		common.HandleAPIError(res, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this project.",
			Status:  http.StatusForbidden,
		})
		return
	}

	common.RenderJSON(res, http.StatusOK, project)
}

// ProjectUpdate handles PUT /projects/:project_id.
func ProjectUpdate(res http.ResponseWriter, req *http.Request) {
	form := new(projectForm)

	if common.BindForm(res, req, form) {
		return
	}

	project, err := GetProject(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err = CheckUserPermission(res, req, project.UserID); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if form.Title != nil {
		project.Title = *form.Title
	}

	if form.Description != nil {
		project.Description = *form.Description
	}

	if form.IsPrivate != nil {
		project.IsPrivate = *form.IsPrivate
	}

	if err := project.Save(); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.RenderJSON(res, http.StatusOK, project)
}

// ProjectDestroy handles DELETE /projects/:project_id.
func ProjectDestroy(res http.ResponseWriter, req *http.Request) {
	project, err := GetProject(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := project.Delete(); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

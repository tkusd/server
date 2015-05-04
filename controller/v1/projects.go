package v1

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

// ProjectList handles GET /users/:user_id/projects.
func ProjectList(res http.ResponseWriter, req *http.Request) {
	userID := types.ParseUUID(common.GetParam(req, userIDParam))
	option := &model.ProjectQueryOption{
		UserID: &userID,
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
	Title       *string       `json:"title"`
	Description *string       `json:"description"`
	IsPrivate   *bool         `json:"is_private"`
	Elements    *[]types.UUID `json:"elements"`
}

func (form *projectForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Title:       "title",
		&form.Description: "description",
		&form.IsPrivate:   "is_private",
		&form.Elements:    "elements",
	}
}

func saveProject(form *projectForm, project *model.Project) error {
	if form.Title != nil {
		project.Title = *form.Title
	}

	if form.Description != nil {
		project.Description = *form.Description
	}

	if form.IsPrivate != nil {
		project.IsPrivate = *form.IsPrivate
	}

	return project.Save()
}

// ProjectCreate handles POST /users/:user_id/projects.
func ProjectCreate(res http.ResponseWriter, req *http.Request) {
	userID := types.ParseUUID(common.GetParam(req, userIDParam))

	if err := CheckUserPermission(res, req, userID); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	form := new(projectForm)

	if common.BindForm(res, req, form) {
		return
	}

	project := &model.Project{UserID: userID}

	if err := saveProject(form, project); err != nil {
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

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := saveProject(form, project); err != nil {
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

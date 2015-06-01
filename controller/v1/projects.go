package v1

import (
	"net/http"
	"strconv"

	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/util"
)

// ProjectList handles GET /users/:user_id/projects.
func ProjectList(res http.ResponseWriter, req *http.Request) error {
	userID, err := GetIDParam(req, userIDParam)

	if err != nil {
		return err
	}

	option := &model.ProjectQueryOption{
		UserID: userID,
	}

	if limit := req.URL.Query().Get("limit"); limit != "" {
		if i, err := strconv.Atoi(limit); err == nil {
			option.Limit = i
		}
	}

	if offset := req.URL.Query().Get("offset"); offset != "" {
		if i, err := strconv.Atoi(offset); err == nil {
			option.Offset = i
		}
	}

	if order := req.URL.Query().Get("order"); order != "" {
		option.Order = order
	}

	if err := CheckUserPermission(res, req, *userID); err == nil {
		option.Private = true
	}

	list, err := model.GetProjectList(option)

	if err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusOK, list)
	return nil
}

type projectForm struct {
	Title       *string                  `json:"title"`
	Description *string                  `json:"description"`
	IsPrivate   *bool                    `json:"is_private"`
	Elements    *[]model.ElementTreeItem `json:"elements"`
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
func ProjectCreate(res http.ResponseWriter, req *http.Request) error {
	userID, err := GetIDParam(req, userIDParam)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, *userID); err != nil {
		return err
	}

	form := new(projectForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	project := &model.Project{UserID: *userID}

	if err := saveProject(form, project); err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusCreated, project)
	return nil
}

func getProjectWithOwner(res http.ResponseWriter, req *http.Request) (*model.Project, error) {
	id, err := GetIDParam(req, projectIDParam)

	if err != nil {
		return nil, err
	}

	project, err := model.GetProjectWithOwner(*id)

	if err != nil {
		return nil, err
	}

	token, _ := CheckToken(res, req)

	if project.IsPrivate && !project.UserID.Equal(token.UserID) {
		return nil, &util.APIError{
			Code:    util.UserForbiddenError,
			Message: "You are forbidden to access this project.",
			Status:  http.StatusForbidden,
		}
	}

	return project, nil
}

// ProjectShow handles GET /projects/:project_id.
func ProjectShow(res http.ResponseWriter, req *http.Request) error {
	project, err := getProjectWithOwner(res, req)

	if err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusOK, project)
	return nil
}

// ProjectUpdate handles PUT /projects/:project_id.
func ProjectUpdate(res http.ResponseWriter, req *http.Request) error {
	form := new(projectForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	project, err := GetProject(res, req)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		return err
	}

	if err := saveProject(form, project); err != nil {
		return err
	}

	if form.Elements != nil {
		option := &model.ElementQueryOption{
			ProjectID: &project.ID,
		}

		if err := model.UpdateElementOrder(option, *form.Elements); err != nil {
			return err
		}
	}

	common.APIResponse(res, req, http.StatusOK, project)
	return nil
}

// ProjectDestroy handles DELETE /projects/:project_id.
func ProjectDestroy(res http.ResponseWriter, req *http.Request) error {
	project, err := GetProject(res, req)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		return err
	}

	if err := project.Delete(); err != nil {
		return err
	}

	res.WriteHeader(http.StatusNoContent)
	return nil
}

func ProjectFull(res http.ResponseWriter, req *http.Request) error {
	project, err := getProjectWithOwner(res, req)

	if err != nil {
		return err
	}

	elements, err := model.GetElementList(&model.ElementQueryOption{
		ProjectID: &project.ID,
	})

	if err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusOK, struct {
		*model.Project
		Elements []*model.Element `json:"elements"`
	}{
		Project:  project,
		Elements: elements,
	})
	return nil
}

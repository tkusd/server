package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

// ProjectList handles GET /users/:user_id/projects.
func ProjectList(c *gin.Context) error {
	userID, err := GetIDParam(c, userIDParam)

	if err != nil {
		return err
	}

	option := &model.ProjectQueryOption{
		UserID: userID,
	}

	if limit := c.Query("limit"); limit != "" {
		if i, err := strconv.Atoi(limit); err == nil {
			option.Limit = i
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if i, err := strconv.Atoi(offset); err == nil {
			option.Offset = i
		}
	}

	if order := c.Query("order"); order != "" {
		option.Order = order
	} else {
		option.Order = "-created_at"
	}

	if err := CheckUserPermission(c, *userID); err == nil {
		option.Private = true
	}

	list, err := model.GetProjectList(option)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, list)
}

type projectForm struct {
	Title       *string       `json:"title"`
	Description *string       `json:"description"`
	IsPrivate   *bool         `json:"is_private"`
	Elements    *[]types.UUID `json:"elements"`
	MainScreen  *types.UUID   `json:"main_screen"`
	Theme       *string       `json:"theme"`
}

func (form *projectForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Title:       "title",
		&form.Description: "description",
		&form.IsPrivate:   "is_private",
		&form.Elements:    "elements",
		&form.MainScreen:  "main_screen",
		&form.Theme:       "theme",
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

	if form.Theme != nil {
		project.Theme = *form.Theme
	}

	return project.Save()
}

// ProjectCreate handles POST /users/:user_id/projects.
func ProjectCreate(c *gin.Context) error {
	userID, err := GetIDParam(c, userIDParam)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, *userID); err != nil {
		return err
	}

	form := new(projectForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	project := &model.Project{UserID: *userID}

	if err := saveProject(form, project); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, project)
}

func getProjectWithOwner(c *gin.Context) (*model.Project, error) {
	id, err := GetIDParam(c, projectIDParam)

	if err != nil {
		return nil, err
	}

	project, err := model.GetProjectWithOwner(*id)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, &util.APIError{
			Code:    util.ProjectNotFoundError,
			Message: "Project not found.",
			Status:  http.StatusNotFound,
		}
	}

	token, _ := CheckToken(c)

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
func ProjectShow(c *gin.Context) error {
	project, err := getProjectWithOwner(c)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, project)
}

// ProjectUpdate handles PUT /projects/:project_id.
func ProjectUpdate(c *gin.Context) error {
	form := new(projectForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	project, err := GetProject(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, project.UserID); err != nil {
		return err
	}

	if form.MainScreen != nil {
		project.MainScreen = *form.MainScreen
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

	return common.APIResponse(c, http.StatusOK, project)
}

// ProjectDestroy handles DELETE /projects/:project_id.
func ProjectDestroy(c *gin.Context) error {
	project, err := GetProject(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, project.UserID); err != nil {
		return err
	}

	if err := project.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

func ProjectFull(c *gin.Context) error {
	project, err := getProjectWithOwner(c)

	if err != nil {
		return err
	}

	option := parseElementListQueryOption(c)
	option.ProjectID = &project.ID
	option.WithEvents = true

	elements, err := model.GetElementList(option)

	if err != nil {
		return err
	}

	assets, err := model.GetAssetList(project.ID)

	if err != nil {
		return err
	}

	actions, err := model.GetActionList(project.ID)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, struct {
		*model.Project
		Elements []*model.Element `json:"elements"`
		Assets   []*model.Asset   `json:"assets"`
		Actions  []*model.Action  `json:"actions"`
	}{
		Project:  project,
		Elements: elements,
		Assets:   assets,
		Actions:  actions,
	})
}

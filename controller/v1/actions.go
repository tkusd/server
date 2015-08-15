package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
)

func ActionList(c *gin.Context) error {
	projectID, err := GetIDParam(c, projectIDParam)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, *projectID, false); err != nil {
		return err
	}

	list, err := model.GetActionList(*projectID)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, list)
}

type actionForm struct {
	Name   *string           `json:"name"`
	Action *string           `json:"action"`
	Data   *types.JSONObject `json:"data"`
}

func (form *actionForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:   "name",
		&form.Action: "action",
		&form.Data:   "data",
	}
}

func saveAction(form *actionForm, action *model.Action) error {
	if form.Name != nil {
		action.Name = *form.Name
	}

	if form.Action != nil {
		action.Action = *form.Action
	}

	if form.Data != nil {
		action.Data = *form.Data
	}

	return action.Save()
}

func ActionCreate(c *gin.Context) error {
	project, err := GetProject(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, project.ID, true); err != nil {
		return err
	}

	form := new(actionForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	action := &model.Action{
		ProjectID: project.ID,
	}

	if err := saveAction(form, action); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, action)
}

func ActionShow(c *gin.Context) error {
	action, err := GetAction(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, action.ProjectID, false); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, action)
}

func ActionUpdate(c *gin.Context) error {
	form := new(actionForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	action, err := GetAction(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, action.ProjectID, true); err != nil {
		return err
	}

	if err := saveAction(form, action); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, action)
}

func ActionDestroy(c *gin.Context) error {
	action, err := GetAction(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, action.ProjectID, true); err != nil {
		return err
	}

	if err := action.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

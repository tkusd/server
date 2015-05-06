package v1

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
)

// ElementList handles GET /projects/:project_id/elements.
func ElementList(res http.ResponseWriter, req *http.Request) {
	projectID := types.ParseUUID(common.GetParam(req, projectIDParam))
	option := &model.ElementQueryOption{
		ProjectID: &projectID,
	}

	if err := CheckProjectPermission(res, req, projectID, false); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	list, err := model.GetElementList(option)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	common.APIResponse(res, req, http.StatusOK, list)
}

func ChildElementList(res http.ResponseWriter, req *http.Request) {
	element, err := GetElement(res, req)
	option := &model.ElementQueryOption{
		ElementID: &element.ID,
	}

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, false); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	list, err := model.GetElementList(option)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	common.APIResponse(res, req, http.StatusOK, list)
}

type elementForm struct {
	Name       *string            `json:"name"`
	Type       *types.ElementType `json:"type"`
	Attributes *types.JSONObject  `json:"attributes"`
	ParentID   *types.UUID        `json:"parent_id"`
	Elements   *[]types.UUID      `json:"elements"`
	Order      *int               `json:"order"`
}

func (form *elementForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:       "name",
		&form.Type:       "type",
		&form.Attributes: "attributes",
		&form.ParentID:   "parent_id",
		&form.Elements:   "elements",
		&form.Order:      "order",
	}
}

func saveElement(form *elementForm, element *model.Element) error {
	if form.Name != nil {
		element.Name = *form.Name
	}

	if form.Type != nil {
		element.Type = *form.Type
	}

	if form.Attributes != nil {
		element.Attributes = *form.Attributes
	}

	return element.Save()
}

// ElementCreate handles POST /projects/:project_id/elements.
func ElementCreate(res http.ResponseWriter, req *http.Request) {
	project, err := GetProject(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	form := new(elementForm)

	if common.BindForm(res, req, form) {
		return
	}

	element := &model.Element{ProjectID: project.ID}

	if err := saveElement(form, element); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	common.APIResponse(res, req, http.StatusCreated, element)
}

func ChildElementCreate(res http.ResponseWriter, req *http.Request) {
	parent, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckProjectPermission(res, req, parent.ProjectID, true); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	form := new(elementForm)

	if common.BindForm(res, req, form) {
		return
	}

	element := &model.Element{
		ProjectID: parent.ProjectID,
		ElementID: parent.ID,
	}

	if err := saveElement(form, element); err != nil {
		common.HandleAPIError(res, req, err)
	}

	common.APIResponse(res, req, http.StatusCreated, element)
}

// ElementShow handles GET /elements/:element_id.
func ElementShow(res http.ResponseWriter, req *http.Request) {
	element, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, false); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	common.APIResponse(res, req, http.StatusOK, element)
}

// ElementUpdate handles PUT /elements/:element_id.
func ElementUpdate(res http.ResponseWriter, req *http.Request) {
	form := new(elementForm)

	if common.BindForm(res, req, form) {
		return
	}

	element, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, true); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if form.ParentID != nil {
		element.ElementID = *form.ParentID
	}

	if err := saveElement(form, element); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if form.Elements != nil {
		// TODO: update element order
		/*
			if err := element.UpdateOrder(*form.Elements); err != nil {
				common.HandleAPIError(res, req, err)
				return
			}*/
	}

	common.APIResponse(res, req, http.StatusOK, element)
}

// ElementDestroy handles DELETE /elements/:element_id.
func ElementDestroy(res http.ResponseWriter, req *http.Request) {
	element, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, true); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	if err := element.Delete(); err != nil {
		common.HandleAPIError(res, req, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

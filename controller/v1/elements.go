package v1

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/tommy351/app-studio-server/controller/common"
	"github.com/tommy351/app-studio-server/model"
	"github.com/tommy351/app-studio-server/model/types"
)

// ElementList handles GET /projects/:project_id/elements.
func ElementList(res http.ResponseWriter, req *http.Request) {
	projectID := types.ParseUUID(common.GetParam(req, projectIDParam))
	option := &model.ElementQueryOption{
		ProjectID: &projectID,
	}

	if err := CheckProjectPermission(res, req, projectID, false); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	list, err := model.GetElementList(option)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.RenderJSON(res, http.StatusOK, list)
}

type elementForm struct {
	Name       *string            `json:"name"`
	Type       *types.ElementType `json:"type"`
	Attributes *types.JSONObject  `json:"attributes"`
	ParentID   *types.UUID        `json:"parent_id"`
	Elements   *[]types.UUID      `json:"elements"`
}

func (form *elementForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:       "name",
		&form.Type:       "type",
		&form.Attributes: "attributes",
		&form.ParentID:   "parent_id",
		&form.Elements:   "elements",
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
		common.HandleAPIError(res, err)
		return
	}

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	form := new(elementForm)

	if common.BindForm(res, req, form) {
		return
	}

	element := &model.Element{ProjectID: project.ID}

	if err := saveElement(form, element); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.RenderJSON(res, http.StatusCreated, element)
}

// ElementShow handles GET /elements/:element_id.
func ElementShow(res http.ResponseWriter, req *http.Request) {
	element, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, false); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	common.RenderJSON(res, http.StatusOK, element)
}

// ElementUpdate handles PUT /elements/:element_id.
func ElementUpdate(res http.ResponseWriter, req *http.Request) {
	form := new(elementForm)

	if common.BindForm(res, req, form) {
		return
	}

	element, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, true); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if form.ParentID != nil {
		element.ElementID = form.ParentID
	}

	if err := saveElement(form, element); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if form.Elements != nil {
		if err := element.UpdateOrder(*form.Elements); err != nil {
			common.HandleAPIError(res, err)
			return
		}
	}

	common.RenderJSON(res, http.StatusOK, element)
}

// ElementDestroy handles DELETE /elements/:element_id.
func ElementDestroy(res http.ResponseWriter, req *http.Request) {
	element, err := GetElement(res, req)

	if err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, true); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	if err := element.Delete(); err != nil {
		common.HandleAPIError(res, err)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

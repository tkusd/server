package v1

import (
	"net/http"
	"strconv"

	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
)

func parseElementListQueryOption(req *http.Request) *model.ElementQueryOption {
	var option model.ElementQueryOption
	query := req.URL.Query()

	if flat, ok := query["flat"]; ok {
		if len(flat) == 0 || (flat[0] != "false" && flat[0] != "0") {
			option.Flat = true
		}
	}

	if depth := query.Get("depth"); depth != "" {
		if i, err := strconv.Atoi(depth); err == nil {
			option.Depth = uint(i)
		}
	}

	// This feature is disabled temporarily since I can't control the returned fields.
	/*
		if sel := query.Get("select"); sel != "" {
			option.Select = util.SplitAndTrim(sel, ",")
		}*/

	return &option
}

// ElementList handles GET /projects/:project_id/elements.
func ElementList(res http.ResponseWriter, req *http.Request) error {
	projectID, err := GetIDParam(req, projectIDParam)

	if err != nil {
		return err
	}

	option := parseElementListQueryOption(req)
	option.ProjectID = projectID

	if err := CheckProjectPermission(res, req, *projectID, false); err != nil {
		return err
	}

	list, err := model.GetElementList(option)

	if err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusOK, list)
	return nil
}

func ChildElementList(res http.ResponseWriter, req *http.Request) error {
	element, err := GetElement(res, req)
	option := parseElementListQueryOption(req)
	option.ElementID = &element.ID

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, false); err != nil {
		return err
	}

	list, err := model.GetElementList(option)

	if err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusOK, list)
	return nil
}

type elementForm struct {
	Name       *string                  `json:"name"`
	Type       *types.ElementType       `json:"type"`
	Attributes *types.JSONObject        `json:"attributes"`
	Elements   *[]model.ElementTreeItem `json:"elements"`
}

func (form *elementForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:       "name",
		&form.Type:       "type",
		&form.Attributes: "attributes",
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
func ElementCreate(res http.ResponseWriter, req *http.Request) error {
	project, err := GetProject(res, req)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(res, req, project.UserID); err != nil {
		return err
	}

	form := new(elementForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	element := &model.Element{ProjectID: project.ID}

	if err := saveElement(form, element); err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusCreated, element)
	return nil
}

func ChildElementCreate(res http.ResponseWriter, req *http.Request) error {
	parent, err := GetElement(res, req)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(res, req, parent.ProjectID, true); err != nil {
		return err
	}

	form := new(elementForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	element := &model.Element{
		ProjectID: parent.ProjectID,
		ElementID: parent.ID,
	}

	if err := saveElement(form, element); err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusCreated, element)
	return nil
}

// ElementShow handles GET /elements/:element_id.
func ElementShow(res http.ResponseWriter, req *http.Request) error {
	element, err := GetElement(res, req)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, false); err != nil {
		return err
	}

	common.APIResponse(res, req, http.StatusOK, element)
	return nil
}

// ElementUpdate handles PUT /elements/:element_id.
func ElementUpdate(res http.ResponseWriter, req *http.Request) error {
	form := new(elementForm)

	if err := common.BindForm(res, req, form); err != nil {
		return err
	}

	element, err := GetElement(res, req)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, true); err != nil {
		return err
	}

	if err := saveElement(form, element); err != nil {
		return err
	}

	if form.Elements != nil {
		option := &model.ElementQueryOption{
			ElementID: &element.ID,
		}

		if err := model.UpdateElementOrder(option, *form.Elements); err != nil {
			return err
		}
	}

	common.APIResponse(res, req, http.StatusOK, element)
	return nil
}

// ElementDestroy handles DELETE /elements/:element_id.
func ElementDestroy(res http.ResponseWriter, req *http.Request) error {
	element, err := GetElement(res, req)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(res, req, element.ProjectID, true); err != nil {
		return err
	}

	if err := element.Delete(); err != nil {
		return err
	}

	res.WriteHeader(http.StatusNoContent)
	return nil
}

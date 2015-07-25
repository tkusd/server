package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
)

func parseElementListQueryOption(c *gin.Context) *model.ElementQueryOption {
	var option model.ElementQueryOption

	if common.QueryExist(c, "flat") {
		flat := c.Query("flat")

		if flat == "" || (flat != "false" && flat != "0") {
			option.Flat = true
		}
	}

	if depth := c.Query("depth"); depth != "" {
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
func ElementList(c *gin.Context) error {
	projectID, err := GetIDParam(c, projectIDParam)

	if err != nil {
		return err
	}

	option := parseElementListQueryOption(c)
	option.ProjectID = projectID

	if err := CheckProjectPermission(c, *projectID, false); err != nil {
		return err
	}

	list, err := model.GetElementList(option)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, list)
}

func ChildElementList(c *gin.Context) error {
	element, err := GetElement(c)

	if err != nil {
		return err
	}

	option := parseElementListQueryOption(c)
	option.ElementID = &element.ID

	if err := CheckProjectPermission(c, element.ProjectID, false); err != nil {
		return err
	}

	list, err := model.GetElementList(option)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, list)
}

type elementForm struct {
	Name       *string                  `json:"name"`
	Type       *string                  `json:"type"`
	Attributes *types.JSONObject        `json:"attributes"`
	Styles     *types.JSONObject        `json:"styles"`
	Elements   *[]model.ElementTreeItem `json:"elements"`
	IsVisible  *bool                    `json:"is_visible"`
}

func (form *elementForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:       "name",
		&form.Type:       "type",
		&form.Attributes: "attributes",
		&form.Styles:     "styles",
		&form.Elements:   "elements",
		&form.IsVisible:  "is_visible",
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

	if form.Styles != nil {
		element.Styles = *form.Styles
	}

	if form.IsVisible != nil {
		element.IsVisible = *form.IsVisible
	}

	return element.Save()
}

// ElementCreate handles POST /projects/:project_id/elements.
func ElementCreate(c *gin.Context) error {
	project, err := GetProject(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, project.UserID); err != nil {
		return err
	}

	form := new(elementForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	element := &model.Element{
		ProjectID: project.ID,
		IsVisible: true,
	}

	if err := saveElement(form, element); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, element)
}

func ChildElementCreate(c *gin.Context) error {
	parent, err := GetElement(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, parent.ProjectID, true); err != nil {
		return err
	}

	form := new(elementForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	element := &model.Element{
		ProjectID: parent.ProjectID,
		ElementID: parent.ID,
		IsVisible: true,
	}

	if err := saveElement(form, element); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, element)
}

// ElementShow handles GET /elements/:element_id.
func ElementShow(c *gin.Context) error {
	element, err := GetElement(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, element.ProjectID, false); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, element)
}

// ElementUpdate handles PUT /elements/:element_id.
func ElementUpdate(c *gin.Context) error {
	form := new(elementForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	element, err := GetElement(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, element.ProjectID, true); err != nil {
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

	return common.APIResponse(c, http.StatusOK, element)
}

// ElementDestroy handles DELETE /elements/:element_id.
func ElementDestroy(c *gin.Context) error {
	element, err := GetElement(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, element.ProjectID, true); err != nil {
		return err
	}

	if err := element.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

func ElementFull(c *gin.Context) error {
	element, err := GetElement(c)

	if err != nil {
		return err
	}

	option := parseElementListQueryOption(c)
	option.ElementID = &element.ID

	if err := CheckProjectPermission(c, element.ProjectID, false); err != nil {
		return err
	}

	list, err := model.GetElementList(option)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, struct {
		*model.Element
		Elements []*model.Element `json:"elements"`
	}{
		Element:  element,
		Elements: list,
	})
}

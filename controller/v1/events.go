package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
)

func EventList(c *gin.Context) error {
	elementID, err := GetIDParam(c, elementIDParam)

	if err != nil {
		return err
	}

	projectID := model.GetProjectIDForElement(*elementID)

	if err := CheckProjectPermission(c, projectID, false); err != nil {
		return err
	}

	list, err := model.GetEventList(*elementID)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, list)
}

type eventForm struct {
	Event    *string     `json:"event"`
	ActionID *types.UUID `json:"action_id"`
}

func (form *eventForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Event:    "event",
		&form.ActionID: "action_id",
	}
}

func saveEvent(form *eventForm, event *model.Event) error {
	if form.Event != nil {
		event.Event = *form.Event
	}

	if form.ActionID != nil {
		event.ActionID = *form.ActionID
	}

	return event.Save()
}

func EventCreate(c *gin.Context) error {
	elementID, err := GetIDParam(c, elementIDParam)

	if err != nil {
		return err
	}

	projectID := model.GetProjectIDForElement(*elementID)

	if err := CheckProjectPermission(c, projectID, true); err != nil {
		return err
	}

	form := new(eventForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	event := &model.Event{
		ElementID: *elementID,
	}

	if err := saveEvent(form, event); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, event)
}

func EventShow(c *gin.Context) error {
	event, err := GetEvent(c)

	if err != nil {
		return err
	}

	projectID := model.GetProjectIDForElement(event.ElementID)

	if err := CheckProjectPermission(c, projectID, false); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, event)
}

func EventUpdate(c *gin.Context) error {
	form := new(eventForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	event, err := GetEvent(c)

	if err != nil {
		return err
	}

	projectID := model.GetProjectIDForElement(event.ElementID)

	if err := CheckProjectPermission(c, projectID, true); err != nil {
		return err
	}

	if err := saveEvent(form, event); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, event)
}

func EventDestroy(c *gin.Context) error {
	event, err := GetEvent(c)

	if err != nil {
		return err
	}

	projectID := model.GetProjectIDForElement(event.ElementID)

	if err := CheckProjectPermission(c, projectID, true); err != nil {
		return err
	}

	if err := event.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

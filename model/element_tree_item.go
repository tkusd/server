package model

import (
	"encoding/json"
	"errors"

	"github.com/tkusd/server/model/types"
)

type ElementTreeItem struct {
	ID       types.UUID        `json:"id"`
	Elements []ElementTreeItem `json:"elements,omitempty"`
}

func (e *ElementTreeItem) UnmarshalJSON(data []byte) error {
	var val interface{}

	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch v := val.(type) {
	case string:
		e.ID = types.ParseUUID(v)

		if !e.ID.Valid() {
			return errors.New("Element ID is not a valid UUID.")
		}

		break

	case map[string]interface{}:
		id, ok := v["id"]

		if !ok {
			return errors.New("Element ID is required.")
		}

		idStr, ok := id.(string)

		if !ok {
			return errors.New("Element ID must be a string.")
		}

		e.ID = types.ParseUUID(idStr)

		if !e.ID.Valid() {
			return errors.New("Element ID is not a valid UUID.")
		}

		if arr, ok := v["elements"]; ok {
			elements, ok := arr.([]interface{})

			if !ok {
				return errors.New("Elements must be an array.")
			}

			var list []ElementTreeItem

			for _, s := range elements {
				b, ok := s.([]byte)

				if !ok {
					return errors.New("Element item is not valid.")
				}

				var item ElementTreeItem

				if err := json.Unmarshal(b, &item); err != nil {
					return err
				}

				list = append(list, item)
			}

			e.Elements = list
		}

		break

	default:
		return errors.New("The item should be either a UUID or an object including UUID and element list.")
	}

	return nil
}

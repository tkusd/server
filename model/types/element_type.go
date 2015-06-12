package types

import "encoding/json"

// ElementType is a enum for element type.
type ElementType int16

// Element types
const (
	ElementTypeScreen ElementType = iota + 1
	ElementTypeText
	ElementTypeLayout
	ElementTypeButton
	ElementTypeInput
	ElementTypeLink
	ElementTypeImage
	ElementTypeList
)

var elementTypeMap = map[ElementType]string{
	ElementTypeScreen: "screen",
	ElementTypeText:   "text",
	ElementTypeLayout: "layout",
	ElementTypeButton: "button",
	ElementTypeInput:  "input",
	ElementTypeLink:   "link",
	ElementTypeImage:  "image",
	ElementTypeList:   "list",
}

var elementTypeReversedMap map[string]ElementType

func init() {
	elementTypeReversedMap = make(map[string]ElementType)

	for key, value := range elementTypeMap {
		elementTypeReversedMap[value] = key
	}
}

func (t ElementType) String() string {
	if value, ok := elementTypeMap[t]; ok {
		return value
	}

	return ""
}

func findElementTypeFromText(str string) ElementType {
	if value, ok := elementTypeReversedMap[str]; ok {
		return value
	}

	return 0
}

// MarshalJSON implements json.Marshaler interface.
func (t ElementType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

// MarshalText implements encoding.TextMarshaler interface.
func (t ElementType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *ElementType) UnmarshalJSON(data []byte) error {
	var str string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*t = findElementTypeFromText(str)
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface.
func (t *ElementType) UnmarshalText(data []byte) error {
	return t.UnmarshalJSON(data)
}

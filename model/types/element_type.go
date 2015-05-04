package types

import "encoding/json"

// ElementType is a enum for element type.
type ElementType int16

// Element types
const (
	ElementTypeScreen ElementType = 1
	ElementTypeText   ElementType = 2
)

// Element types in string
const (
	elementTypeStringScreen = "screen"
	elementTypeStringText   = "text"
)

func (t ElementType) String() string {
	switch t {
	case ElementTypeScreen:
		return elementTypeStringScreen

	case ElementTypeText:
		return elementTypeStringText
	}
	return ""
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

func findElementTypeFromText(str string) ElementType {
	switch str {
	case elementTypeStringScreen:
		return ElementTypeScreen

	case elementTypeStringText:
		return ElementTypeText
	}

	return 0
}

package util

import "net/http"

const (
	// 1000: Unknown error
	UnknownError = 1000
	ServerError  = 1001
	// 1100: Validation error
	RequiredError        = 1100
	ContentTypeError     = 1101
	DeserializationError = 1102
	TypeError            = 1103
	EmailError           = 1104
	LengthError          = 1105
	URLError             = 1106
	// 1200: Resource error
	UserNotFoundError  = 1200
	TokenNotFoundError = 1201
	// 1300: Data error
	WrongPasswordError = 1300
	EmailUsedError     = 1301
)

type APIError struct {
	Code    int    `json:"error"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"-"`
	Field   string `json:"field,omitempty"`
}

func (err *APIError) Error() string {
	return err.Message
}

func HandleAPIError(res http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *APIError:
		if e.Status == 0 {
			e.Status = http.StatusBadRequest
		}

		RenderJSON(res, e.Status, e)
		break

	default:
		RenderJSON(res, http.StatusInternalServerError, &APIError{
			Code:    ServerError,
			Message: "Server error",
		})
	}
}

package util

// 1000: Unknown error
const (
	UnknownError  = 1000
	ServerError   = 1001
	NotFoundError = 1002
)

// 1100: Validation error
const (
	RequiredError        = 1100
	ContentTypeError     = 1101
	DeserializationError = 1102
	TypeError            = 1103
	EmailError           = 1104
	LengthError          = 1105
	URLError             = 1106
)

// 1200: Resource error
const (
	UserNotFoundError    = 1200
	TokenNotFoundError   = 1201
	ProjectNotFoundError = 1202
)

// 1300: Data error
const (
	WrongPasswordError = 1300
	EmailUsedError     = 1301
	TokenRequiredError = 1302
	TokenInvalidError  = 1303
	UserForbiddenError = 1304
)

// APIError represents an API error.
type APIError struct {
	Code    int    `json:"error"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"-"`
	Field   string `json:"field,omitempty"`
}

func (err APIError) Error() string {
	return err.Message
}

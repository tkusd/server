package common

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/tkusd/server/util"
)

// BindForm binds data to the struct.
func BindForm(res http.ResponseWriter, req *http.Request, mapper binding.FieldMapper) error {
	if err := binding.Bind(req, mapper); err != nil {
		e := err[0]
		code := util.UnknownError
		status := http.StatusBadRequest
		field := ""

		switch e.Classification {
		case binding.ContentTypeError:
			code = util.ContentTypeError
			status = http.StatusUnsupportedMediaType
			break

		case binding.DeserializationError:
			code = util.DeserializationError
			break

		case binding.RequiredError:
			code = util.RequiredError
			break

		case binding.TypeError:
			code = util.TypeError
			break
		}

		if len(e.Fields()) > 0 {
			field = e.Fields()[0]
		}

		return &util.APIError{
			Code:    code,
			Message: e.Message,
			Status:  status,
			Field:   field,
		}
	}

	return nil
}

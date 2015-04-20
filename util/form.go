package util

import (
	"net/http"

	"github.com/mholt/binding"
)

func BindForm(res http.ResponseWriter, req *http.Request, mapper binding.FieldMapper) bool {
	if err := binding.Bind(req, mapper); err != nil {
		firstErr := err[0]
		code := UnknownError

		switch firstErr.Classification {
		case binding.ContentTypeError:
			code = ContentTypeError
			break

		case binding.DeserializationError:
			code = DeserializationError
			break

		case binding.RequiredError:
			code = RequiredError
			break

		case binding.TypeError:
			code = TypeError
			break
		}

		RenderJSON(res, http.StatusBadRequest, APIError{
			Code:    code,
			Message: firstErr.Message,
		})

		return true
	}

	return false
}

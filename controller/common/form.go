package common

import (
	"net/http"

	"github.com/mholt/binding"
	"github.com/tommy351/app-studio-server/util"
)

// BindForm binds data to the struct.
func BindForm(res http.ResponseWriter, req *http.Request, mapper binding.FieldMapper) bool {
	if err := binding.Bind(req, mapper); err != nil {
		e := err[0]
		code := util.UnknownError

		switch e.Classification {
		case binding.ContentTypeError:
			code = util.ContentTypeError
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

		HandleAPIError(res, &util.APIError{
			Code:    code,
			Message: e.Message,
		})

		return true
	}

	return false
}

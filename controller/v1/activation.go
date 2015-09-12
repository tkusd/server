package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/util"
)

func ActivateUser(c *gin.Context) error {
	activationToken := c.Param(activationIDParam)
	user, err := model.GetUserByActivationToken(activationToken)

	if err != nil {
		return &util.APIError{
			Code:    util.UserActivationTokenMismatchError,
			Message: "Activation token mismatch",
		}
	}

	if user.IsActivated {
		return &util.APIError{
			Code:    util.UserAlreadyActivatedError,
			Message: "User has already been activated",
		}
	}

	user.SetActivated(true)

	if err := user.Save(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

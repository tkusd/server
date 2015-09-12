package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

type passwordResetCreateForm struct {
	Email string `json:"email"`
}

func (form *passwordResetCreateForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Email: "email",
	}
}

func PasswordResetCreate(c *gin.Context) error {
	form := new(passwordResetCreateForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	if form.Email == "" {
		return &util.APIError{
			Field:   "email",
			Code:    util.RequiredError,
			Message: "Email is required.",
		}
	}

	user, err := model.GetUserByEmail(form.Email)

	if err != nil {
		return &util.APIError{
			Field:   "email",
			Code:    util.UserNotFoundError,
			Message: "User not found",
		}
	}

	user.PasswordResetToken = types.NewRandomUUID()
	user.PasswordResetAt = types.Now()

	if err := user.Save(); err != nil {
		return err
	}

	msg := util.Mailgun.NewMessage(
		"Diff <noreply@tkusd.zespia.tw>",
		"Reset your password",
		"Click this link to reset your password: http://tkusd.zespia.tw/reset_password/"+user.PasswordResetToken.String(),
		user.Email,
	)

	go util.Mailgun.Send(msg)

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

type passwordResetUpdateForm struct {
	Password string `json:"password"`
}

func (form *passwordResetUpdateForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Password: "password",
	}
}

func PasswordResetUpdate(c *gin.Context) error {
	form := new(passwordResetUpdateForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	if form.Password == "" {
		return &util.APIError{
			Field:   "password",
			Code:    util.RequiredError,
			Message: "Password is required.",
		}
	}

	resetToken := c.Param(passwordResetIDParam)
	user, err := model.GetUserByPasswordResetToken(resetToken)

	if err != nil {
		return &util.APIError{
			Code:    util.PasswordResetTokenMismatchError,
			Message: "Password reset token mismatch",
		}
	}

	if user.PasswordResetAt.Time.Add(time.Hour * 6).Before(time.Now()) {
		return &util.APIError{
			Code:    util.PasswordResetTokenExpiredError,
			Message: "Password reset token was expired",
		}
	}

	if err := user.GeneratePassword(form.Password); err != nil {
		return err
	}

	user.PasswordResetToken = types.UUID{}

	if err := user.Save(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

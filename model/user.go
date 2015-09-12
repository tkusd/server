package model

import (
	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
	"github.com/tkusd/server/config"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

// User represents the data structure of a user.
type User struct {
	ID                 types.UUID `json:"id"`
	Name               string     `json:"name"`
	Password           []byte     `json:"-"`
	Email              string     `json:"email"`
	Avatar             string     `json:"avatar"`
	CreatedAt          types.Time `json:"created_at"`
	UpdatedAt          types.Time `json:"updated_at"`
	IsActivated        bool       `json:"is_activated"`
	ActivationToken    types.UUID `json:"-"`
	Language           string     `json:"language"`
	PasswordResetToken types.UUID `json:"-"`
	PasswordResetAt    types.Time `json:"-"`
}

// PublicProfile returns the data for public display.
func (u *User) PublicProfile() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"avatar":     u.Avatar,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
		"language":   u.Language,
	}
}

// Save creates or updates data in the database.
func (u *User) Save() error {
	u.Name = govalidator.Trim(u.Name, "")
	u.Email = govalidator.Trim(u.Email, "")
	u.Avatar = govalidator.Trim(u.Avatar, "")
	u.Language = govalidator.Trim(u.Language, "")

	if u.Name == "" {
		return &util.APIError{
			Field:   "name",
			Code:    util.RequiredError,
			Message: "Name is required.",
		}
	}

	if len(u.Name) > 100 {
		return &util.APIError{
			Field:   "name",
			Code:    util.LengthError,
			Message: "Maximum length of name is 100.",
		}
	}

	if !govalidator.IsEmail(u.Email) {
		return &util.APIError{
			Field:   "email",
			Code:    util.EmailError,
			Message: "Email is invalid.",
		}
	}

	if u.Avatar != "" && !govalidator.IsURL(u.Avatar) {
		return &util.APIError{
			Field:   "avatar",
			Code:    util.URLError,
			Message: "Avatar URL is invalid.",
		}
	}

	if u.Avatar == "" {
		u.Avatar = util.Gravatar(u.Email)
	}

	if u.Language == "" {
		u.Language = "en"
	}

	if len(u.Language) > 35 {
		return &util.APIError{
			Field:   "language",
			Code:    util.LengthError,
			Message: "Maximum length of language is 35.",
		}
	}

	if err := db.Save(u).Error; err != nil {
		switch e := err.(type) {
		case *pq.Error:
			switch e.Code.Name() {
			case UniqueViolation:
				return &util.APIError{
					Code:    util.EmailUsedError,
					Message: "Email has been used.",
					Field:   "email",
				}
			}
		}

		return err
	}

	return nil
}

func (u *User) AfterCreate() error {
	if !u.IsActivated && config.Config.EmailActivation {
		msg := util.Mailgun.NewMessage(
			"Diff <noreply@tkusd.zespia.tw>",
			"Activate your account",
			"Click this link to activate your account: http://tkusd.zespia.tw/activation/"+u.ActivationToken.String(),
			u.Email,
		)

		go util.Mailgun.Send(msg)
	}

	return nil
}

// Delete deletes data from the database.
func (u *User) Delete() error {
	return db.Delete(u).Error
}

// Exists returns true if the record exists.
func (u *User) Exists() bool {
	return exists("users", u.ID.String())
}

// SetActivated updates the activation status of a user.
func (u *User) SetActivated(activated bool) {
	if activated {
		u.IsActivated = true
		u.ActivationToken = types.UUID{}
	} else {
		u.IsActivated = false
		u.ActivationToken = types.NewRandomUUID()
	}
}

func validatePassword(password string) error {
	if password == "" {
		return &util.APIError{
			Code:    util.RequiredError,
			Field:   "password",
			Message: "Password is required.",
		}
	}

	if len(password) < 6 || len(password) > 50 {
		return &util.APIError{
			Code:    util.LengthError,
			Field:   "password",
			Message: "The length of password must be between 6 to 50.",
		}
	}

	return nil
}

// GeneratePassword generates the bcrypt password for a user.
func (u *User) GeneratePassword(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	hash, err := util.GenerateBcryptHash(password)

	if err != nil {
		return err
	}

	u.Password = hash

	return nil
}

// Authenticate authenticates a user.
func (u *User) Authenticate(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	if err := util.CompareBcryptHash(u.Password, password); err != nil {
		return &util.APIError{
			Field:   "password",
			Code:    util.WrongPasswordError,
			Message: "Password is wrong.",
		}
	}

	return nil
}

// GetUser returns the user data.
func GetUser(id types.UUID) (*User, error) {
	user := new(User)

	if err := db.Where("id = ?", id.String()).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail searchs user data by email.
func GetUserByEmail(email string) (*User, error) {
	user := new(User)

	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByActivationToken(str string) (*User, error) {
	user := new(User)

	if err := db.Where("activation_token = ?", str).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByPasswordResetToken(str string) (*User, error) {
	user := new(User)

	if err := db.Where("password_reset_token = ?", str).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

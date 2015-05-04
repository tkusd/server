package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"

	"code.google.com/p/go-uuid/uuid"
)

// User represents the data structure of a user.
type User struct {
	ID              types.UUID       `json:"id"`
	Name            string           `json:"name"`
	Password        []byte           `json:"-"`
	Email           string           `json:"email"`
	Avatar          string           `json:"avatar"`
	CreatedAt       types.Time       `json:"created_at"`
	UpdatedAt       types.Time       `json:"updated_at"`
	IsActivated     bool             `json:"is_activated"`
	ActivationToken types.Base64Hash `json:"-"`
}

// PublicProfile returns the data for public display.
func (u *User) PublicProfile() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"avatar":     u.Avatar,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}

// BeforeSave is called when the data is about to be saved.
func (u *User) BeforeSave() error {
	u.UpdatedAt = types.Now()
	return nil
}

// BeforeCreate is called when the data is about to be created.
func (u *User) BeforeCreate() error {
	u.CreatedAt = types.Now()
	return nil
}

// Save creates or updates data in the database.
func (u *User) Save() error {
	u.Name = govalidator.Trim(u.Name, "")
	u.Email = govalidator.Trim(u.Email, "")
	u.Avatar = govalidator.Trim(u.Avatar, "")

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
	} else {
		u.IsActivated = false
		u.ActivationToken = types.Base64Hash{types.SHA256(u.Email, time.Now().String(), uuid.New())}
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

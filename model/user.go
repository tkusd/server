package model

import (
	"encoding/json"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/lib/pq"
	"github.com/tommy351/app-studio-server/util"

	"net/http"

	"code.google.com/p/go-uuid/uuid"
)

type User struct {
	ID              util.UUID
	Name            string
	Password        []byte
	Email           string
	Avatar          string
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IsActivated     bool      `json:"is_activated"`
	ActivationToken util.Hash
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":           u.ID,
		"name":         u.Name,
		"email":        u.Email,
		"avatar":       u.Avatar,
		"created_at":   util.ISOTime(u.CreatedAt),
		"updated_at":   util.ISOTime(u.UpdatedAt),
		"is_activated": u.IsActivated,
	})
}

func (u *User) PublicProfile() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"avatar":     u.Avatar,
		"created_at": util.ISOTime(u.CreatedAt),
		"updated_at": util.ISOTime(u.UpdatedAt),
	}
}

func (u *User) BeforeSave() error {
	u.Name = govalidator.Trim(u.Name, "")
	u.Email = govalidator.Trim(u.Email, "")
	u.Avatar = govalidator.Trim(u.Avatar, "")
	u.UpdatedAt = time.Now().UTC()

	return nil
}

func (u *User) BeforeCreate() error {
	u.CreatedAt = time.Now().UTC()
	return nil
}

func (u *User) SetActivated(activated bool) {
	if activated {
		u.IsActivated = true
	} else {
		u.IsActivated = false
		u.ActivationToken = util.SHA256(u.Email, time.Now().String(), uuid.New())
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

func (u *User) GeneratePassword(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	if hash, err := util.GenerateBcryptHash(password); err != nil {
		return err
	} else {
		u.Password = hash
	}

	return nil
}

func (u *User) Authenticate(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}

	if err := util.CompareBcryptHash(u.Password, password); err != nil {
		return &util.APIError{
			Field:   "password",
			Code:    util.WrongPasswordError,
			Message: "Password is wrong",
			Status:  http.StatusUnauthorized,
		}
	}

	return nil
}

func handleUserDBError(err error) error {
	switch e := err.(type) {
	case *pq.Error:
		switch e.Code.Name() {
		case UniqueViolation:
			return &util.APIError{
				Code:    util.EmailUsedError,
				Status:  http.StatusBadRequest,
				Message: "Email has been used.",
				Field:   "email",
			}
		}

		return e
	}

	return err
}

func validateUser(u *User) error {
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
	} else {
		u.Avatar = util.Gravatar(u.Email)
	}

	return nil
}

func GetUser(id util.UUID) (*User, error) {
	user := new(User)

	if err := db.Where("id = ?", id.String()).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUser(user *User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	if err := db.Create(user).Error; err != nil {
		return handleUserDBError(err)
	}

	return nil
}

func UpdateUser(user *User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	if err := db.Save(user).Error; err != nil {
		return handleUserDBError(err)
	}

	return nil
}

func DeleteUser(user *User) error {
	return db.Delete(user).Error
}

func GetUserByEmail(email string) (*User, error) {
	user := new(User)

	if err := db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

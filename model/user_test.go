package model

import (
	"net/http"
	"testing"

	"time"

	"log"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/util"
)

func createTestUser() (*User, error) {
	user := &User{
		Name:  "abc",
		Email: "abc@example.com",
	}

	if err := user.GeneratePassword("123456"); err != nil {
		return nil, err
	}

	user.SetActivated(false)

	if err := CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func TestUser(t *testing.T) {
	Convey("BeforeSave", t, func() {
		Convey("Trim name", func() {
			user := &User{Name: "  abc  "}
			user.BeforeSave()
			So(user.Name, ShouldEqual, "abc")
		})

		Convey("Trim email", func() {
			user := &User{Email: "  abc@example.com   "}
			user.BeforeSave()
			So(user.Email, ShouldEqual, "abc@example.com")
		})

		Convey("Trim avatar", func() {
			user := &User{Avatar: "   http://example.com/test.jpg  "}
			user.BeforeSave()
			So(user.Avatar, ShouldEqual, "http://example.com/test.jpg")
		})

		Convey("Update modified time", func() {
			now := time.Now()
			user := &User{}
			user.BeforeSave()
			So(user.UpdatedAt, ShouldHappenOnOrAfter, now)
		})
	})

	Convey("BeforeCreate", t, func() {
		Convey("Give created time", func() {
			now := time.Now()
			user := &User{}
			user.BeforeCreate()
			So(user.CreatedAt, ShouldHappenOnOrAfter, now)
		})
	})

	Convey("SetActivated", t, func() {
		Convey("true", func() {
			user := &User{}
			user.SetActivated(true)
			So(user.IsActivated, ShouldBeTrue)
		})

		Convey("false", func() {
			user := &User{}
			user.SetActivated(false)
			So(user.IsActivated, ShouldBeFalse)
			So(user.ActivationToken, ShouldNotBeEmpty)
		})
	})

	Convey("GeneratePassword", t, func() {
		user := &User{}
		user.GeneratePassword("123456")
		So(user.Password, ShouldNotBeEmpty)
	})

	Convey("Authenticate", t, func() {
		user := &User{}
		password := "123456"
		user.GeneratePassword(password)

		Convey("Success", func() {
			err := user.Authenticate(password)
			So(err, ShouldBeNil)
		})

		Convey("Failed", func() {
			err := user.Authenticate("abcdef")
			So(err, ShouldResemble, &util.APIError{
				Field:   "password",
				Code:    util.WrongPasswordError,
				Message: "Password is wrong",
				Status:  http.StatusUnauthorized,
			})
		})
	})
}

func TestValidatePassword(t *testing.T) {
	Convey("Success", t, func() {
		err := validatePassword("123456")
		So(err, ShouldBeNil)
	})

	Convey("Password is required", t, func() {
		err := validatePassword("")

		So(err, ShouldResemble, &util.APIError{
			Code:    util.RequiredError,
			Field:   "password",
			Message: "Password is required.",
		})
	})

	Convey("Minimum length of password is 6", t, func() {
		err := validatePassword("12345")

		So(err, ShouldResemble, &util.APIError{
			Code:    util.LengthError,
			Field:   "password",
			Message: "The length of password must be between 6 to 50.",
		})
	})

	Convey("Maximum length of password is 50", t, func() {
		err := validatePassword("12346545341354354135431354354dgfmgdjfgoierterterttt")

		So(err, ShouldResemble, &util.APIError{
			Code:    util.LengthError,
			Field:   "password",
			Message: "The length of password must be between 6 to 50.",
		})
	})
}

func TestValidateUser(t *testing.T) {
	Convey("Name is required", t, func() {
		user := &User{}
		err := validateUser(user)

		So(err, ShouldResemble, &util.APIError{
			Field:   "name",
			Code:    util.RequiredError,
			Message: "Name is required.",
		})
	})

	Convey("Maximum length of name is 100", t, func() {
		user := &User{Name: "adijosjfosejfoeijfroiejrowijerowjerowiejrowiejrowiejrowiejroweirjoweijrwoeirjwerwrjwioerjwoerwerrrrrr"}
		err := validateUser(user)

		So(err, ShouldResemble, &util.APIError{
			Field:   "name",
			Code:    util.LengthError,
			Message: "Maximum length of name is 100.",
		})
	})

	Convey("Email is invalid", t, func() {
		user := &User{
			Name:  "abc",
			Email: "abc@",
		}
		err := validateUser(user)

		So(err, ShouldResemble, &util.APIError{
			Field:   "email",
			Code:    util.EmailError,
			Message: "Email is invalid.",
		})
	})

	Convey("Avatar URL is invalid", t, func() {
		user := &User{
			Name:   "abc",
			Email:  "abc@example.com",
			Avatar: "ht://erw",
		}
		err := validateUser(user)

		So(err, ShouldResemble, &util.APIError{
			Field:   "avatar",
			Code:    util.URLError,
			Message: "Avatar URL is invalid.",
		})
	})

	Convey("Define avatar from email", t, func() {
		user := &User{
			Name:  "abc",
			Email: "abc@example.com",
		}
		validateUser(user)
		So(user.Avatar, ShouldEqual, util.Gravatar(user.Email))
	})
}

func TestGetUser(t *testing.T) {
	user, err := createTestUser()
	defer DeleteUser(user)

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		result, err := GetUser(user.ID)

		if err != nil {
			log.Fatal(err)
		}

		So(result.ID.String(), ShouldEqual, user.ID.String())
		So(result.Name, ShouldEqual, user.Name)
		So(result.Email, ShouldEqual, user.Email)
		So(result.Avatar, ShouldEqual, user.Avatar)
		So(result.Password, ShouldResemble, user.Password)
		So(util.ISOTime(result.CreatedAt), ShouldResemble, util.ISOTime(user.CreatedAt))
		So(util.ISOTime(result.UpdatedAt), ShouldResemble, util.ISOTime(user.UpdatedAt))
		So(user.IsActivated, ShouldEqual, user.IsActivated)
		So(user.ActivationToken, ShouldResemble, user.ActivationToken)
	})

	Convey("Failed", t, func() {
		_, err := GetUser(util.NewRandomUUID())
		So(err, ShouldNotBeNil)
	})
}

func TestCreateUser(t *testing.T) {
	Convey("Success", t, func() {
		user, err := createTestUser()
		defer DeleteUser(user)

		if err != nil {
			log.Fatal(err)
		}

		So(user, ShouldNotBeNil)
	})

	Convey("Email is used", t, func() {
		user, err := createTestUser()
		defer DeleteUser(user)

		if err != nil {
			log.Fatal(err)
		}

		_, err = createTestUser()

		So(err, ShouldResemble, &util.APIError{
			Code:    util.EmailUsedError,
			Status:  http.StatusBadRequest,
			Message: "Email has been used.",
			Field:   "email",
		})
	})
}

func TestUpdateUser(t *testing.T) {
	user, err := createTestUser()
	defer DeleteUser(user)

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		user.Name = "new name"
		UpdateUser(user)

		newUser, err := GetUser(user.ID)

		if err != nil {
			log.Fatal(err)
		}

		So(newUser.Name, ShouldEqual, user.Name)
	})

	SkipConvey("Email is used", t, nil)
}

func TestDeleteUser(t *testing.T) {
	user, err := createTestUser()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		DeleteUser(user)
		user, _ := GetUser(user.ID)
		So(user, ShouldBeNil)
	})
}

func TestGetUserByEmail(t *testing.T) {
	user, err := createTestUser()
	defer DeleteUser(user)

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		u, err := GetUserByEmail(user.Email)

		if err != nil {
			log.Fatal(err)
		}

		So(u.ID, ShouldResemble, user.ID)
	})

	Convey("Failed", t, func() {
		u, err := GetUserByEmail("nothing@nothing.com")

		So(u, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

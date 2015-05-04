package model

import (
	"strings"
	"testing"

	"time"

	"log"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/model/types"
	"github.com/tommy351/app-studio-server/util"
)

func createTestUser(data fixtureUser) (*User, error) {
	user := &User{
		Name:  data.Name,
		Email: data.Email,
	}

	if err := user.GeneratePassword(data.Password); err != nil {
		return nil, err
	}

	user.SetActivated(false)

	if err := user.Save(); err != nil {
		return nil, err
	}

	return user, nil
}

func TestUser(t *testing.T) {
	Convey("BeforeSave", t, func() {
		now := time.Now()
		user := &User{}
		user.BeforeSave()
		So(user.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("BeforeCreate", t, func() {
		now := time.Now()
		user := &User{}
		user.BeforeCreate()
		So(user.CreatedAt.Time, ShouldHappenOnOrAfter, now)
	})

	Convey("Save", t, func() {
		Convey("Trim name", func() {
			user := &User{Name: "   abc   "}
			user.Save()
			defer user.Delete()
			So(user.Name, ShouldEqual, "abc")
		})

		Convey("Trim email", func() {
			user := &User{Email: "   abc@we.com  "}
			user.Save()
			defer user.Delete()
			So(user.Email, ShouldEqual, "abc@we.com")
		})

		Convey("Trim avatar", func() {
			user := &User{Avatar: "   http://example.com/test.jpg   "}
			user.Save()
			defer user.Delete()
			So(user.Avatar, ShouldEqual, "http://example.com/test.jpg")
		})

		Convey("Name is required", func() {
			user := &User{}
			err := user.Save()
			So(err, ShouldResemble, &util.APIError{
				Field:   "name",
				Code:    util.RequiredError,
				Message: "Name is required.",
			})
		})

		Convey("Maximum length of name is 100", func() {
			user := &User{Name: strings.Repeat("a", 101)}
			err := user.Save()

			So(err, ShouldResemble, &util.APIError{
				Field:   "name",
				Code:    util.LengthError,
				Message: "Maximum length of name is 100.",
			})
		})

		Convey("Email is invalid", func() {
			user := &User{
				Name:  "abc",
				Email: "abc@",
			}
			err := user.Save()

			So(err, ShouldResemble, &util.APIError{
				Field:   "email",
				Code:    util.EmailError,
				Message: "Email is invalid.",
			})
		})

		Convey("Avatar URL is invalid", func() {
			user := &User{
				Name:   "abc",
				Email:  "abc@example.com",
				Avatar: "ht://erw",
			}
			err := user.Save()

			So(err, ShouldResemble, &util.APIError{
				Field:   "avatar",
				Code:    util.URLError,
				Message: "Avatar URL is invalid.",
			})
		})

		Convey("Create", func() {
			now := time.Now().Truncate(time.Second)
			user, err := createTestUser(fixtureUsers[0])
			defer user.Delete()

			if err != nil {
				log.Fatal(err)
			}

			u, _ := GetUser(user.ID)
			So(u.ID, ShouldResemble, user.ID)
			So(u.Name, ShouldEqual, user.Name)
			So(u.Email, ShouldEqual, user.Email)
			So(u.Avatar, ShouldEqual, util.Gravatar(u.Email))
			So(u.CreatedAt.Time, ShouldHappenOnOrAfter, now)
			So(u.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
			So(u.IsActivated, ShouldBeFalse)
			So(u.ActivationToken, ShouldNotBeEmpty)
		})

		Convey("Update", func() {
			user, err := createTestUser(fixtureUsers[0])
			defer user.Delete()

			if err != nil {
				log.Fatal(err)
			}

			newName := "new name"
			now := time.Now().Truncate(time.Second)
			user.Name = newName

			if err := user.Save(); err != nil {
				log.Fatal(err)
			}

			u, _ := GetUser(user.ID)
			So(u.Name, ShouldEqual, newName)
			So(u.UpdatedAt.Time, ShouldHappenOnOrAfter, now)
		})

		Convey("Email is used", func() {
			user, err := createTestUser(fixtureUsers[0])
			defer user.Delete()

			if err != nil {
				log.Fatal(err)
			}

			_, err = createTestUser(fixtureUsers[0])
			So(err, ShouldResemble, &util.APIError{
				Code:    util.EmailUsedError,
				Message: "Email has been used.",
				Field:   "email",
			})
		})
	})

	Convey("Delete", t, func() {
		user, err := createTestUser(fixtureUsers[0])

		if err != nil {
			log.Fatal(err)
		}

		user.Delete()
		u, _ := GetUser(user.ID)
		So(u, ShouldBeNil)
	})

	Convey("Exists", t, func() {
		Convey("true", func() {
			user, err := createTestUser(fixtureUsers[0])
			defer user.Delete()

			if err != nil {
				log.Fatal(err)
			}

			So(user.Exists(), ShouldBeTrue)
		})

		Convey("false", func() {
			user := &User{ID: types.NewRandomUUID()}
			So(user.Exists(), ShouldBeFalse)
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

	var validatePasswordTests = []struct {
		name     string
		password string
		err      error
	}{
		{"Password is required", "", &util.APIError{
			Code:    util.RequiredError,
			Field:   "password",
			Message: "Password is required.",
		}},
		{"Minimum length of password is 6", strings.Repeat("a", 5), &util.APIError{
			Code:    util.LengthError,
			Field:   "password",
			Message: "The length of password must be between 6 to 50.",
		}},
		{"Maximum length of password is 50", strings.Repeat("a", 51), &util.APIError{
			Code:    util.LengthError,
			Field:   "password",
			Message: "The length of password must be between 6 to 50.",
		}},
	}

	Convey("GeneratePassword", t, func() {
		Convey("Success", func() {
			user := &User{}
			user.GeneratePassword("123456")
			So(user.Password, ShouldNotBeEmpty)
		})

		for _, test := range validatePasswordTests {
			Convey(test.name, func() {
				user := &User{}
				err := user.GeneratePassword(test.password)
				So(err, ShouldResemble, test.err)
			})
		}
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
				Message: "Password is wrong.",
			})
		})

		for _, test := range validatePasswordTests {
			Convey(test.name, func() {
				user := &User{}
				err := user.Authenticate(test.password)
				So(err, ShouldResemble, test.err)
			})
		}
	})
}

func TestGetUser(t *testing.T) {
	user, err := createTestUser(fixtureUsers[0])
	defer user.Delete()

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
		So(result.CreatedAt.ISOTime(), ShouldResemble, user.CreatedAt.ISOTime())
		So(result.UpdatedAt.ISOTime(), ShouldResemble, user.UpdatedAt.ISOTime())
		So(user.IsActivated, ShouldEqual, user.IsActivated)
		So(user.ActivationToken, ShouldResemble, user.ActivationToken)
	})

	Convey("Failed", t, func() {
		_, err := GetUser(types.NewRandomUUID())
		So(err, ShouldNotBeNil)
	})
}

func TestGetUserByEmail(t *testing.T) {
	user, err := createTestUser(fixtureUsers[0])
	defer user.Delete()

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

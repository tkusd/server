package model

import (
	"testing"

	"log"

	. "github.com/smartystreets/goconvey/convey"
)

func createTestToken(u *User) (*Token, error) {
	token := &Token{UserID: u.ID}

	if err := token.Save(); err != nil {
		return nil, err
	}

	return token, nil
}

func TestToken(t *testing.T) {
	user, err := createTestUser(fixtureUsers[0])
	defer user.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Save", t, func() {
		token, err := createTestToken(user)
		defer token.Delete()

		if err != nil {
			log.Fatal(err)
		}

		So(token.ID, ShouldNotBeEmpty)
		So(token.UserID, ShouldResemble, user.ID)
	})

	Convey("Delete", t, func() {
		token, err := createTestToken(user)

		if err != nil {
			log.Fatal(err)
		}

		token.Delete()

		t, _ := GetToken(token.ID)
		So(t, ShouldBeNil)
	})
}

func TestGetToken(t *testing.T) {
	user, err := createTestUser(fixtureUsers[0])
	defer user.Delete()

	if err != nil {
		log.Fatal(err)
	}

	token, err := createTestToken(user)
	defer token.Delete()

	if err != nil {
		log.Fatal(err)
	}

	Convey("Hash ID", t, func() {
		t, err := GetToken(token.ID)

		if err != nil {
			log.Fatal(err)
		}

		So(t.ID, ShouldResemble, token.ID)
	})

	Convey("Base64 ID", t, func() {
		t, err := GetTokenBase64(token.ID.String())

		if err != nil {
			log.Fatal(err)
		}

		So(t.ID, ShouldResemble, token.ID)
	})
}

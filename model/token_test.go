package model

import (
	"testing"

	"time"

	"encoding/base64"

	"log"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/tommy351/app-studio-server/util"
)

func createTestToken(u *User) (*Token, error) {
	token := &Token{UserID: u.ID}

	if err := CreateToken(token); err != nil {
		return nil, err
	} else {
		return token, nil
	}
}

func TestToken(t *testing.T) {
	Convey("BeforeCreate", t, func() {
		Convey("Give created time", func() {
			now := time.Now()
			token := &Token{}
			token.BeforeCreate()
			So(token.CreatedAt, ShouldHappenOnOrAfter, now)
		})

		Convey("Give ID", func() {
			token := &Token{}
			token.BeforeCreate()
			So(token.ID, ShouldNotBeEmpty)
		})
	})

	Convey("GetBase64Key", t, func() {
		token := &Token{
			ID: util.SHA256("test"),
		}

		So(token.GetBase64ID(), ShouldResemble, base64.URLEncoding.EncodeToString(token.ID))
	})
}

func TestCreateToken(t *testing.T) {
	user, err := createTestUser()
	defer DeleteUser(user)

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		token, err := createTestToken(user)
		defer DeleteToken(token)

		if err != nil {
			log.Fatal(err)
		}

		So(token.UserID, ShouldResemble, user.ID)
	})
}

func TestDeleteToken(t *testing.T) {
	user, err := createTestUser()
	defer DeleteUser(user)

	if err != nil {
		log.Fatal(err)
	}

	Convey("Success", t, func() {
		token, err := createTestToken(user)

		if err != nil {
			log.Fatal(err)
		}

		DeleteToken(token)

		_, err = GetToken(token.ID)
		So(err, ShouldNotBeNil)
	})
}

func TestGetToken(t *testing.T) {
	user, err := createTestUser()
	defer DeleteUser(user)

	if err != nil {
		log.Fatal(err)
	}

	token, err2 := createTestToken(user)
	defer DeleteToken(token)

	if err2 != nil {
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
		t, err := GetTokenBase64(token.GetBase64ID())

		if err != nil {
			log.Fatal(err)
		}

		So(t.ID, ShouldResemble, token.ID)
	})
}

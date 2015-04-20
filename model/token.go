package model

import (
	"time"

	"encoding/json"

	"encoding/base64"

	"code.google.com/p/go-uuid/uuid"
	"github.com/tommy351/app-studio-server/util"
)

type Token struct {
	ID        util.Hash
	UserID    util.UUID
	CreatedAt time.Time `json:"created_at"`
}

func (t Token) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         t.GetBase64ID(),
		"user_id":    t.UserID,
		"created_at": util.ISOTime(t.CreatedAt),
	})
}

func (t *Token) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if id, err := base64.URLEncoding.DecodeString(obj["id"].(string)); err != nil {
		return err
	} else {
		t.ID = id
	}

	t.UserID = util.ParseUUID(obj["user_id"].(string))

	if date, err := util.ParseISOTime(obj["created_at"].(string)); err != nil {
		return err
	} else {
		t.CreatedAt = date
	}

	return nil
}

func (t *Token) BeforeCreate() error {
	t.CreatedAt = time.Now().UTC()
	t.ID = util.SHA256(t.UserID.String(), t.CreatedAt.String(), uuid.New())

	return nil
}

func (t *Token) GetBase64ID() string {
	return base64.URLEncoding.EncodeToString(t.ID)
}

func CreateToken(token *Token) error {
	if err := db.Create(token).Error; err != nil {
		return err
	}

	return nil
}

func DeleteToken(token *Token) error {
	if err := db.Delete(token).Error; err != nil {
		return err
	}

	return nil
}

func GetToken(id util.Hash) (*Token, error) {
	token := new(Token)

	if err := db.Where("id = ?", id.String()).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func GetTokenBase64(key string) (*Token, error) {
	if id, err := base64.URLEncoding.DecodeString(key); err != nil {
		return nil, err
	} else {
		return GetToken(util.Hash(id))
	}
}

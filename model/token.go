package model

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/tkusd/server/model/types"
)

// Token represents the data structure of a token.
type Token struct {
	ID        types.Base64Hash `json:"id"`
	UserID    types.UUID       `json:"user_id"`
	CreatedAt types.Time       `json:"created_at"`
}

// BeforeCreate is called when the data is about to be created.
func (t *Token) BeforeCreate() error {
	t.CreatedAt = types.Now()
	return nil
}

// Save creates or updates data in the database.
func (t *Token) Save() error {
	if t.ID.IsEmpty() {
		t.ID = types.Base64Hash{types.SHA256(t.UserID.String(), t.CreatedAt.String(), uuid.New())}
		return db.Create(t).Error
	}

	return db.Save(t).Error
}

// Delete deletes data from the database.
func (t *Token) Delete() error {
	return db.Delete(t).Error
}

// GetToken returns the token data.
func GetToken(id types.Base64Hash) (*Token, error) {
	token := new(Token)

	if err := db.Where("id = ?", id.HexString()).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

// GetTokenBase64 searchs token data by Base64 ID.
func GetTokenBase64(key string) (*Token, error) {
	hash, err := types.DecodeBase64(key)

	if err != nil {
		return nil, err
	}

	return GetToken(*hash)
}

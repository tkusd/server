package model

import (
	"math/rand"
	"time"

	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

const (
	secretLength = 64
)

// Token represents the data structure of a token.
type Token struct {
	ID        types.UUID       `json:"id"`
	UserID    types.UUID       `json:"user_id"`
	Secret    types.Base64Hash `json:"secret"`
	CreatedAt types.Time       `json:"created_at"`
	UpdatedAt types.Time       `json:"updated_at"`
}

func (t *Token) WithoutSecret() map[string]interface{} {
	return map[string]interface{}{
		"id":         t.ID,
		"user_id":    t.UserID,
		"created_at": t.CreatedAt,
		"updated_at": t.UpdatedAt,
	}
}

func (t *Token) BeforeCreate() error {
	randomStr := util.RandStringBytesMaskImprSrc("0123456789abcdef", rand.NewSource(time.Now().UnixNano()), secretLength)
	hash, err := types.DecodeHash(randomStr)

	if err != nil {
		return err
	}

	t.Secret = types.Base64Hash{hash}

	return nil
}

// Save creates or updates data in the database.
func (t *Token) Save() error {
	return db.LogMode(true).Save(t).Error
}

// Delete deletes data from the database.
func (t *Token) Delete() error {
	return db.Delete(t).Error
}

// GetToken returns the token data.
func GetToken(id types.UUID) (*Token, error) {
	token := new(Token)

	if err := db.Where("id = ?", id.String()).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

func GetTokenBySecret(str string) (*Token, error) {
	hash, err := types.DecodeBase64(str)

	if err != nil {
		return nil, err
	}

	token := new(Token)

	if err := db.Where("secret = ?", hash.HexString()).First(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

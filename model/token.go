package model

import "github.com/tkusd/server/model/types"

// Token represents the data structure of a token.
type Token struct {
	ID        types.UUID `json:"id"`
	UserID    types.UUID `json:"user_id"`
	CreatedAt types.Time `json:"created_at"`
	UpdatedAt types.Time `json:"updated_at"`
}

// Save creates or updates data in the database.
func (t *Token) Save() error {
	return db.Save(t).Error
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

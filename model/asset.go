package model

import (
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

type Asset struct {
	ID          types.UUID `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProjectID   types.UUID `json:"project_id"`
	CreatedAt   types.Time `json:"created_at"`
	UpdatedAt   types.Time `json:"updated_at"`
	Size        int64      `json:"size"`
	Type        string     `json:"type"`
	Slug        string     `json:"-"`
	Width       int        `json:"width,omitempty"`
	Height      int        `json:"height,omitempty"`
	Hash        types.Hash `json:"hash"`
}

func (asset *Asset) Save() error {
	asset.Name = govalidator.Trim(asset.Name, "")

	if len(asset.Name) > 255 {
		return &util.APIError{
			Field:   "name",
			Code:    util.LengthError,
			Message: "Maximum length of name is 255.",
		}
	}

	return db.Save(asset).Error
}

func (asset *Asset) Delete() error {
	if err := asset.DeleteAsset(); err != nil {
		return err
	}

	return db.Delete(asset).Error
}

func (asset *Asset) DeleteAsset() error {
	path := util.GetAssetFilePath(asset.Slug)

	// File does not exist. Skip deletion
	if !util.IsAssetExist(asset.Slug) {
		return nil
	}

	return os.Remove(path)
}

func (asset *Asset) Exists() bool {
	return exists("assets", asset.ID.String())
}

func GetAssetList(projectID types.UUID) ([]*Asset, error) {
	var assets []*Asset

	if err := db.Where("project_id = ?", projectID.String()).Order("created_at").Find(&assets).Error; err != nil {
		return nil, err
	}

	if assets == nil {
		assets = make([]*Asset, 0)
	}

	return assets, nil
}

func GetAsset(id types.UUID) (*Asset, error) {
	asset := new(Asset)

	if err := db.Where("id = ?", id.String()).First(asset).Error; err != nil {
		return nil, err
	}

	return asset, nil
}

package model

import (
	"database/sql"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

var (
	rAssetBase = regexp.MustCompile(`^(.+?)(?: *\((\d+)\))?$`)
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

func (asset *Asset) BeforeSave(tx *gorm.DB) error {
	asset.Name = govalidator.Trim(asset.Name, "")

	if asset.Name == "" {
		return &util.APIError{
			Code:    util.RequiredError,
			Field:   "name",
			Message: "Name is required.",
		}
	}

	var similarName string
	var serial int

	// Get the base name
	ext := filepath.Ext(asset.Name)
	originalBase := asset.Name[:len(asset.Name)-len(ext)]
	base := originalBase
	match := rAssetBase.FindStringSubmatch(base)

	// Get the serial number from the base name
	if match[2] != "" {
		base = util.Slugize(match[1])
		serial, _ = strconv.Atoi(match[2])
	} else {
		base = util.Slugize(base)
	}

	// Build the regular expression for search
	exp := "^" + regexp.QuoteMeta(base) + " *(\\(\\d+\\))?" + regexp.QuoteMeta(ext) + "$"

	// Search for the simliar name in the database
	scope := tx.Table("assets").
		Select("name").
		Where("name ~ ? AND project_id = ?", exp, asset.ProjectID.String()).
		Order("name desc").
		Limit(1)

	// Exclude the asset itself
	if asset.ID.Valid() {
		scope = scope.Not("id", asset.ID.String())
	}

	err := scope.Row().
		Scan(&similarName)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Update the serial number if there's a simliar name
	if similarName != "" {
		similarExt := filepath.Ext(similarName)
		similarBase := similarName[:len(similarName)-len(similarExt)]
		match := rAssetBase.FindStringSubmatch(similarBase)

		if match[2] != "" {
			serial, _ = strconv.Atoi(match[2])
			serial++
		} else {
			serial = 1
		}
	}

	// Append the serial number to the asset name
	if serial > 0 {
		asset.Name = base + " (" + strconv.Itoa(serial) + ")" + ext
	} else {
		asset.Name = base + ext
	}

	return nil
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

package v1

import (
	"bytes"
	"crypto/sha1"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"image"
	// Image packages
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

func AssetList(c *gin.Context) error {
	projectID, err := GetIDParam(c, projectIDParam)

	if err != nil {
		return err
	}

	list, err := model.GetAssetList(*projectID)

	if err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, list)
}

type assetForm struct {
	Name        *string               `json:"name"`
	Description *string               `json:"description"`
	Data        *multipart.FileHeader `json:"data"`
}

func (form *assetForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&form.Name:        "name",
		&form.Description: "description",
		&form.Data:        "data",
	}
}

func saveAsset(form *assetForm, asset *model.Asset) error {
	if form.Name != nil {
		asset.Name = *form.Name
	}

	if form.Description != nil {
		asset.Description = *form.Description
	}

	if form.Data != nil {
		var err error

		// Set the asset name
		if form.Name == nil && !asset.ID.Valid() {
			if asset.Name, err = url.QueryUnescape(form.Data.Filename); err != nil {
				return err
			}
		}

		// Detect the mime type
		extname := filepath.Ext(form.Data.Filename)
		asset.Type = mime.TypeByExtension(extname)

		// Open the multipart file stream
		var fh io.ReadCloser

		if fh, err = form.Data.Open(); err != nil {
			return err
		}

		defer fh.Close()

		// Read the file into a buffer
		buf := bytes.Buffer{}

		if asset.Size, err = buf.ReadFrom(fh); err != nil {
			return err
		}

		// Read the image dimensions
		switch asset.Type {
		case "image/png", "image/jpeg", "image/gif":
			var img image.Image
			reader := bytes.NewReader(buf.Bytes())

			if img, _, err = image.Decode(reader); err != nil {
				return err
			}

			size := img.Bounds().Size()
			asset.Width = size.X
			asset.Height = size.Y

			break
		}

		// Create the upload directory
		if err := util.EnsureUploadDir(); err != nil {
			return err
		}

		// Write the file
		asset.Slug = types.NewRandomUUID().String() + extname
		uploadPath := util.GetUploadFilePath(asset.Slug)
		file, err := os.Create(uploadPath)
		h := sha1.New()
		writer := io.MultiWriter(file, h)

		if err != nil {
			return err
		}

		defer file.Close()

		if _, err := buf.WriteTo(writer); err != nil {
			return err
		}

		asset.Hash = h.Sum(nil)
	}

	return asset.Save()
}

func AssetCreate(c *gin.Context) error {
	project, err := GetProject(c)

	if err != nil {
		return err
	}

	if err := CheckUserPermission(c, project.UserID); err != nil {
		return err
	}

	form := new(assetForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	if form.Data == nil {
		return &util.APIError{
			Code:    util.RequiredError,
			Field:   "data",
			Message: "Data is required.",
		}
	}

	asset := &model.Asset{
		ProjectID: project.ID,
	}

	if err := saveAsset(form, asset); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusCreated, asset)
}

func AssetShow(c *gin.Context) error {
	asset, err := GetAsset(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, asset.ProjectID, false); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, asset)
}

func AssetUpdate(c *gin.Context) error {
	form := new(assetForm)

	if err := common.BindForm(c, form); err != nil {
		return err
	}

	asset, err := GetAsset(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, asset.ProjectID, true); err != nil {
		return err
	}

	if err := saveAsset(form, asset); err != nil {
		return err
	}

	return common.APIResponse(c, http.StatusOK, asset)
}

func AssetDestroy(c *gin.Context) error {
	asset, err := GetAsset(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, asset.ProjectID, true); err != nil {
		return err
	}

	if err := asset.Delete(); err != nil {
		return err
	}

	c.Writer.WriteHeader(http.StatusNoContent)
	return nil
}

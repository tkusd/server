package v1

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"image"
	// Image packages
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/daddye/vips"
	"github.com/gin-gonic/gin"
	"github.com/mholt/binding"
	"github.com/tkusd/server/controller/common"
	"github.com/tkusd/server/model"
	"github.com/tkusd/server/model/types"
	"github.com/tkusd/server/util"
)

const (
	headerETag            = "ETag"
	headerCacheControl    = "Cache-Control"
	headerLastModified    = "Last-Modified"
	headerIfModifiedSince = "If-Modified-Since"
	headerIfNoneMatch     = "If-None-Match"
	headerContentType     = "Content-Type"

	defaultThumbSize = "medium"
)

var thumbSize = map[string]int{
	"small":  160,
	"medium": 320,
	"large":  640,
	"huge":   1024,
}

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

		if asset.Slug != "" {
			filepath := util.GetAssetFilePath(asset.Slug)

			if err = os.Remove(filepath); err != nil {
				return err
			}
		}

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

		// Write the file
		asset.Slug = types.NewRandomUUID().String() + extname
		uploadPath := util.GetAssetFilePath(asset.Slug)
		dir, _ := filepath.Split(uploadPath)

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

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

func shouldUseAssetCache(c *gin.Context, asset *model.Asset) bool {
	if etag := c.Request.Header.Get(headerIfNoneMatch); etag != "" {
		if unquotedEtag, err := strconv.Unquote(etag); err == nil {
			if hash, err := base64.StdEncoding.DecodeString(unquotedEtag); err == nil {
				return bytes.Compare(hash, asset.Hash) == 0
			}
		}
	}

	if modified := c.Request.Header.Get(headerIfModifiedSince); modified != "" {
		if modifiedTime, err := http.ParseTime(modified); err == nil {
			return modifiedTime.Truncate(time.Second).Equal(asset.UpdatedAt.Time.Truncate(time.Second))
		}
	}

	return false
}

func addAssetBlobCacheHeader(c *gin.Context, asset *model.Asset) {
	etag := base64.StdEncoding.EncodeToString(asset.Hash)

	c.Header(headerETag, strconv.Quote(etag))
	c.Header(headerCacheControl, "private, must-revalidate, max-age=31536000") // 1 year
	c.Header(headerLastModified, asset.UpdatedAt.UTC().Format(http.TimeFormat))
}

func AssetBlob(c *gin.Context) error {
	asset, err := GetAsset(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, asset.ProjectID, false); err != nil {
		return err
	}

	addAssetBlobCacheHeader(c, asset)

	if shouldUseAssetCache(c, asset) {
		c.Writer.WriteHeader(http.StatusNotModified)
		return nil
	}

	path := util.GetAssetFilePath(asset.Slug)

	http.ServeFile(c.Writer, c.Request, path)
	return nil
}

func AssetThumbnail(c *gin.Context) error {
	asset, err := GetAsset(c)

	if err != nil {
		return err
	}

	if err := CheckProjectPermission(c, asset.ProjectID, false); err != nil {
		return err
	}

	addAssetBlobCacheHeader(c, asset)

	if shouldUseAssetCache(c, asset) {
		c.Writer.WriteHeader(http.StatusNotModified)
		return nil
	}

	switch asset.Type {
	case "image/jpeg", "image/png", "image/gif":
		var size, width, height int
		var ok bool
		path := util.GetAssetFilePath(asset.Slug)

		if size, ok = thumbSize[c.Query("size")]; !ok {
			size = thumbSize[defaultThumbSize]
		}

		if asset.Width > asset.Height {
			width = size
		} else {
			height = size
		}
		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		buf, _ := ioutil.ReadAll(file)
		thumb, err := vips.Resize(buf, vips.Options{
			Width:        width,
			Height:       height,
			Crop:         false,
			Extend:       vips.EXTEND_WHITE,
			Interpolator: vips.BILINEAR,
			Gravity:      vips.CENTRE,
			Quality:      70,
		})

		if err != nil {
			return err
		}

		c.Header(headerContentType, "image/jpeg")
		c.Writer.Write(thumb)
		return nil
	}

	return &util.APIError{
		Code:    util.ContentTypeError,
		Message: "Thumbnail only support for image/jpeg, image/png and image/gif content type.",
	}
}

package handlers

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"strings"

	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/pkg/memer"
	"github.com/labstack/echo/v5"
)

const (
	maxFileSize = 1 * 1024 * 1024 // 1 mb max
)

const (
	ErrorText500         = "Cant create file"
	ErrorTextUnknownType = "Unknown file type, only png and jpeg are available"
)

var mimeToType = map[string]memer.ImageType{
	"image/png":  memer.Png,
	"image/jpeg": memer.Jpeg,
}

type CreateMemeHandler struct {
	mem *memer.Memer
}

func NewCreateMemeHandler(mem *memer.Memer) *CreateMemeHandler {
	return &CreateMemeHandler{
		mem: mem,
	}
}

func (h *CreateMemeHandler) CreateMeme(c *echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return customerrors.NewAppHTTPError(http.StatusNotFound, "File not found", err)
	}

	imgType, err := validateImage(file)
	if err != nil {
		return err
	}

	imgSrc, err := decodeImage(file, imgType)
	if err != nil {
		return customerrors.NewAppHTTPError(http.StatusInternalServerError, ErrorText500, err)
	}

	text := c.FormValue("text")
	if text == "" || strings.TrimSpace(text) == "" {
		return customerrors.NewAppHTTPError(http.StatusNotFound, "Text not provided", nil)
	}

	resultImg, err := h.mem.Generate(imgSrc, strings.ToUpper(text))
	if err != nil {
		return customerrors.NewAppHTTPError(http.StatusInternalServerError, ErrorText500, err)
	}

	c.Response().Header().Set("Content-Type", "image/jpeg")
	return jpeg.Encode(c.Response(), resultImg, &jpeg.Options{Quality: 90})
}

func validateImage(f *multipart.FileHeader) (memer.ImageType, error) {
	if f.Size > maxFileSize {
		return "", customerrors.NewAppHTTPError(
			http.StatusForbidden, "Max file size should be less than 1MB", nil,
		)
	}

	ct := f.Header.Get("Content-Type")
	if ct == "" {
		return "", customerrors.NewAppHTTPError(
			http.StatusForbidden,
			"Unknown file type",
			nil,
		)
	}

	imgType, ok := mimeToType[ct]
	if !ok {
		return "", customerrors.NewAppHTTPError(
			http.StatusForbidden,
			"Unsupported file type, only PNG and JPEG are allowed",
			nil,
		)
	}

	return imgType, nil
}

func decodeImage(fh *multipart.FileHeader, imgType memer.ImageType) (image.Image, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var img image.Image

	switch imgType {
	case memer.Jpeg:
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	case memer.Png:
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported image type")
	}

	return img, nil
}

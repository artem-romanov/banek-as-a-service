package handlers

import (
	"errors"
	"net/http"

	"baneks.com/internal/api/baneks/dto"
	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/loaders/banekloader"
	util "baneks.com/internal/utils"
	"github.com/labstack/echo/v4"
)

type HandlerRequest struct {
	Slug string `param:"slug"`
}

func GetBanekBySlug(c echo.Context) error {
	request := new(HandlerRequest)
	if err := c.Bind(request); err != nil {
		customerrors.NewAppHTTPError(http.StatusBadRequest, "Incorrect slug", err)
	}
	if err := util.Validate(c, request); err != nil {
		return err
	}
	loader := banekloader.NewBaneksSiteLoader()
	banek, err := loader.GetBanekBySlug(request.Slug)
	if err != nil {
		var notFoundError *customerrors.NotFoundRequestError
		switch {
		case errors.As(err, &notFoundError):
			return customerrors.NewAppHTTPError(http.StatusNotFound, "Banek not found", err)
		default:
			return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Banek download error", err)
		}
	}

	return c.JSON(http.StatusOK, dto.BanekToResponse(&banek))
}

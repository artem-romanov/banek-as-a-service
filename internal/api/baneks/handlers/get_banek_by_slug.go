package handlers

import (
	"errors"
	"net/http"

	"baneks.com/internal/api/baneks/dto"
	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/loaders/banekloader"
	customvalidator "baneks.com/internal/utils/validator"
	"github.com/labstack/echo/v5"
)

type HandlerRequest struct {
	Slug string `param:"slug"`
}

func GetBanekBySlug(c *echo.Context) error {
	ctx := c.Request().Context()

	requestParams := new(HandlerRequest)
	if err := c.Bind(requestParams); err != nil {
		customerrors.NewAppBindError(err)
	}
	httpError := customvalidator.ValidateRequest(c.Validate, requestParams)
	if httpError != nil {
		return httpError
	}
	loader := banekloader.NewBaneksSiteLoader()
	banek, err := loader.GetBanekBySlug(ctx, requestParams.Slug)
	if err != nil {
		if _, ok := errors.AsType[*customerrors.NotFoundRequestError](err); ok {
			return customerrors.NewAppHTTPError(http.StatusNotFound, "Banek not found", err)
		}

		return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Banek download error", err)
	}

	return c.JSON(http.StatusOK, dto.BanekToResponse(&banek))
}

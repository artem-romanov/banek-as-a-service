package handlers

import (
	"errors"
	"net/http"

	"baneks.com/internal/api/baneks/dto"
	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/loaders/banekloader"
	"github.com/labstack/echo/v5"
)

func GetRandomBanek(c *echo.Context) error {
	ctx := c.Request().Context()

	balancer := banekloader.GetBalancer()
	banekLoader := balancer.GetLoader()
	banek, err := banekLoader.GetRandomBanek(ctx)
	if err != nil {
		if _, ok := errors.AsType[*customerrors.NotFoundRequestError](err); ok {
			return customerrors.NewAppHTTPError(http.StatusNotFound, "Banek not found", err)
		}
		return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Banek download error", err)
	}

	banekResponse := dto.BanekToResponse(&banek)
	return c.JSON(http.StatusOK, banekResponse)
}

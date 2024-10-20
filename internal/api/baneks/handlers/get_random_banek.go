package handlers

import (
	"errors"
	"net/http"

	"baneks.com/internal/api/baneks/dto"
	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/loaders/banekloader"
	"github.com/labstack/echo/v4"
)

func GetRandomBanek(c echo.Context) error {
	balancer := banekloader.GetBalancer()
	banekLoader := balancer.GetLoader()
	banek, err := banekLoader.GetRandomBanek()
	if err != nil {
		var notFoundError *customerrors.NotFoundRequestError
		switch {
		case errors.As(err, &notFoundError):
			return customerrors.NewAppHTTPError(http.StatusNotFound, "Banek not found", err)
		default:
			return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Banek download error", err)
		}
	}

	banekResponse := dto.BanekToResponse(&banek)
	return c.JSON(http.StatusOK, banekResponse)
}

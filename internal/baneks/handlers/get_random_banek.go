package handlers

import (
	"net/http"

	"baneks.com/internal/banek_loader"
	"baneks.com/internal/baneks/dto"
	"github.com/labstack/echo/v4"
)

func GetRandomBanek(c echo.Context) error {
	balancer := banek_loader.GetBalancer()
	banekLoader := balancer.GetLoader()
	banek, err := banekLoader.GetRandomBanek()
	if err != nil {
		return echo.ErrBadRequest
	}

	banekResponse := dto.BanekToResponse(&banek)
	return c.JSON(http.StatusOK, banekResponse)
}

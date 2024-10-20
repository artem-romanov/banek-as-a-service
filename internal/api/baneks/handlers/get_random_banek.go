package handlers

import (
	"net/http"

	"baneks.com/internal/api/baneks/dto"
	"baneks.com/internal/loaders/banekloader"
	"github.com/labstack/echo/v4"
)

func GetRandomBanek(c echo.Context) error {
	banekLoader := banekloader.NewBanekRuLoader()
	// banekLoader := banekloader.NewBaneksSiteLoader()
	// balancer := banekloader.GetBalancer()
	// banekLoader := balancer.GetLoader()
	banek, err := banekLoader.GetRandomBanek()
	if err != nil {
		return echo.ErrBadRequest
	}

	banekResponse := dto.BanekToResponse(&banek)
	return c.JSON(http.StatusOK, banekResponse)
}

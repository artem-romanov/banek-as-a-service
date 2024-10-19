package handlers

import (
	"net/http"

	"baneks.com/internal/baneks/dto"
	"baneks.com/internal/scraper"
	"github.com/labstack/echo/v4"
)

func GetRandomBanek(c echo.Context) error {
	banek, err := scraper.GetRandomBanek()
	if err != nil {
		return echo.ErrBadRequest
	}

	banekResponse := dto.BanekToResponse(&banek)
	return c.JSON(http.StatusOK, banekResponse)
}

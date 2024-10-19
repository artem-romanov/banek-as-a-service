package handlers

import (
	"net/http"

	"baneks.com/internal/baneks/dto"
	"baneks.com/internal/scraper"
	util "baneks.com/internal/utils"
	"github.com/labstack/echo/v4"
)

type HandlerRequest struct {
	Slug string `param:"slug"`
}

func GetBanekBySlug(c echo.Context) error {
	request := new(HandlerRequest)
	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect slug")
	}
	if err := util.Validate(c, request); err != nil {
		return err
	}

	banek, err := scraper.GetBanekBySlug(request.Slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Banek not found")
	}

	return c.JSON(http.StatusOK, dto.BanekToResponse(&banek))
}

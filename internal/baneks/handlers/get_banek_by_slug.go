package handlers

import (
	"net/http"

	"baneks.com/internal/banek_loader"
	"baneks.com/internal/baneks/dto"
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
	loader := banek_loader.NewBaneksSiteLoader()
	banek, err := loader.GetBanekBySlug(request.Slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Banek not found")
	}

	return c.JSON(http.StatusOK, dto.BanekToResponse(&banek))
}

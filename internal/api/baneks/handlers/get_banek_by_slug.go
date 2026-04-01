package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"baneks.com/internal/api/baneks/dto"
	"baneks.com/internal/api/middlewares"
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

	fmt.Println("inside slug!")
	slog.Info("fuck!")
	data, ok := middlewares.GetUser(ctx)
	if !ok {
		fmt.Println("data is nil:", data)
	} else {
		fmt.Println("data is:", data)
	}

	timeCh := time.After(time.Second * 4)

	select {
	case <-ctx.Done():
		return customerrors.NewAppHTTPError(500, "canceled context", errors.New("ctx canceled"))
	case <-timeCh:
		// noop
	}

	requestParams := new(HandlerRequest)
	if err := c.Bind(requestParams); err != nil {
		return customerrors.NewAppBindError(err)
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

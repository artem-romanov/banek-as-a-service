package handlers

import (
	"errors"
	"net/http"

	"baneks.com/internal/api/memes/dto"
	customerrors "baneks.com/internal/custom_errors"
	memesloader "baneks.com/internal/loaders/memes_loader"
	customvalidator "baneks.com/internal/utils/validator"
	"github.com/labstack/echo/v5"
)

// Struct for handling the request parameters for the GetRandomMemes endpoint.
// The 'year' parameter can be any integer from the current year to 2015. Non-required.
type GetRandomMemesRequest struct {
	Year int `query:"year" validate:"omitempty,is-correct-meme-year"`
}

func GetRandomMemes(c *echo.Context) error {
	var requestParams GetRandomMemesRequest

	err := c.Bind(&requestParams)
	if err != nil {
		return customerrors.NewAppBindError(err)
	}
	httpError := customvalidator.ValidateRequest(c.Validate, requestParams)
	if httpError != nil {
		return httpError
	}

	var memeLoader memesloader.MemeLoader = memesloader.NewQablydauMemeLoader()
	memes, err := memeLoader.GetRandomMemesWithConfig(memesloader.RandomMemesConfig{
		Year: requestParams.Year,
	})
	if err != nil {
		var notFoundError *customerrors.NotFoundRequestError
		switch {
		case errors.As(err, &notFoundError):
			return customerrors.NewAppHTTPError(http.StatusNotFound, "Memes not found", err)
		default:
			return customerrors.NewAppHTTPError(http.StatusInternalServerError, "Memes loading error", err)
		}
	}
	responseMemes := []dto.MemeResponse{}

	for _, meme := range memes {
		responseMemes = append(responseMemes, dto.MemeResponse{
			ImageUri:        meme.ImageUri,
			OriginalPostUri: meme.OriginalUri,
		})
	}

	finalResponse := dto.MemesResponse{
		Memes: responseMemes,
	}

	return c.JSON(http.StatusOK, finalResponse)
}

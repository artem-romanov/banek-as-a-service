package handlers

import (
	"errors"
	"net/http"

	"baneks.com/internal/api/memes/dto"
	customerrors "baneks.com/internal/custom_errors"
	memesloader "baneks.com/internal/loaders/memes_loader"
	"github.com/labstack/echo/v4"
)

func GetRandomMemes(c echo.Context) error {
	var memeLoader memesloader.MemeLoader = memesloader.NewQablydauMemeLoader()
	memes, err := memeLoader.GetRandomMemes()
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

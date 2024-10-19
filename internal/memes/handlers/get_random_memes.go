package handlers

import (
	"net/http"

	memesloader "baneks.com/internal/loaders/memes_loader"
	"baneks.com/internal/memes/dto"
	"github.com/labstack/echo/v4"
)

func GetRandomMemes(c echo.Context) error {
	var memeLoader memesloader.MemeLoader = memesloader.NewQablydauMemeLoader()
	memes, err := memeLoader.GetRandomMemes()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Memes download error")
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

package memesloader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/model"
	"github.com/PuerkitoBio/goquery"
)

const (
	minimumMemeYear int = 2015
)

const (
	baseMemesUri string = "https://idiod.qabyldau.com"
)

type QablydauMemeLoader struct {
	baseUri string
}

// Base response for pages such as /random, /random/:year, /top, /top/:year, etc
type JsonResponse struct {
	Props struct {
		Items struct {
			Data []struct {
				Id       uint   `json:"id"`
				User     string `json:"user"`
				PostLink string `json:"post_link"`
				Path     string `json:"path"`
			}
		} `json:"items"`
	} `json:"props"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

func NewQablydauMemeLoader() *QablydauMemeLoader {
	return &QablydauMemeLoader{
		baseUri: baseMemesUri,
	}
}

/// GET RANDOM MEMES FUNCTIONS

var DefaultGetRandomMemesConfig = RandomMemesConfig{}

// GetRandomMemesWithConfig returns a slice of random memes from the provided year.
//
// This is a more versatile version of GetRandomMemes, which allows you to
// customize the request.
//
// The function will return an error if the provided year is not inbetween
// the current year and the minimum available year on the website.
//
// The function will also return an error if there was a problem with the
// request or parsing the response.
func (loader *QablydauMemeLoader) GetRandomMemesWithConfig(config RandomMemesConfig) ([]model.Meme, error) {
	// right now this is ok just to take year value without any checks
	// validator of GetRandomMemes will check it for us
	year := config.Year

	// but in future config updates - check them there
	// ...

	return loader.getRandomMemes(year)
}

// GetRandomMemes returns a slice of random memes from any year.
//
// This is just a wrapper around GetRandomMemesWithConfig with default config.
func (loader *QablydauMemeLoader) GetRandomMemes() ([]model.Meme, error) {
	return loader.GetRandomMemesWithConfig(DefaultGetRandomMemesConfig)
}

// Gets memes based on provided args.
// This is the base functions, use provded wrappers to get memes.
func (loader *QablydauMemeLoader) getRandomMemes(year int) ([]model.Meme, error) {
	var uri string
	if year == 0 { // default unititilized value for a year
		uri = loader.baseUri + "/random"
	} else if year > time.Now().Year() || year < minimumMemeYear {
		return nil, &customerrors.InvalidInputError{
			Err: fmt.Errorf("invalid year, should be inbetween now and %d. provided: %d", minimumMemeYear, year),
		}
	} else {
		uri = loader.baseUri + "/random/" + strconv.Itoa(year)
	}

	response, err := http.Get(uri)
	if err != nil {
		return nil, &customerrors.HttpNetworkError{
			Err: err,
			Uri: response.Request.URL.String(),
		}
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, &customerrors.NotFoundRequestError{
			Uri: response.Request.URL.String(),
		}
	}
	if response.StatusCode != http.StatusOK {
		return nil, &customerrors.DownloadRequestError{
			Uri:        response.Request.URL.String(),
			StatusCode: response.StatusCode,
		}
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, &customerrors.ParseDataError{
			Err: err,
		}
	}

	app := doc.Find("#app").First()
	data, exists := app.Attr("data-page")
	if !exists {
		return nil, &customerrors.ParseDataError{
			Err: fmt.Errorf("data-page attribute not found on %s", response.Request.URL.String()),
		}
	}

	finalData := &JsonResponse{}
	err = json.Unmarshal([]byte(data), &finalData)
	if err != nil {
		return nil, &customerrors.ParseDataError{
			Err: fmt.Errorf("can't unmarshal JSON: %w", err),
		}
	}

	memes := []model.Meme{}
	for _, item := range finalData.Props.Items.Data {
		meme := model.Meme{
			OriginalId:  item.Id,
			OriginalUri: item.PostLink,
			ImageUri:    item.Path,
		}
		memes = append(memes, meme)
	}
	return memes, nil
}

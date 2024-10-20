package memesloader

import (
	"encoding/json"
	"fmt"
	"net/http"

	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/model"
	"github.com/PuerkitoBio/goquery"
)

type QablydauMemeLoader struct {
	uri string
}

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
		uri: "https://idiod.qabyldau.com",
	}
}

func (loader *QablydauMemeLoader) GetRandomMemes() ([]model.Meme, error) {
	uri := loader.uri + "/random"
	finalData := &JsonResponse{}
	response, err := http.Get(uri)

	if err != nil {
		return nil, &customerrors.HttpNetworkError{
			Err: err,
			Uri: response.Request.RequestURI,
		}
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return nil, &customerrors.NotFoundRequestError{
			Uri: response.Request.RequestURI,
			Err: err,
		}
	}
	if response.StatusCode != http.StatusOK {
		return nil, &customerrors.DownloadRequestError{
			Uri:        response.Request.RequestURI,
			StatusCode: response.StatusCode,
			Err:        err,
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
			Err: fmt.Errorf("data-page not found in %s", response.Request.URL),
		}
	}

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

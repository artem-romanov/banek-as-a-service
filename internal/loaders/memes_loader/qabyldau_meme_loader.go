package memesloader

import (
	"encoding/json"
	"errors"
	"net/http"

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
		return nil, errors.New("error loading memes")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("memes not found")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, errors.New("error parsing Body")
	}

	app := doc.Find("#app").First()
	data, exists := app.Attr("data-page")
	if !exists {
		return nil, errors.New("data field not exists in html")
	}

	err = json.Unmarshal([]byte(data), &finalData)
	if err != nil {
		return nil, err
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

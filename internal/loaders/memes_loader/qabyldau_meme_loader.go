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

func NewQablydauMemeLoader() *QablydauMemeLoader {
	return &QablydauMemeLoader{
		uri: "https://idiod.qabyldau.com",
	}
}

func (loader *QablydauMemeLoader) GetRandomMemes() ([]model.Meme, error) {
	uri := loader.uri + "/random"
	var finalData JsonResponse = JsonResponse{}
	response, err := http.Get(uri)
	if err != nil {
		return nil, errors.New("Error loading memes")
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Memes not found")
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, errors.New("Error parsing Body")
	}

	app := doc.Find("#app").First()
	data, exists := app.Attr("data-page")
	if exists == false {
		return nil, errors.New("Memes not found")
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

// TODO: Messy code, think where to put these structs later
type Item struct {
	Id       uint   `json:"id"`
	User     string `json:"user"`
	PostLink string `json:"post_link"`
	Path     string `json:"path"`
}

type Items struct {
	Data []Item `json:"data"`
}

type Props struct {
	Items Items `json:"items"`
}
type JsonResponse struct {
	Props   Props  `json:"props"`
	Version string `json:"version"`
	Url     string `json:"url"`
}

package banek_loader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"baneks.com/internal/models"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type BaneksRuLoader struct {
	siteUri string
}

var banekRuUri string = "https://baneks.ru"

func NewBanekRuLoader() *BaneksRuLoader {
	return &BaneksRuLoader{
		siteUri: banekRuUri,
	}
}

func (loader *BaneksRuLoader) GetRandomBanek() (models.Banek, error) {
	response, err := http.Get(loader.siteUri + "/random")
	if err != nil {
		return models.Banek{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return models.Banek{}, errors.New(
			fmt.Sprintf(
				"error loading banek, status code: %d",
				response.StatusCode,
			),
		)
	}
	banek, err := loader.extractBanekFromBody(response.Body)
	if err != nil {

	}

	return banek, nil
}

func (loader *BaneksRuLoader) extractBanekFromBody(body io.ReadCloser) (models.Banek, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return models.Banek{}, err
	}

	banekText, _ := loader.extractText(doc)
	banekLikes, err := loader.extractLikes(doc)
	if err != nil {
		// There might be a better way,
		// but for now likes count is not so important.
		// If text presented - we dont really matter about likes
		banekLikes = 0
	}

	return models.Banek{
		Text:  banekText,
		Likes: banekLikes,
	}, nil

}

func (loader *BaneksRuLoader) extractText(doc *goquery.Document) (string, error) {
	selector := ".anek-view article > p"
	rawText := doc.Find(selector).First()
	var textBuilder strings.Builder

	rawText.Contents().Each(func(i int, s *goquery.Selection) {
		switch s.Nodes[0].Type {
		case html.TextNode:
			textBuilder.WriteString(s.Text())
		case html.ElementNode:
			if s.Nodes[0].Data == "br" {
				textBuilder.WriteString("\n")
			}
		}
	})

	return strings.TrimSpace(textBuilder.String()), nil
}

func (loader *BaneksRuLoader) extractLikes(doc *goquery.Document) (int, error) {
	selector := ".rating-counter"
	strLikes := doc.Find(selector).First().Text()

	likes, err := strconv.Atoi(strLikes)
	if err != nil {
		return -1, err
	}
	return likes, nil
}

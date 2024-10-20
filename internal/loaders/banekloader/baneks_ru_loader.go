package banekloader

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	customerrors "baneks.com/internal/custom_errors"
	"baneks.com/internal/model"
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

func (loader *BaneksRuLoader) GetRandomBanek() (model.Banek, error) {
	response, err := http.Get(loader.siteUri + "/random")
	if err != nil {
		return model.Banek{}, &customerrors.HttpNetworkError{
			Err: err,
			Uri: response.Request.RequestURI,
		}
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return model.Banek{}, &customerrors.NotFoundRequestError{
			Uri: response.Request.RequestURI,
		}
	}
	if response.StatusCode != http.StatusOK {
		return model.Banek{}, &customerrors.DownloadRequestError{
			Uri:        response.Request.RequestURI,
			StatusCode: response.StatusCode,
		}
	}

	banek, err := loader.extractBanekFromBody(response.Body)
	if err != nil {
		return model.Banek{}, &customerrors.ParseDataError{
			Err: fmt.Errorf("Banek parsing error: %w", err),
		}
	}

	return banek, nil
}

func (loader *BaneksRuLoader) extractBanekFromBody(body io.ReadCloser) (model.Banek, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return model.Banek{}, err
	}

	banekText, _ := loader.extractText(doc)
	banekLikes, err := loader.extractLikes(doc)
	if err != nil {
		// There might be a better way,
		// but for now likes count is not so important.
		// If text presented - we dont really matter about likes
		banekLikes = 0
	}

	return model.Banek{
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
			text := s.Text()
			text = strings.TrimSpace(text)
			textBuilder.WriteString(text)
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

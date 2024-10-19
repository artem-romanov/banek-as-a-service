package banekloader

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"baneks.com/internal/model"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const banekSiteUri = "https://baneks.site"

type BaneksSiteLoader struct {
	siteUri string
}

func NewBaneksSiteLoader() *BaneksSiteLoader {
	return &BaneksSiteLoader{
		siteUri: banekSiteUri,
	}
}

func (loader *BaneksSiteLoader) GetRandomBanek() (model.Banek, error) {
	randomBanekUri := loader.siteUri + "/random"
	resp, err := http.Get(randomBanekUri)
	body := resp.Body // making explicit body so defer make sense
	defer body.Close()
	if err != nil {
		return model.Banek{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return model.Banek{}, errors.New("Banek couldn't be downloaded")
	}

	banek, err := loader.parseBanekPage(body)
	if err != nil {
		log.Printf("Random banek error download: %s", err.Error())
		return model.Banek{}, errors.New("couldn't download Banek")
	}

	return banek, nil
}

func (loader *BaneksSiteLoader) GetBanekBySlug(slug string) (model.Banek, error) {
	banekUri := banekSiteUri + "/" + slug
	resp, err := http.Get(banekUri)
	body := resp.Body // making explicit body so defer make sense
	defer body.Close()
	if err != nil {
		return model.Banek{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return model.Banek{}, errors.New("Banek not found")
	}

	banek, err := loader.parseBanekPage(resp.Body)
	if err != nil {
		return model.Banek{}, err
	}

	return banek, nil
}

func (loader *BaneksSiteLoader) parseBanekPage(body io.Reader) (banek model.Banek, err error) {
	html, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return model.Banek{}, errors.New("couldn't download Banek")
	}

	banekText := loader.extractBanekText(html)
	banekLikes, err := loader.extractBanekLikes(html)
	if err != nil {
		return model.Banek{}, errors.New("couldn't parse Banek")
	}

	return model.Banek{
		Text:  banekText,
		Likes: banekLikes,
	}, nil
}

func (loader *BaneksSiteLoader) extractBanekText(doc *goquery.Document) string {
	banekTextSelector := "article > section[itemprop='description'] > p"
	var textBuilder strings.Builder

	data := doc.Find(banekTextSelector).First()
	data.Contents().Each(func(i int, s *goquery.Selection) {
		switch s.Nodes[0].Type {
		case html.TextNode:
			textBuilder.WriteString(s.Text())
			return
		case html.ElementNode:
			if goquery.NodeName(s) == "br" {
				textBuilder.WriteString("\n")
			}
		}
	})
	return textBuilder.String()
}

func (loader *BaneksSiteLoader) extractBanekLikes(html *goquery.Document) (int, error) {
	banekLikesSelector := ".clickable.like-statistic > span.likes"
	element := html.Find(banekLikesSelector).First()
	likesStr := element.Text()
	finalLikes, err := strconv.Atoi(likesStr)
	if err != nil {
		return -1, errors.New("can't parse likes")
	}
	return finalLikes, nil
}

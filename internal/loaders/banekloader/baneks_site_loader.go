package banekloader

import (
	"context"
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

const banekSiteUri = "https://baneks.site"

type BaneksSiteLoader struct {
	siteUri string
}

func NewBaneksSiteLoader() *BaneksSiteLoader {
	return &BaneksSiteLoader{
		siteUri: banekSiteUri,
	}
}

func (loader *BaneksSiteLoader) GetRandomBanek(ctx context.Context) (model.Banek, error) {
	randomBanekUri := loader.siteUri + "/random"

	req, err := http.NewRequestWithContext(ctx, "GET", randomBanekUri, nil)
	if err != nil {
		return model.Banek{}, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Banek{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return model.Banek{}, &customerrors.NotFoundRequestError{
			Uri: req.URL.String(),
		}
	}
	if resp.StatusCode != http.StatusOK {
		return model.Banek{}, &customerrors.DownloadRequestError{
			Uri:        req.URL.String(),
			StatusCode: resp.StatusCode,
		}
	}

	banek, err := loader.parseBanekPage(resp.Body)
	if err != nil {
		return model.Banek{}, &customerrors.ParseDataError{
			Err: fmt.Errorf("Banek parsing error: %w", err),
		}
	}

	return banek, nil
}

func (loader *BaneksSiteLoader) GetBanekBySlug(ctx context.Context, slug string) (model.Banek, error) {
	banekUri := banekSiteUri + "/" + slug
	req, err := http.NewRequestWithContext(ctx, "GET", banekUri, nil)
	if err != nil {
		return model.Banek{}, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return model.Banek{}, &customerrors.HttpNetworkError{
			Err: err,
			Uri: req.URL.String(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return model.Banek{}, &customerrors.NotFoundRequestError{
			Uri: req.URL.String(),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return model.Banek{}, &customerrors.DownloadRequestError{
			Uri:        req.URL.String(),
			StatusCode: resp.StatusCode,
		}
	}

	banek, err := loader.parseBanekPage(resp.Body)
	if err != nil {
		return model.Banek{}, &customerrors.ParseDataError{
			Err: fmt.Errorf("banek parsing error: %w", err),
		}
	}

	return banek, nil
}

func (loader *BaneksSiteLoader) parseBanekPage(body io.Reader) (banek model.Banek, err error) {
	html, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return model.Banek{}, err
	}

	banekText := loader.extractBanekText(html)
	banekLikes, err := loader.extractBanekLikes(html)
	if err != nil {
		return model.Banek{}, err
	}

	return model.Banek{
		Text:  banekText,
		Likes: banekLikes,
	}, nil
}

func (loader *BaneksSiteLoader) extractBanekText(doc *goquery.Document) string {
	var textBuilder strings.Builder

	// banek.site section might contain multiple paragraphs (<p>)
	banekSectionSelector := "article > section[itemprop='description']"

	// we need only first selector, others are user comments
	section := doc.Find(banekSectionSelector).First()
	section.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Nodes[0].Type != html.ElementNode {
			return
		}
		if goquery.NodeName(s) != "p" {
			return
		}

		// inside each <p> we should get text
		// and add line breaks instead of <br />
		s.Contents().Each(func(i int, s *goquery.Selection) {
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
		textBuilder.WriteString("\n")
	})

	result := strings.TrimSpace(textBuilder.String())
	return result
}

func (loader *BaneksSiteLoader) extractBanekLikes(html *goquery.Document) (int, error) {
	banekLikesSelector := ".clickable.like-statistic > span.likes"
	element := html.Find(banekLikesSelector).First()
	likesStr := element.Text()
	finalLikes, err := strconv.Atoi(likesStr)
	if err != nil {
		return -1, err
	}
	return finalLikes, nil
}

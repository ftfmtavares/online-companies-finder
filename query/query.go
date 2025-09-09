package query

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func DuckDuckGoFirstResult(query string) (string, error) {
	searchURL := "https://duckduckgo.com/html/?q=" + url.QueryEscape(query)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:142.0) Gecko/20100101 Firefox/142.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	result := ""
	doc.Find(".result__a").First().Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			result = href
		}
	})

	if result == "" {
		return "", fmt.Errorf("no result found")
	}

	result, err = extractFinalURL(result)
	if err != nil {
		return "", err
	}

	return result, nil
}

func extractFinalURL(duckduckgoLink string) (string, error) {
	u, err := url.Parse(duckduckgoLink)
	if err != nil {
		return "", err
	}
	uddg := u.Query().Get("uddg")
	if uddg == "" {
		return duckduckgoLink, nil
	}
	finalURL, err := url.QueryUnescape(uddg)
	if err != nil {
		return "", err
	}
	return finalURL, nil
}

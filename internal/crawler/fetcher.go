package crawler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/automatedtomato/go-web-crawler/internal/ratelimiter"
)

// Functions to retrieve web pages

type Fetcher struct {
	client      *http.Client
	ratelimiter *ratelimiter.HostLimiter
	userAgent   string
}

func NewFetcher(ratelimiter *ratelimiter.HostLimiter, userAgent string) *Fetcher {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Fetcher{
		client:      client,
		ratelimiter: ratelimiter,
		userAgent:   userAgent,
	}
}

func (f *Fetcher) FetchURL(urlStr string) (*goquery.Document, error) {
	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Get host name and apply rate limit
	host := parsedURL.Host
	f.ratelimiter.Wait(host)

	// Creating HTTP request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	// Set header
	req.Header.Set("User-Agent", f.userAgent)
	req.Header.Set("Accept", "text/html")

	// Perform request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, http.ErrAbortHandler
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad status: " + resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (f *Fetcher) ExtractLink(doc *goquery.Document, baseURL string) []string {
	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return nil
	}

	var links []string
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" || strings.HasPrefix(href, "#") {
			return
		}

		// Translate relative path to absolute path
		absURL, err := f.resolveURL(baseURLParsed, href)
		if err != nil {
			return
		}

		// Add link only in the same domain
		if absURL.Host == baseURLParsed.Host {
			links = append(links, absURL.String())
		}
	})

	return links
}

func (f *Fetcher) resolveURL(base *url.URL, href string) (*url.URL, error) {
	return base.Parse(href)
}

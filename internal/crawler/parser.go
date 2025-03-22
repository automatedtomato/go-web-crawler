package crawler

import (
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/automatedtomato/go-web-crawler/internal/models"
)

// Extract articles from URL
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Extract articles from HTML
func (p *Parser) ParseArticle(doc *goquery.Document, urlStr string) (*models.Article, error) {
	article := &models.Article{
		URL:         urlStr,
		PublishedAt: time.Now(),
	}

	// Extract source(domain) from URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		article.Source = parsedURL.Host
	}

	// Extract title
	// Try general selectors
	title := doc.Find("meta[property='og:title]").AttrOr("content", "")
	if title == "" {
		title = doc.Find("meta[name='twitter:title]").AttrOr("content", "")
	}
	if title == "" {
		title = doc.Find("title").Text()
	}
	if title == "" {
		title = doc.Find("h1").First().Text()
	}
	article.Title = strings.TrimSpace(title)

	// Extract content
	// Try general content selectors
	content := ""
	contentSelectors := []string{
		"article",
		".article-body",
		".story-body",
		"[itemprop='articleBody']",
		".entry-content",
		"#content",
	}

	for _, selector := range contentSelectors {
		if content == "" {
			content = strings.TrimSpace(doc.Find(selector).First().Text())
		}
	}

	// Check if content is empty
	if content == "" {
		// Try to get all the body text
		content = strings.TrimSpace(doc.Find("body").Text())
	}
	article.Content = p.cleanText(content)

	// Extract author info
	author := doc.Find("meta[name='author']").AttrOr("content", "")
	if author == "" {
		author = doc.Find("[itemprop='author']").Text()
	}
	article.Author = strings.TrimSpace(author)

	// Extract image URL
	imageURL := doc.Find("meta[property='og:image']").AttrOr("content", "")
	if imageURL == "" {
		imageURL = doc.Find("meta[name='twitter:image']").AttrOr("content", "")
	}
	if imageURL == "" {
		// Use first large image
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			if imageURL == "" {
				src, exists := s.Attr("src")
				if exists && src != "" {
					width, _ := s.Attr("width")
					height, _ := s.Attr("height")
					// Select image that is not too small
					if width > "200" && height > "200" {
						imageURL = src
					}
				}
			}
		})
	}
	article.ImageURL = imageURL

	// Extract published date
	pubDate := doc.Find("meta[property='article:published_time']").AttrOr("content", "")
	if pubDate != "" {
		if t, err := time.Parse(time.RFC3339, pubDate); err == nil {
			article.PublishedAt = t
		}
	}

	return article, nil
}

// Helper method to clear blank spaces
func (p *Parser) cleanText(text string) string {
	// Turn multiple spaces into single space
	text = strings.Join(strings.Fields(text), " ")
	return text
}

// Check if a page is an article
func (p *Parser) IsArticlePage(doc *goquery.Document) bool {
	// Extract page features
	// 1. Has article tag
	hasArticleTag := doc.Find("article").Length() > 0

	// 2. Has article meta tag
	hasArticleMeta := doc.Find("meta[property='og:type'][content='article']").Length() > 0

	// 3. Has article structure
	hasArticleStructure := doc.Find(".article-body, .story-body, [itemprop='articleBody']").Length() > 0

	return hasArticleTag || hasArticleMeta || hasArticleStructure
}

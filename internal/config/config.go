package config

type Config struct {
	SeedURLs          []string // Initial crawl URL
	NumWorkers        int      // Num of concurrent workers
	MaxDepth          int      // Max crawl depth
	RequestsPerSecond float64
	DBPath            string
	UserAgent         string
	Timeout           int // Max crawling time (m)
}

func NewDefaultConfig() *Config {
	return &Config{
		SeedURLs: []string{
			"https://news.google.com/",
			"https://www.bbc.com/news",
		},
		NumWorkers:        5,
		MaxDepth:          3,
		RequestsPerSecond: 1,
		DBPath:            "./crawler.db",
		UserAgent:         "GoCrawler/1.0",
		Timeout:           30,
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/automatedtomato/go-web-crawler/internal/config"
	"github.com/automatedtomato/go-web-crawler/internal/crawler"
)

func main() {
	// Analyze command line arguments
	configPath := flag.String("config", "", "Path to config file (optional)")
	timeout := flag.Int("timeout", 30, "Timeout in minutes (default: 30)")
	workers := flag.Int("workers", 5, "Number of concurrent workers (default: 5)")
	depth := flag.Int("depth", 3, "Maximum crawl depth (default: 3)")
	rps := flag.Float64("rps", 1.0, "Maximum requests per second (default: 1.0)")
	dbPath := flag.String("db", "./crawler.db", "Path to database (default: ./crawler.db)")
	flag.Parse()

	// Seed URLs
	seedURLs := flag.Args()
	if len(seedURLs) == 0 {
		// Use default seed URLs
		seedURLs = []string{
			"https://news.google.com/",
			"https://www.bbc.com/news",
		}
	}

	// Initialize config
	cfg := config.NewDefaultConfig()
	if *configPath != "" {
		// TODO: Load config from file
		log.Println("TODO: Load config from file")
	}

	// Update config with command line arguments
	cfg.SeedURLs = seedURLs
	cfg.NumWorkers = *workers
	cfg.MaxDepth = *depth
	cfg.RequestsPerSecond = *rps
	cfg.DBPath = *dbPath
	cfg.Timeout = *timeout

	// Initialize crawler
	c, err := crawler.NewCrawler(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize crawler: %v", err)
	}

	// Signal handling (Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Received signal, shutting down...")
		c.Stop()
	}()

	// Start crawler
	log.Printf("Starting crawler... (worker: %d, depth: %d, timeout: %d min)", cfg.NumWorkers, cfg.MaxDepth, cfg.Timeout)
	if err != nil {
		log.Fatalf("Error while performing crawl: %v", err)
	}

	log.Println("Crawler exited successfully")
}

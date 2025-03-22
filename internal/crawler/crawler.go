package crawler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/automatedtomato/go-web-crawler/internal/config"
	"github.com/automatedtomato/go-web-crawler/internal/models"
	"github.com/automatedtomato/go-web-crawler/internal/queue"
	"github.com/automatedtomato/go-web-crawler/internal/ratelimiter"
	"github.com/automatedtomato/go-web-crawler/internal/storage"
	"golang.org/x/time/rate"
)

type Crawler struct {
	config      *config.Config
	urlQueue    *queue.URLQueue
	db          *storage.Database
	fetcher     *Fetcher
	parser      *Parser
	rateLimiter *ratelimiter.HostLimiter
	jobQueue    chan CrawlJob
	resultChan  chan *models.Article
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewCrawler(cfg *config.Config) (*Crawler, error) {
	// Initialize DB
	db, err := storage.NewDatabase(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	// Initialize URL queue
	urlQueue := queue.NewURLQueue(cfg.SeedURLs)

	// Set up rate limiter
	ratelimiter := ratelimiter.NewHostLimiter(rate.Limit(cfg.RequestsPerSecond), 1)

	// Set up context
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.Timeout)*time.Minute,
	)

	crawler := &Crawler{
		config:      cfg,
		urlQueue:    urlQueue,
		db:          db,
		rateLimiter: ratelimiter,
		jobQueue:    make(chan CrawlJob, cfg.NumWorkers*10), // buffered
		resultChan:  make(chan *models.Article, cfg.NumWorkers*5),
		ctx:         ctx,
		cancel:      cancel,
	}

	crawler.fetcher = NewFetcher(ratelimiter, cfg.UserAgent)
	crawler.parser = NewParser()

	return crawler, nil
}

func (c *Crawler) Start() error {
	log.Println("Starting crawler...")

	// Goroutine for processing results

	c.wg.Add(1)
	go c.processResults()

	// Start worker pool
	for i := 0; i < c.config.NumWorkers; i++ {
		c.wg.Add(1)
		worker := NewWorker(
			i+1,
			c.jobQueue,
			c,
			c.fetcher,
			c.parser,
			c.db,
			&c.wg,
			c.urlQueue,
			c.config.MaxDepth,
			c.resultChan,
		)
		worker.Start()
	}

	// Add seed URLs to the queue
	for _, url := range c.config.SeedURLs {
		c.SubmitURL(url, 0)
	}

	// Main loop
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				// Time out or cancel
				log.Println("Crawler exiting...")
				close(c.jobQueue)
				close(c.resultChan)
				return
			case <-time.After(5 * time.Second):
				// Exit if no URLs in queue
				if c.urlQueue.IsEmpty() && len(c.jobQueue) == 0 {
					log.Println("Finished processing all URLs. Exiting...")
					c.cancel() // Cancel context and exit
				}
			}
		}
	}()

	// Wait for all goroutines to finish
	c.wg.Wait()
	return nil
}

// Add URL to the queue
func (c *Crawler) SubmitURL(url string, depth int) {
	select {
	case <-c.ctx.Done():
		// Do nothing if context is done
		return
	case c.jobQueue <- CrawlJob{url, depth}:
		// Add job to the queue
	}
}

// Process collected results
func (c *Crawler) processResults() {
	defer c.wg.Done()

	var count int
	for article := range c.resultChan {
		// Save article to DB
		err := c.db.SaveArticle(article)
		if err != nil {
			log.Printf("Error saving article: %v", err)
			continue
		}

		count++
		if count%10 == 0 {
			log.Printf("Processed %d articles", count)
		}
	}

	log.Printf("Processed %d articles in total", count)
}

func (c *Crawler) Stop() {
	c.cancel()
	c.wg.Wait()
	c.db.Close()
}

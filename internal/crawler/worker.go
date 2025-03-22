package crawler

import (
	"log"
	"sync"

	"github.com/automatedtomato/go-web-crawler/internal/models"
	"github.com/automatedtomato/go-web-crawler/internal/queue"
	"github.com/automatedtomato/go-web-crawler/internal/storage"
)

// Represent crawled object
type CrawlJob struct {
	URL   string
	Depth int
}

// Represent crawl worker
type Worker struct {
	id         int
	jobQueue   chan CrawlJob
	crawler    *Crawler
	fetcher    *Fetcher
	parser     *Parser
	db         *storage.Database
	wg         *sync.WaitGroup
	urlQueue   *queue.URLQueue
	maxDepth   int
	resultChan chan *models.Article
}

func NewWorker(
	id int,
	jobQueue chan CrawlJob,
	crawler *Crawler,
	fetcher *Fetcher,
	parser *Parser,
	db *storage.Database,
	wg *sync.WaitGroup,
	urlQueue *queue.URLQueue,
	maxDepth int,
	resultChan chan *models.Article,
) *Worker {
	return &Worker{
		id:         id,
		jobQueue:   jobQueue,
		crawler:    crawler,
		fetcher:    fetcher,
		parser:     parser,
		db:         db,
		wg:         wg,
		urlQueue:   urlQueue,
		maxDepth:   maxDepth,
		resultChan: resultChan,
	}
}

func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()

		log.Printf("Worker %d: start", w.id)

		for job := range w.jobQueue {
			w.processJob(job)
		}

		log.Printf("Worker %d: done", w.id)
	}()
}

// Perform single job
func (w *Worker) processJob(job CrawlJob) {
	log.Printf("Worker %d: processing %s (depth: %d)", w.id, job.URL, job.Depth)

	// Check depth
	if job.Depth > w.maxDepth {
		return
	}

	// Fetch URL
	doc, err := w.fetcher.FetchURL(job.URL)
	if err != nil {
		log.Printf("Worker %d: failed to fetch %s: %v", w.id, job.URL, err)
		return
	}

	// Check if it's an article page
	if w.parser.IsArticlePage(doc) {
		// Extract article info
		article, err := w.parser.ParseArticle(doc, job.URL)
		if err != nil {
			log.Printf("Worker %d: failed to parse %s: %v", w.id, job.URL, err)
		} else if article != nil {
			// Send result via channel
			w.resultChan <- article
		}
	}

	// Extract links and add to queue
	links := w.fetcher.ExtractLink(doc, job.URL)
	for _, link := range links {
		if w.urlQueue.Push(link) {
			// Send new link to crawler
			w.crawler.SubmitURL(link, job.Depth+1)
		}
	}
}

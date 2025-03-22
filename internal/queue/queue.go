package queue

import "sync"

/*
 * Data structure of QUEUE implementation
 */

// Queue URLs for crawling target
type URLQueue struct {
	urls  []string
	mutex sync.Mutex
	seen  map[string]bool
}

func NewURLQueue(initialURLs []string) *URLQueue {
	seen := make(map[string]bool)
	for _, url := range initialURLs {
		seen[url] = true
	}

	return &URLQueue{
		urls: initialURLs,
		seen: seen,
	}
}

// Push: add URL to the queue
func (q *URLQueue) Push(url string) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.seen[url] {
		return false
	}

	q.urls = append(q.urls, url)
	q.seen[url] = true
	return true
}

// Pop: retrieve URL from the queue
func (q *URLQueue) Pop() (string, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.urls) == 0 {
		return "", false
	}

	url := q.urls[0]
	q.urls = q.urls[1:]
	return url, true
}

func (q *URLQueue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.urls) == 0
}

func (q *URLQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.urls)
}

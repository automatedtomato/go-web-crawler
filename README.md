# Go Concurrent Web Crawler

A highly performant concurrent web crawler built with Go, designed to efficiently collect and store news articles from various sources. This project implements worker pool pattern, rate limiting, and database integration, making it suitable as a data collection component in an MLOps pipeline.

## Features

- **Concurrent Processing**: Implements worker pool pattern for efficient parallel processing
- **Channel-based Synchronization**: Uses Go channels for communication between workers
- **Rate Limiting**: Prevents overloading websites with configurable per-host request limits
- **Depth Control**: Configurable crawling depth to manage scope
- **Database Integration**: Stores collected articles in a SQLite database using GORM
- **Configurable Settings**: Command-line flags for customizing crawler behavior
- **Graceful Shutdown**: Handles termination signals for clean shutdown
- **Article Detection**: Identifies and extracts information from news article pages
- **Duplicate Prevention**: Avoids processing the same URL multiple times

## Installation

### Requirements

- Go 1.16 or higher
- SQLite3

### Build and Install

```bash
# Clone the repository
git clone https://github.com/your-username/go-concurrent-crawler.git
cd go-concurrent-crawler

# Install dependencies
go mod download

# Build the project
go build -o crawler ./cmd
```

## Usage

### Basic Usage

Run the crawler with default settings:

```bash
./crawler
```

### Advanced Usage

Customize crawler behavior with command-line flags:

```bash
./crawler -workers 10 -depth 5 -timeout 60 -rps 2.0 -db ./news.db https://example.com/news
```

### Available Options

- `-workers`: Number of concurrent workers (default: 5)
- `-depth`: Maximum crawling depth (default: 3)
- `-timeout`: Maximum execution time in minutes (default: 30)
- `-rps`: Maximum requests per second per host (default: 1.0)
- `-db`: Database file path (default: ./crawler.db)
- Additional arguments are treated as seed URLs

## Project Structure

```
crawler/
├── cmd/
│   └── main.go               # Entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── crawler/
│   │   ├── crawler.go        # Core crawler functionality
│   │   ├── fetcher.go        # HTTP request handling
│   │   ├── parser.go         # HTML parsing functions
│   │   └── worker.go         # Worker pool implementation
│   ├── models/
│   │   └── article.go        # Data models
│   ├── queue/
│   │   └── queue.go          # URL queue management
│   ├── ratelimiter/
│   │   └── ratelimiter.go    # Rate limiting implementation
│   └── storage/
│       └── db.go             # Database operations
├── go.mod
└── go.sum
```

## Learning Objectives

This project was created to learn the following Golang concepts/features:

- **Concurrent Programming**: Using goroutines and channels to handle multiple tasks simultaneously
- **Worker Pool Pattern**: Implementing resource-controlled concurrency with a pool of workers
- **Data Synchronization**: Managing shared resources with mutexes and channels
- **Context Package**: Controlling cancellation and timeout across multiple goroutines
- **Rate Limiting**: Implementing polite crawling with per-host request limits
- **GORM Integration**: Using an ORM for database operations with SQLite
- **HTML Parsing**: Extracting structured data from web pages
- **Error Handling**: Robust error management in concurrent applications
- **Signal Handling**: Capturing OS signals for graceful shutdown

## Possible Extensions

Ideas for extending this project:

1. **Distributed Crawling**: Implement distributed crawling across multiple machines
2. **Content Classification**: Add ML-based content categorization
3. **Sentiment Analysis**: Analyze the sentiment of crawled articles
4. **API Layer**: Add a REST API to access the collected data
5. **Advanced Scheduling**: Implement intelligent recrawling based on site update patterns
6. **Content Deduplication**: Identify and merge similar articles from different sources
7. **Image Analysis**: Extract and analyze images from articles
8. **Named Entity Recognition**: Extract people, places, and organizations mentioned in articles
9. **Topic Modeling**: Automatically identify article topics using NLP techniques
10. **Integration with ML Pipelines**: Connect with Airflow or other ML workflow tools

## License

MIT License

## Contributing

While this project was created for learning purposes, improvement suggestions and bug reports are welcome. Feel free to create an Issue or Pull Request.

## Author

Hikaru Tomizawa

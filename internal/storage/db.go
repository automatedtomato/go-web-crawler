package storage

import (
	"log"
	"sync"

	"github.com/automatedtomato/go-web-crawler/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Wrapper for database
type Database struct {
	db    *gorm.DB
	mutex sync.Mutex
}

// Create a new DB
func NewDatabase(dbPath string) (*Database, error) {

	// Connect to DB
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Perform migrations
	err = db.AutoMigrate((&models.Article{}))
	if err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

// Close the DB
func (d *Database) Close() error {
	// Retrieve SQL DB then close
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Save article to DB
func (d *Database) SaveArticle(article *models.Article) error {
	// Lock DB preventing simultaneous writes
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if article already exists
	var existingArticle models.Article
	result := d.db.Where("url = ?", article.URL).First(&existingArticle)

	if result.Error == nil {
		log.Printf("Article already exists: %s", article.URL)
		return nil
	}

	return d.db.Create(article).Error
}

func (d *Database) GetArticles(limit int) ([]models.Article, error) {
	var articles []models.Article
	err := d.db.Limit(limit).Find(&articles).Error
	return articles, err
}

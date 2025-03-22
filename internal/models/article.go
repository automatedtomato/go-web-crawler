package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	URL         string    `json:"url" gorm:"unique"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
	Author      string    `json:"author"`
	ImageURL    string    `json:"image_url"`
}

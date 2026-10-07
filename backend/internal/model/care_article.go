package model

import "time"

// CareArticle is a gardening knowledge article.
type CareArticle struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Content   string    `gorm:"type:longtext" json:"content"`
	Cover     string    `gorm:"size:255" json:"cover"`
	TopicTag  string    `gorm:"size:32;index;not null" json:"topic_tag"`
	Status    string    `gorm:"size:16;default:published" json:"status"`
	ViewCount int       `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

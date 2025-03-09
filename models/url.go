package models

import (
	"time"
)

// URL 模型定义
type URL struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OriginalURL  string    `gorm:"size:2048;not null" json:"original_url"`
	ShortCode    string    `gorm:"size:10;unique;not null" json:"short_code"`
	CustomAlias  bool      `gorm:"default:false" json:"custom_alias"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	AccessCount  int64     `gorm:"default:0" json:"access_count"`
	LastAccessAt time.Time `json:"last_access_at"`
}

// URLStats 访问统计
type URLStats struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URLID     uint      `json:"url_id"`
	AccessIP  string    `gorm:"size:50" json:"access_ip"`
	UserAgent string    `gorm:"size:512" json:"user_agent"`
	Referer   string    `gorm:"size:512" json:"referer"`
	AccessAt  time.Time `json:"access_at"`
}

package services

import (
	"time"

	"github.com/keenJoe/go-url-shortener/database"
	"github.com/keenJoe/go-url-shortener/models"
)

// URLStatsData 统计数据结构
type URLStatsData struct {
	TotalAccess  int64       `json:"total_access"`
	LastAccessAt time.Time   `json:"last_access_at"`
	DailyStats   []DailyStat `json:"daily_stats"`
}

// DailyStat 每日统计
type DailyStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// GetURLStats 获取URL访问统计
func GetURLStats(shortCode string) (*URLStatsData, error) {
	var url models.URL
	if err := database.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return nil, err
	}

	// 获取过去30天的每日统计
	var dailyStats []DailyStat
	rows, err := database.DB.Raw(`
		SELECT DATE(access_at) as date, COUNT(*) as count 
		FROM url_stats 
		WHERE url_id = ? AND access_at > DATE_SUB(NOW(), INTERVAL 30 DAY)
		GROUP BY DATE(access_at)
		ORDER BY date DESC
	`, url.ID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stat DailyStat
		rows.Scan(&stat.Date, &stat.Count)
		dailyStats = append(dailyStats, stat)
	}

	return &URLStatsData{
		TotalAccess:  url.AccessCount,
		LastAccessAt: url.LastAccessAt,
		DailyStats:   dailyStats,
	}, nil
}

// RecordURLAccess 记录URL访问
func RecordURLAccess(urlID uint, ip, userAgent, referer string) error {
	stats := models.URLStats{
		URLID:     urlID,
		AccessIP:  ip,
		UserAgent: userAgent,
		Referer:   referer,
		AccessAt:  time.Now(),
	}
	return database.DB.Create(&stats).Error
}

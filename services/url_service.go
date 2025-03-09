package services

import (
	"errors"
	"time"

	"github.com/keenJoe/go-url-shortener/cache"
	"github.com/keenJoe/go-url-shortener/database"
	"github.com/keenJoe/go-url-shortener/models"
	"github.com/keenJoe/go-url-shortener/utils"
	"gorm.io/gorm"
)

// CreateShortURL 创建短链接
func CreateShortURL(originalURL string, customAlias string, expiration time.Duration) (string, error) {
	// 检查URL是否已存在
	var existingURL models.URL
	result := database.DB.Where("original_url = ?", originalURL).First(&existingURL)
	if result.Error == nil {
		// URL已存在，检查是否过期
		if existingURL.ExpiresAt.After(time.Now()) {
			return existingURL.ShortCode, nil
		}
	}

	var shortCode string
	if customAlias != "" {
		// 检查自定义别名是否合法
		if !utils.IsValidShortCode(customAlias) {
			return "", errors.New("自定义别名不合法")
		}

		// 检查自定义别名是否已被使用
		var existingAlias models.URL
		result := database.DB.Where("short_code = ?", customAlias).First(&existingAlias)
		if result.Error == nil {
			return "", errors.New("自定义别名已被使用")
		}

		shortCode = customAlias
	} else {
		// 生成随机短码
		for i := 0; i < 5; i++ { // 尝试5次
			shortCode = utils.GenerateRandomShortCode()
			var existingCode models.URL
			result := database.DB.Where("short_code = ?", shortCode).First(&existingCode)
			if result.Error == gorm.ErrRecordNotFound {
				break
			}
			if i == 4 {
				return "", errors.New("无法生成唯一短码")
			}
		}
	}

	// 设置过期时间
	var expiresAt time.Time
	if expiration > 0 {
		expiresAt = time.Now().Add(expiration)
	} else {
		// 默认不过期，设置为100年后
		expiresAt = time.Now().AddDate(100, 0, 0)
	}

	// 创建URL记录
	url := models.URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CustomAlias: customAlias != "",
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}

	if err := database.DB.Create(&url).Error; err != nil {
		return "", err
	}

	// 添加到缓存
	cache.SetURL(shortCode, originalURL, expiration)
	cache.SetURLLocal(shortCode, originalURL, expiration)

	// 添加到布隆过滤器
	utils.ShortCodeFilter.Add(shortCode)
	utils.OriginalURLFilter.Add(originalURL)

	return shortCode, nil
}

// GetOriginalURL 获取原始URL
func GetOriginalURL(shortCode string) (string, error) {
	// 检查短码是否合法
	if !utils.IsValidShortCode(shortCode) {
		return "", errors.New("短码不合法")
	}

	// 检查布隆过滤器
	if !utils.ShortCodeFilter.Contains(shortCode) {
		return "", errors.New("短码不存在")
	}

	// 先查本地缓存
	if originalURL, found := cache.GetURLLocal(shortCode); found {
		go updateAccessStats(shortCode)
		return originalURL, nil
	}

	// 查Redis缓存
	originalURL, err := cache.GetURL(shortCode)
	if err == nil {
		// 更新本地缓存
		cache.SetURLLocal(shortCode, originalURL, 30*time.Minute)
		go updateAccessStats(shortCode)
		return originalURL, nil
	}

	// 查数据库
	var url models.URL
	result := database.DB.Where("short_code = ?", shortCode).First(&url)
	if result.Error != nil {
		return "", result.Error
	}

	// 检查是否过期
	if url.ExpiresAt.Before(time.Now()) {
		return "", errors.New("链接已过期")
	}

	// 更新缓存
	cache.SetURL(shortCode, url.OriginalURL, time.Until(url.ExpiresAt))
	cache.SetURLLocal(shortCode, url.OriginalURL, 30*time.Minute)

	// 异步更新访问统计
	go updateAccessStats(shortCode)

	return url.OriginalURL, nil
}

// 更新访问统计
func updateAccessStats(shortCode string) {
	// 增加计数器
	cache.IncrementCounter(shortCode)

	// 更新数据库访问计数和最后访问时间
	database.DB.Model(&models.URL{}).
		Where("short_code = ?", shortCode).
		Updates(map[string]interface{}{
			"access_count":   gorm.Expr("access_count + 1"),
			"last_access_at": time.Now(),
		})
}

// DeleteExpiredURLs 删除过期URL
func DeleteExpiredURLs() error {
	return database.DB.Where("expires_at < ?", time.Now()).Delete(&models.URL{}).Error
}

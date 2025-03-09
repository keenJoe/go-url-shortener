package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/database"
	"github.com/keenJoe/go-url-shortener/models"
	"github.com/keenJoe/go-url-shortener/services"
)

// RedirectURL 重定向到原始URL
func RedirectURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// 获取原始URL
	originalURL, err := services.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在或已过期"})
		return
	}

	// 异步记录访问统计
	go func() {
		var url models.URL
		database.DB.Where("short_code = ?", shortCode).First(&url)

		services.RecordURLAccess(
			url.ID,
			c.ClientIP(),
			c.Request.UserAgent(),
			c.Request.Referer(),
		)
	}()

	// 重定向到原始URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

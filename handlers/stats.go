package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/services"
)

// GetURLStats 获取URL访问统计
func GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	stats, err := services.GetURLStats(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "链接不存在"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

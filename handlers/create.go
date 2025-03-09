package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/services"
)

// CreateURLRequest 创建URL请求
type CreateURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomAlias string `json:"custom_alias"`
	ExpiresIn   int64  `json:"expires_in"` // 过期时间（秒）
}

// CreateURLResponse 创建URL响应
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

// CreateURL 创建短链接
func CreateURL(c *gin.Context) {
	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 设置过期时间
	var expiration time.Duration
	if req.ExpiresIn > 0 {
		expiration = time.Duration(req.ExpiresIn) * time.Second
	}

	// 创建短链接
	shortCode, err := services.CreateShortURL(req.OriginalURL, req.CustomAlias, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建短链接URL
	baseURL := "http://" + c.Request.Host
	shortURL := baseURL + "/" + shortCode

	// 设置过期时间
	var expiresAt string
	if req.ExpiresIn > 0 {
		expiresAt = time.Now().Add(expiration).Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, CreateURLResponse{
		ShortCode:   shortCode,
		ShortURL:    shortURL,
		OriginalURL: req.OriginalURL,
		ExpiresAt:   expiresAt,
	})
}

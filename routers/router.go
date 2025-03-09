package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/handlers"
)

// Router 路由接口
type Router interface {
	Register(engine *gin.Engine)
}

// APIRouter API路由
type APIRouter struct{}

// Register 注册API路由
func (r *APIRouter) Register(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		api.POST("/shorten", handlers.CreateURL)
		api.GET("/stats/:shortCode", handlers.GetURLStats)
	}

	// 重定向路由
	engine.GET("/:shortCode", handlers.RedirectURL)
}

// InitRouter 初始化路由
func InitRouter() Router {
	return &APIRouter{}
}

// RegisterMiddleware 注册中间件
func RegisterMiddleware(engine *gin.Engine) {
	engine.Use(gin.Recovery())
	// engine.Use(middleware.Logger())
	// engine.Use(middleware.RateLimit())
}

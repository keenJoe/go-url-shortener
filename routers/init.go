package routers

import "github.com/gin-gonic/gin"

// 初始化所有路由
func InitRouter() *RouterGroup {
	return &RouterGroup{
		Routers: []IRouter{
			NewUserRouter(),
			// 在这里添加更多路由模块
		},
	}
}

// 注册全局中间件
func RegisterMiddleware(engine *gin.Engine) {
	// 全局中间件
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	// 可以添加更多全局中间件
}

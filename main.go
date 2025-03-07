package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/config"
	"github.com/keenJoe/go-url-shortener/pkg/database"
	"github.com/keenJoe/go-url-shortener/routers"
)

func main() {
	// 加载配置
	conf := config.LoadConfig()
	// 初始化数据库
	if err := database.InitDB(conf); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建 gin 实例
	app := gin.New()

	// 注册中间件
	routers.RegisterMiddleware(app)
	// 注册所有路由
	routerGroup := routers.InitRouter()
	routerGroup.Register(app)

	// 启动服务
	app.Run(":8080")
}

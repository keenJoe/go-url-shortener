package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/cache"
	"github.com/keenJoe/go-url-shortener/config"
	"github.com/keenJoe/go-url-shortener/database"
	"github.com/keenJoe/go-url-shortener/routers"
)

func main() {
	// 加载配置
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置gin模式
	gin.SetMode(conf.Server.Mode)

	// 初始化数据库
	if err := database.InitDB(conf); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化Redis
	if err := cache.InitRedis(conf); err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}

	// 创建gin实例
	app := gin.New()

	// 注册中间件
	routers.RegisterMiddleware(app)

	// 注册路由
	routerGroup := routers.InitRouter()
	routerGroup.Register(app)

	// 启动服务
	serverAddr := fmt.Sprintf(":%d", conf.Server.Port)
	if err := app.Run(serverAddr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

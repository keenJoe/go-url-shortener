package routers

import "github.com/gin-gonic/gin"

// 路由注册接口
type IRouter interface {
	Register(r *gin.Engine)
}

// 路由组注册器
type RouterGroup struct {
	Routers []IRouter
}

// 注册所有路由
func (rg *RouterGroup) Register(e *gin.Engine) {
	for _, router := range rg.Routers {
		router.Register(e)
	}
}

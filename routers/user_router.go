package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/controllers"
)

type UserRouter struct {
	userController *controllers.UserController
}

func NewUserRouter() *UserRouter {
	return &UserRouter{
		userController: controllers.NewUserController(),
	}
}

func (r *UserRouter) Register(e *gin.Engine) {
	api := e.Group("/api")
	users := api.Group("/users")
	{
		users.GET("", r.userController.GetUsers)
		users.POST("", r.userController.CreateUser)
		users.GET("/:id", r.userController.GetUser)
		users.PUT("/:id", r.userController.UpdateUser)
		users.DELETE("/:id", r.userController.DeleteUser)
	}
}

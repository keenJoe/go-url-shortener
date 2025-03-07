package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/keenJoe/go-url-shortener/models"
	"github.com/keenJoe/go-url-shortener/pkg/response"
	"github.com/keenJoe/go-url-shortener/services"
)

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(), // 注入服务
	}
}

// GetUsers 获取用户列表
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		response.Error(ctx, 500, "获取用户列表失败")
		return
	}
	response.Success(ctx, users)
}

// GetUser 获取单个用户
func (c *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.userService.GetUserByID(id)
	if err != nil {
		response.Error(ctx, 404, "用户不存在")
		return
	}
	response.Success(ctx, user)
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误："+err.Error())
		return
	}

	user, err := c.userService.CreateUser(req)
	if err != nil {
		response.Error(ctx, 500, "创建用户失败")
		return
	}
	response.Success(ctx, user)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误："+err.Error())
		return
	}

	user, err := c.userService.UpdateUser(id, req)
	if err != nil {
		response.Error(ctx, 500, "更新用户失败")
		return
	}
	response.Success(ctx, user)
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.userService.DeleteUser(id)
	if err != nil {
		response.Error(ctx, 500, "删除用户失败")
		return
	}
	response.Success(ctx, nil)
}

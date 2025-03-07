package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// 预定义业务状态码
const (
	CodeSuccess = 0    // 成功
	CodeError   = 1000 // 错误
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    CodeError,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithCode 带业务状态码的错误响应
func ErrorWithCode(c *gin.Context, httpStatus, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ValidationError 参数验证错误响应
func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "服务器内部错误")
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

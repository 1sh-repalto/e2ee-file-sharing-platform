package router

import (
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/handler"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := rg.Group("users")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)
	}
}

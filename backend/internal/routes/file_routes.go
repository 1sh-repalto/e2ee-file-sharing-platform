package router

import (
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/handler"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/middleware"
	"github.com/gin-gonic/gin"
)

func FileRoutes(rg *gin.RouterGroup, fileHandler *handler.FileHandler) {
	files := rg.Group("/files")
	files.Use(middleware.JWTAuthMiddleware())
	{
		files.POST("/", fileHandler.Upload)
		files.GET("/:id", fileHandler.GetByID)
		files.GET("/", fileHandler.ListByOwner)
		files.DELETE("/:id", fileHandler.Delete)
	}
}

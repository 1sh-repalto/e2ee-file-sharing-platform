package router

import (
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler *handler.UserHandler, fileHandler *handler.FileHandler) *gin.Engine {

	api := r.Group("/api")
	{
		AuthRoutes(api, userHandler)
		FileRoutes(api, fileHandler)
	}

	return r
}

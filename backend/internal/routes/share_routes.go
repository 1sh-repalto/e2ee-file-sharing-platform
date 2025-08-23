package router

import (
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/handler"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ShareRoutes(rg *gin.RouterGroup, shareHandler *handler.ShareHandler) {
	shares := rg.Group("/shares")
	shares.Use(middleware.JWTAuthMiddleware())
	{
		shares.POST("/", shareHandler.ShareFile)
		shares.GET("/", shareHandler.ListShares)
		shares.GET("/:file_id", shareHandler.GetShare)
		shares.DELETE("/:share_id", shareHandler.Unshare)
	}
}
package main

import (
	"log"
	"os"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/handler"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/repository"
	router "github.com/1sh-repalto/e2ee-file-sharing-platform/internal/routes"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/usecase"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found.")
	}

	db := config.NewPostgresPool()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	fileUsecase := usecase.NewFileUsecase(fileRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	fileHandler := handler.NewFileHandler(fileUsecase)

	r := gin.Default()
	router.SetupRouter(r, userHandler, fileHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server at : %v", port)
	}
}

package main

import (
	"log"
	"os"

	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/handler"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/repository"
	router "github.com/1sh-repalto/e2ee-file-sharing-platform/internal/routes"
	"github.com/1sh-repalto/e2ee-file-sharing-platform/internal/storage"
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

	minioStorage, err := storage.NewMinioStorage(
		"localhost:9000",
		os.Getenv("MINIO_ROOT_USER"),
		os.Getenv("MINIO_ROOT_PASSWORD"),
		false,
	)
	if err != nil {
		log.Fatalf("failed to init minio: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)
	shareRepo := repository.NewShareRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	fileUsecase := usecase.NewFileUsecase(fileRepo, shareRepo, minioStorage)
	shareUsecase := usecase.NewShareUsecase(shareRepo, fileRepo)

	userHandler := handler.NewUserHandler(userUsecase)
	fileHandler := handler.NewFileHandler(fileUsecase)
	shareHandler := handler.NewShareHandler(shareUsecase)

	r := gin.Default()
	router.SetupRouter(r, userHandler, fileHandler, shareHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server at : %v", port)
	}
}

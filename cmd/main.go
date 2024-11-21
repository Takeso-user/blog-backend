package main

import (
	"github.com/Takeso-user/blog-backend/config"
	"github.com/Takeso-user/blog-backend/pkg"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка переменных окружения
	config.LoadEnv()

	// Подключение к MongoDB
	cfg, err := config.ConnectToMongo()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer cfg.CloseMongo()

	// Инициализация репозиториев и сервисов
	userRepo := pkg.NewUserRepository(cfg.Database.Collection("users"))
	postRepo := pkg.NewPostRepository(cfg.Database.Collection("posts"))
	commentRepo := pkg.NewCommentRepository(cfg.Database.Collection("comments"))

	userService := pkg.NewUserService(userRepo)
	postService := pkg.NewPostService(postRepo)
	commentService := pkg.NewCommentService(commentRepo, userService)

	handler := pkg.NewHandler(postService, commentService, userService)

	// Настройка маршрутов
	router := gin.Default()
	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)

	api := router.Group("/api").Use(pkg.JWTMiddleware())
	api.POST("/posts", handler.CreatePost)
	api.GET("/posts", handler.GetPosts)
	api.POST("/posts/:id/comments", handler.AddComment)

	// Запуск сервера
	log.Println("Server is running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Takeso-user/blog-backend/config"
	"github.com/Takeso-user/blog-backend/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Loading environment variables...")
	config.LoadEnv()

	log.Println("Connecting to MongoDB...")
	cfg, err := config.ConnectToMongo()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		log.Println("Closing MongoDB connection...")
		cfg.CloseMongo()
	}()

	log.Println("Initializing repositories...")
	repository := pkg.NewRepository(cfg.Database)

	log.Println("Initializing services...")
	userService := pkg.NewUserService(repository.UserRepositoryInterface)
	postService := pkg.NewPostService(repository.PostRepositoryInterface)
	commentService := pkg.NewCommentService(repository.CommentRepositoryInterface, userService)

	log.Println("Initializing handlers...")
	handler := pkg.NewHandler(postService, commentService, userService)

	log.Println("Setting up router...")
	router := gin.Default()
	{
		router.POST("/auth/register", handler.Register)
		router.POST("/auth/login", handler.Login)
		router.GET("/auth/users", handler.GetUsers) //.Use(pkg.OwnerOrAdminMiddleware(postService))
	}
	api := router.Group("/api").Use(pkg.JWTMiddleware())
	{
		{
			api.POST("/posts", handler.CreatePost)
			api.GET("/posts", handler.GetPosts)
			api.GET("/posts/:id", handler.GetPostById)
			api.PATCH("/posts/:id", pkg.OwnerOrAdminMiddleware(postService), handler.UpdatePost)
			api.DELETE("/posts/:id", pkg.OwnerOrAdminMiddleware(postService), handler.DeletePost)
		}
		{
			api.POST("/posts/:id/comments", handler.AddComment)
			api.GET("/posts/:id/comments", handler.GetComments)
			api.GET("/posts/comments/", handler.GetAllComment)
			api.DELETE("/posts/comments/:commentID", pkg.OwnerOrAdminMiddleware(postService), handler.DeleteComment)
			api.PATCH("/posts/comments/:commentID", pkg.OwnerOrAdminMiddleware(postService), handler.UpdateComment)
		}
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Starting server on :8080...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

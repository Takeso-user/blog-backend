package main

import (
	"github.com/Takeso-user/blog-backend/config"
	"github.com/Takeso-user/blog-backend/pkg"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadEnv()

	cfg, err := config.ConnectToMongo()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer cfg.CloseMongo()

	userRepo := pkg.NewUserRepository(cfg.Database.Collection("users"))
	postRepo := pkg.NewPostRepository(cfg.Database.Collection("posts"))
	commentRepo := pkg.NewCommentRepository(cfg.Database.Collection("comments"))

	userService := pkg.NewUserService(userRepo)
	postService := pkg.NewPostService(postRepo)
	commentService := pkg.NewCommentService(commentRepo, userService)

	handler := pkg.NewHandler(postService, commentService, userService)

	router := gin.Default()
	router.POST("/auth/register", handler.Register)
	router.POST("/auth/login", handler.Login)
	router.GET("/auth/users", handler.GetUsers) //.Use(pkg.RoleMiddleware("Admin"))

	api := router.Group("/api").Use(pkg.JWTMiddleware())
	api.POST("/posts", handler.CreatePost)
	api.GET("/posts", handler.GetPosts)
	api.GET("/posts/:id", handler.GetPostById)
	api.PATCH("/posts/:id", pkg.OwnerOrAdminMiddleware(postService), handler.UpdatePost)
	api.DELETE("/posts/:id", pkg.OwnerOrAdminMiddleware(postService), handler.DeletePost)
	api.POST("/posts/:id/comments", handler.AddComment)
	api.GET("/posts/:id/comments", handler.GetComments)
	api.GET("/posts/comments/", handler.GetAllComment)
	api.DELETE("/posts/comments/:commentID", pkg.OwnerOrAdminMiddleware(postService), handler.DeleteComment)
	api.PATCH("/posts/comments/:commentID", pkg.OwnerOrAdminMiddleware(postService), handler.UpdateComment)

	log.Println("Server is running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

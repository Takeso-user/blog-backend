package main

import (
	"context"
	"errors"
	"github.com/Takeso-user/in-mem-cache/cache"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Takeso-user/blog-backend/config"
	_ "github.com/Takeso-user/blog-backend/docs"
	"github.com/Takeso-user/blog-backend/pkg"
)

func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}

//	@title			Blog API
//	@version		1.0
//	@description	This is a simple blog API

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	initLogger()
	logrus.Println("Loading environment variables...")
	config.LoadEnv()

	logrus.Println("Connecting to MongoDB...")
	cfg, err := config.ConnectToMongo()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		logrus.Println("Closing MongoDB connection...")
		cfg.CloseMongo()
	}()

	logrus.Println("Initializing repositories...")
	repository := pkg.NewRepository(cfg.Database)

	logrus.Println("Initializing cache...")
	cacheInstance := cache.NewCache(5 * time.Minute)

	logrus.Println("Initializing services...")
	userService := pkg.NewUserService(repository.UserRepositoryInterface, cacheInstance)
	postService := pkg.NewPostService(repository.PostRepositoryInterface, cacheInstance)
	commentService := pkg.NewCommentService(repository.CommentRepositoryInterface, userService, cacheInstance)

	logrus.Println("Initializing handlers...")
	handler := pkg.NewHandler(postService, commentService, userService)

	logrus.Println("Setting up router...")
	router := mux.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/auth/users", handler.GetUsers).Methods("GET")
	api := router.PathPrefix("/api").Subrouter()
	api.Use(pkg.JWTMiddleware)
	api.HandleFunc("/posts", handler.CreatePost).Methods("POST")
	api.HandleFunc("/posts", handler.GetPosts).Methods("GET")
	api.HandleFunc("/posts/{id}", handler.GetPostById).Methods("GET")
	api.Handle("/posts/{id}", pkg.OwnerOrAdminMiddleware(postService)(http.HandlerFunc(handler.UpdatePost))).Methods("PATCH")
	api.Handle("/posts/{id}", pkg.OwnerOrAdminMiddleware(postService)(http.HandlerFunc(handler.DeletePost))).Methods("DELETE")
	api.HandleFunc("/posts/{id}/comments", handler.AddComment).Methods("POST")
	api.HandleFunc("/posts/{id}/comments", handler.GetComments).Methods("GET")
	api.HandleFunc("/comments", handler.GetAllComment).Methods("GET")
	api.Handle("/posts/comments/{commentID}", pkg.OwnerOrAdminMiddleware(postService)(http.HandlerFunc(handler.DeleteComment))).Methods("DELETE")
	api.Handle("/posts/comments/{commentID}", pkg.OwnerOrAdminMiddleware(postService)(http.HandlerFunc(handler.UpdateComment))).Methods("PATCH")

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logrus.Println("Starting server on :8080...")
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Println("Server exiting")
}

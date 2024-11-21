package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	postService    *PostService
	commentService *CommentService
	userService    *UserService
}

func NewHandler(postService *PostService, commentService *CommentService, userService *UserService) *Handler {
	return &Handler{
		postService:    postService,
		commentService: commentService,
		userService:    userService,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Хэшируем пароль
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = hashedPassword
	if input.Role == "" {
		input.Role = "user"
	}

	// Сохраняем пользователя в базе
	if err := h.userService.CreateUser(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем пользователя из базы
	user, err := h.userService.GetUserByUsername(input.Username)
	log.Printf("in habdler getting user: %v", user)
	if err != nil {
		fmt.Printf("error %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Проверяем пароль
	err = CheckPassword(user.Password, input.Password)
	fmt.Printf("user.Password %s", user.Password)
	fmt.Printf("input.Password %s", input.Password)
	if err != nil {
		passwordError := fmt.Sprintf("Invalid username or password. Failed to check password: %v", err)

		c.JSON(http.StatusUnauthorized, gin.H{"error": passwordError})
		return
	}

	// Генерируем JWT
	token, err := GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
func (h *Handler) CreatePost(c *gin.Context) {
	var input Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.postService.CreatePost(input.Title, input.Content, input.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.postService.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) AddComment(c *gin.Context) {
	postID := c.Param("id")

	var input struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.commentService.AddComment(postID, input.UserID, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
}

func (h *Handler) GetComments(c *gin.Context) {
	postID := c.Param("id")

	comments, err := h.commentService.GetComments(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

package pkg

import (
	_ "github.com/Takeso-user/blog-backend/docs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type Handler struct {
	PostService    *PostService
	CommentService *CommentService
	UserService    *UserService
}

func NewHandler(postService *PostService, commentService *CommentService, userService *UserService) *Handler {
	return &Handler{
		PostService:    postService,
		CommentService: commentService,
		UserService:    userService,
	}
}

type Response map[string]interface{}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			input	body		User	true	"User object"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = hashedPassword
	if input.Role == "" {
		input.Role = "user"
	}

	if err := h.UserService.CreateUser(input); err != nil {
		log.Printf("Failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	log.Println("User registered successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login godoc
//
//	@Summary		Login a user
//
//	@Description	Login a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			input	body		User	true	"User object"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserService.GetUserByUsername(input.Username)
	log.Printf("Getting user: %v", user)
	if err != nil {
		log.Printf("Invalid username or password: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = CheckPassword(user.Password, input.Password)
	if err != nil {
		log.Printf("Failed to check password: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Println("User logged in successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// CreatePost godoc
//
//	@Summary		Create a new post
//
//	@Description	Create a new post
//
//	@Security		ApiKeyAuth
//
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			input	body		Post	true	"Post object"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/api/posts [post]
func (h *Handler) CreatePost(c *gin.Context) {
	var input Post
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.PostService.CreatePost(input.Title, input.Content, input.AuthorID)
	if err != nil {
		log.Printf("Unable to create post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create post"})
		return
	}

	log.Println("Post created successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

// GetPosts godoc
//
//	@Summary		Get all posts
//
//	@Description	Get all posts
//
//	@Security		ApiKeyAuth
//
//	@Tags			posts
//	@Produce		json
//	@Success		200	{array}		Post
//	@Failure		500	{object}	Response
//	@Router			/api/posts [get]
func (h *Handler) GetPosts(c *gin.Context) {
	posts, err := h.PostService.GetPosts()
	if err != nil {
		log.Printf("Unable to fetch posts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// AddComment godoc
//
//	@Summary		Add a comment to a post
//
//	@Description	Add a comment to a post
//
//	@Security		ApiKeyAuth
//
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Post ID"
//	@Param			input	body		Comment	true	"Comment object"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/api/posts/{id}/comments [post]
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

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.CommentService.AddComment(postID, userID.Hex(), input.Content)
	if err != nil {
		log.Printf("Unable to add comment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add comment"})
		return
	}

	log.Println("Comment added successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
}

// GetComments godoc
//
//	@Summary		Get comments for a post
//
//	@Description	Get comments for a post
//
//	@Security		ApiKeyAuth
//
//	@Tags			comments
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{array}		Comment
//	@Failure		500	{object}	Response
//	@Router			/api/posts/{id}/comments [get]
func (h *Handler) GetComments(c *gin.Context) {
	postID := c.Param("id")
	comments, err := h.CommentService.GetComments(postID)
	if err != nil {
		log.Printf("Unable to fetch comments: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetUsers godoc
//
//	@Summary		Get all users
//
//	@Description	Get all users
//
//	@Security		ApiKeyAuth
//
//	@Tags			users
//	@Produce		json
//	@Success		200	{array}		User
//	@Failure		500	{object}	Response
//	@Router			/auth/users [get]
func (h *Handler) GetUsers(context *gin.Context) {
	users, err := h.UserService.GetUsers()
	if err != nil {
		log.Printf("Unable to fetch users: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch users"})
		return
	}

	context.JSON(http.StatusOK, users)
}

// GetPostById godoc
//
//	@Summary		Get a post by ID
//	@Description	Get a post by ID
//	@Security		ApiKeyAuth
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	Post
//	@Failure		500	{object}	Response
//	@Router			/api/posts/{id} [get]
func (h *Handler) GetPostById(context *gin.Context) {
	postID := context.Param("id")
	post, err := h.PostService.GetPostById(postID)
	if err != nil {
		log.Printf("Unable to fetch post: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch post"})
		return
	}
	context.JSON(http.StatusOK, post)
}

// DeletePost godoc
//
//	@Summary		Get a post by ID
//	@Description	Get a post by ID
//	@Security		ApiKeyAuth
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		string	true	"Post ID"
//	@Success		200	{object}	Post
//	@Failure		500	{object}	Response
//	@Router			/api/posts/{id} [delete]
func (h *Handler) DeletePost(context *gin.Context) {
	postID := context.Param("id")
	err := h.PostService.DeletePost(postID)
	if err != nil {
		log.Printf("Unable to delete post: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete post"})
		return
	}
	log.Println("Post deleted successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetAllComment godoc
//
//	@Summary		Get all comments
//	@Description	Get all comments
//	@Security		ApiKeyAuth
//	@Tags			comments
//	@Produce		json
//	@Success		200	{array}		Comment
//	@Failure		500	{object}	Response
//	@Router			/api/posts/comments [get]
func (h *Handler) GetAllComment(context *gin.Context) {
	comments, err := h.CommentService.GetAllComment()
	if err != nil {
		log.Printf("Unable to fetch comments: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch comments"})
		return
	}
	context.JSON(http.StatusOK, comments)
}

// DeleteComment godoc
//
//	@Summary		Delete a comment
//	@Description	Delete a comment
//	@Tags			comments
//	@Produce		json
//	@Param			commentID	path		string	true	"Comment ID"
//	@Success		200			{object}	Response
//	@Failure		500			{object}	Response
//	@Router			/api/posts/comments/{commentID} [delete]
func (h *Handler) DeleteComment(context *gin.Context) {
	commentID := context.Param("commentID")
	err := h.CommentService.DeleteComment(commentID)
	if err != nil {
		log.Printf("Unable to delete comment: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete comment"})
		return
	}
	log.Println("Comment deleted successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update a post
//	@Security		ApiKeyAuth
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Post ID"
//	@Param			input	body		Post	true	"Post object"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/api/posts/{id} [patch]
func (h *Handler) UpdatePost(context *gin.Context) {
	postID := context.Param("id")

	var input Post
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.PostService.UpdatePost(objectID, input)
	if err != nil {
		log.Printf("Unable to update post: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update post",
			"err":   err.Error(),
		})
		return
	}
	log.Println("Post updated successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Post updated successfully", "post": post})
}

// UpdateComment godoc
//
//	@Summary		Update a comment
//	@Description	Update a comment
//	@Security		ApiKeyAuth
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			commentID	path		string	true	"Comment ID"
//	@Param			input		body		Comment	true	"Comment object"
//	@Success		200			{object}	Response
//	@Failure		400			{object}	Response
//	@Failure		500			{object}	Response
//	@Router			/api/posts/comments/{commentID} [patch]
func (h *Handler) UpdateComment(context *gin.Context) {
	commentID := context.Param("commentID")

	var input Comment
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.CommentService.UpdateComment(objectID, input)
	if err != nil {
		log.Printf("Unable to update comment: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update comment",
			"err":   err.Error(),
		})
		return
	}
	log.Println("Comment updated successfully")
	context.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully", "comment": comment})
}

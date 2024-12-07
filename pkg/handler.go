package pkg

import (
	"encoding/json"
	_ "github.com/Takeso-user/blog-backend/docs"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (h *Handler) SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "path/to/swagger/index.html")
}

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
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		logrus.Printf("Failed to hash password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	input.Password = hashedPassword
	if input.Role == "" {
		input.Role = "user"
	}

	if err := h.UserService.CreateUser(input); err != nil {
		logrus.Printf("Failed to register user: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	logrus.Println("User registered successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "User registered successfully"})
	if err != nil {
		return
	}
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
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserService.GetUserByUsername(input.Username)
	logrus.Printf("Getting user: %v", user)
	if err != nil {
		logrus.Printf("Invalid username or password: %v", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = CheckPassword(user.Password, input.Password)
	if err != nil {
		logrus.Printf("Failed to check password: %v", err)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		logrus.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	logrus.Println("User logged in successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"token": token})
	if err != nil {
		return
	}
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
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var input Post
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.PostService.CreatePost(input.Title, input.Content, input.AuthorID)
	if err != nil {
		logrus.Printf("Unable to create post: %v", err)
		http.Error(w, "Unable to create post", http.StatusInternalServerError)
		return
	}

	logrus.Println("Post created successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "Post created successfully"})
	if err != nil {
		return
	}
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
func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetPosts()
	if err != nil {
		logrus.Printf("Unable to fetch posts: %v", err)
		http.Error(w, "Unable to fetch posts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		return
	}
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

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	var input struct {
		UserID  string `json:"user_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.CommentService.AddComment(postID, userID.Hex(), input.Content)
	if err != nil {
		logrus.Printf("Unable to add comment: %v", err)
		http.Error(w, "Unable to add comment", http.StatusInternalServerError)
		return
	}

	logrus.Println("Comment added successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "Comment added successfully"})
	if err != nil {
		return
	}
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
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	comments, err := h.CommentService.GetComments(postID)
	if err != nil {
		logrus.Printf("Unable to fetch comments: %v", err)
		http.Error(w, "Unable to fetch comments", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		return
	}
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
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UserService.GetUsers()
	if err != nil {
		logrus.Printf("Unable to fetch users: %v", err)
		http.Error(w, "Unable to fetch users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		return
	}
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
func (h *Handler) GetPostById(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	logrus.Printf("!!!Getting post by ID: %s", postID)
	objectID, err := primitive.ObjectIDFromHex(postID)
	logrus.Printf("!!!Getting post by ID: %v", objectID)
	if err != nil {
		logrus.Printf("Error converting postID to ObjectID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetPostById(objectID.Hex())
	if err != nil {
		logrus.Printf("Error getting post by ID: %v", err)
		http.Error(w, "Unable to fetch post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		return
	}
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
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	err := h.PostService.DeletePost(postID)
	if err != nil {
		logrus.Printf("Unable to delete post: %v", err)
		http.Error(w, "Unable to delete post", http.StatusInternalServerError)
		return
	}
	logrus.Println("Post deleted successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "Post deleted successfully"})
	if err != nil {
		return
	}
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
func (h *Handler) GetAllComment(w http.ResponseWriter, r *http.Request) {
	r = nil
	comments, err := h.CommentService.GetAllComment()
	if err != nil {
		logrus.Printf("Unable to fetch comments: %v", err)
		http.Error(w, "Unable to fetch comments", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		return
	}
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
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["commentID"]

	err := h.CommentService.DeleteComment(commentID)
	if err != nil {
		logrus.Printf("Unable to delete comment: %v", err)
		http.Error(w, "Unable to delete comment", http.StatusInternalServerError)
		return
	}
	logrus.Println("Comment deleted successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "Comment deleted successfully"})
	if err != nil {
		return
	}
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

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	var input Post
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.UpdatePost(objectID, input)
	if err != nil {
		logrus.Printf("Unable to update post: %v", err)
		http.Error(w, "Unable to update post", http.StatusInternalServerError)
		return
	}
	logrus.Println("Post updated successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "Post updated successfully", "post": post})
	if err != nil {
		return
	}
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

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["commentID"]

	var input Comment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	objectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	comment, err := h.CommentService.UpdateComment(objectID, input)
	if err != nil {
		logrus.Printf("Unable to update comment: %v", err)
		http.Error(w, "Unable to update comment", http.StatusInternalServerError)
		return
	}
	logrus.Println("Comment updated successfully")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{"message": "Comment updated successfully", "comment": comment})
	if err != nil {
		return
	}
}

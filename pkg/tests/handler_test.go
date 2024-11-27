package tests

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/Takeso-user/blog-backend/pkg"
	"github.com/Takeso-user/blog-backend/pkg/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserService, globalCache)
	mockUserService.EXPECT().CreateUser(gomock.Any()).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{UserService: userService}
	router.POST("/auth/register", handler.Register)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(`{"username":"testuser","password":"password123"}`))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserService, globalCache)

	// Hash the password used in the test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUserService.EXPECT().GetUserByUsername("testuser").Return(pkg.User{Username: "testuser", Password: string(hashedPassword)}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{UserService: userService}
	router.POST("/auth/login", handler.Login)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{"username":"testuser","password":"password123"}`))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestGetUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserService, globalCache)
	mockUserService.EXPECT().GetUsers().Return([]pkg.User{{Username: "testuser"}}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{UserService: userService}
	router.GET("/auth/users", handler.GetUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
}

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostService := mocks.NewMockPostRepositoryInterface(ctrl)
	postService := pkg.NewPostService(mockPostService, globalCache)
	mockPostService.EXPECT().CreatePost(gomock.Any()).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{PostService: postService}
	router.POST("/posts", handler.CreatePost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", strings.NewReader(`{"title":"Test Title","content":"Test Content","author_id":"authorID"}`))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Post created successfully")
}

func TestAddComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	commentService := pkg.NewCommentService(mockCommentService, pkg.NewUserService(mockUserService, globalCache), globalCache)

	// Convert userID to primitive.ObjectID
	userID, _ := primitive.ObjectIDFromHex("000000000000000000000000")
	mockUserService.EXPECT().GetUserByID(userID.Hex()).Return(pkg.User{ID: userID, Username: "testuser"}, nil)
	mockCommentService.EXPECT().AddComment(gomock.Any()).Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{CommentService: commentService}
	router.POST("/posts/:id/comments", handler.AddComment)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts/postID/comments", strings.NewReader(`{"user_id":"000000000000000000000000","content":"Test Comment"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Comment added successfully")
}

func TestGetComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	commentService := pkg.NewCommentService(mockCommentService, pkg.NewUserService(mockUserService, globalCache), globalCache)
	mockCommentService.EXPECT().GetComments("postID").Return([]pkg.Comment{{Content: "Test Comment"}}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{CommentService: commentService}
	router.GET("/posts/:id/comments", handler.GetComments)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/postID/comments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Comment")
}

func TestGetPostById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostService := mocks.NewMockPostRepositoryInterface(ctrl)
	postService := pkg.NewPostService(mockPostService, globalCache)
	mockPostService.EXPECT().GetPostByID("postID").Return(pkg.Post{Title: "Test Title"}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{PostService: postService}
	router.GET("/posts/:id", handler.GetPostById)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/postID", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Title")
}

func TestDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostService := mocks.NewMockPostRepositoryInterface(ctrl)
	postService := pkg.NewPostService(mockPostService, globalCache)
	mockPostService.EXPECT().DeletePost("postID").Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{PostService: postService}
	router.DELETE("/posts/:id", handler.DeletePost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/postID", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Post deleted successfully")
}

func TestGetAllComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	commentService := pkg.NewCommentService(mockCommentService, pkg.NewUserService(mockUserService, globalCache), globalCache)
	mockCommentService.EXPECT().GetAllComment().Return([]pkg.Comment{{Content: "Test Comment"}}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{CommentService: commentService}
	router.GET("/comments", handler.GetAllComment)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/comments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Comment")
}

func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentService := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	commentService := pkg.NewCommentService(mockCommentService, pkg.NewUserService(mockUserService, globalCache), globalCache)
	mockCommentService.EXPECT().DeleteComment("commentID").Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{CommentService: commentService}
	router.DELETE("/comments/:commentID", handler.DeleteComment)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/comments/commentID", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Comment deleted successfully")
}

func TestUpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostService := mocks.NewMockPostRepositoryInterface(ctrl)
	postService := pkg.NewPostService(mockPostService, globalCache)

	validObjectID := primitive.NewObjectID()
	testPost := pkg.Post{
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	mockPostService.EXPECT().GetPostByID(validObjectID.Hex()).Return(pkg.Post{}, nil)
	mockPostService.EXPECT().UpdatePost(validObjectID, gomock.Any()).Return(testPost, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{PostService: postService}
	router.PUT("/posts/:id", handler.UpdatePost)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/"+validObjectID.Hex(), strings.NewReader(`{"title":"Updated Title","content":"Updated Content"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Post updated successfully")
	assert.Contains(t, w.Body.String(), "Updated Title")
}

func TestUpdateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validObjectID := primitive.NewObjectID()
	mockCommentService := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserService := mocks.NewMockUserRepositoryInterface(ctrl)
	commentService := pkg.NewCommentService(mockCommentService, pkg.NewUserService(mockUserService, globalCache), globalCache)

	// Expect the UpdateComment call with the correct arguments
	mockCommentService.EXPECT().UpdateComment(
		context.TODO(),
		bson.M{"_id": validObjectID},
		bson.M{"$set": bson.M{"content": "Updated Comment"}},
	).Return(pkg.Comment{Content: "Updated Comment"}, nil)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := &pkg.Handler{CommentService: commentService}
	router.PUT("/comments/:commentID", handler.UpdateComment)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/comments/"+validObjectID.Hex(), strings.NewReader(`{"content":"Updated Comment"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Comment updated successfully")
}

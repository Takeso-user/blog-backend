package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/Takeso-user/blog-backend/pkg"
	"github.com/Takeso-user/blog-backend/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_UserService_CreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userID, err := primitive.ObjectIDFromHex("000000000000000000000000")
	if err != nil {
		log.Printf("Error converting userID to ObjectID: %v", err)
	}
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserRepo)

	user := pkg.User{ID: userID, Username: "testuser"}

	mockUserRepo.EXPECT().CreateUser(user).Return(nil)

	err = userService.CreateUser(user)
	assert.NoError(t, err)
}

func Test_UserService_GetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserRepo)
	userID, err := primitive.ObjectIDFromHex("000000000000000000000000")
	if err != nil {
		log.Printf("Error converting userID to ObjectID: %v", err)
	}
	expectedUser := pkg.User{ID: userID, Username: "testuser"}

	mockUserRepo.EXPECT().GetUserByID(userID.Hex()).Return(expectedUser, nil)

	user, err := userService.GetUserByID(userID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}
func Test_PostService_CreatePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepositoryInterface(ctrl)
	postService := pkg.NewPostService(mockPostRepo)
	id, _ := primitive.ObjectIDFromHex("6745fdd700023c89744bd4e8")
	fixedTime := time.Date(2024, time.November, 26, 18, 1, 33, 0, time.UTC)
	post := pkg.Post{
		ID:        id,
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorID:  "000000000000000000000000",
		CreatedAt: fixedTime,
	}

	mockPostRepo.EXPECT().CreatePost(gomock.Any()).DoAndReturn(func(p pkg.Post) error {
		assert.Equal(t, post.Title, p.Title)
		assert.Equal(t, post.Content, p.Content)
		assert.Equal(t, post.AuthorID, p.AuthorID)
		return nil
	})

	err := postService.CreatePost(post.Title, post.Content, post.AuthorID)
	assert.NoError(t, err)
}

func Test_CommentService_AddComment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentRepo := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserRepo)
	commentService := pkg.NewCommentService(mockCommentRepo, userService)

	userID, err := primitive.ObjectIDFromHex("000000000000000000000000")
	if err != nil {
		log.Printf("Error converting userID to ObjectID: %v", err)
	}
	postID := "postID"
	content := "Test Comment"
	user := pkg.User{ID: userID, Username: "testuser"}

	mockUserRepo.EXPECT().GetUserByID(userID.Hex()).Return(user, nil)
	mockCommentRepo.EXPECT().AddComment(gomock.Any()).Return(nil)

	err = commentService.AddComment(postID, userID.Hex(), content)
	assert.NoError(t, err)
}

func Test_CommentService_UpdateComment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommentRepo := mocks.NewMockCommentRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := pkg.NewUserService(mockUserRepo)
	commentService := pkg.NewCommentService(mockCommentRepo, userService)

	commentID := primitive.NewObjectID()
	input := pkg.Comment{Content: "Updated Comment"}

	mockCommentRepo.EXPECT().UpdateComment(
		context.TODO(),
		bson.M{"_id": commentID},
		bson.M{"$set": bson.M{"content": input.Content}},
	).Return(pkg.Comment{Content: "Updated Comment"}, nil)

	updatedComment, err := commentService.UpdateComment(commentID, input)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Comment", updatedComment.Content)
}

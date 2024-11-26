package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepositoryInterface interface {
	CreatePost(post Post) error
	GetPosts() ([]Post, error)
	GetPostByID(postID string) (Post, error)
	DeletePost(postID string) error
	UpdatePost(id primitive.ObjectID, updateFields bson.M) (Post, error)
}

type CommentRepositoryInterface interface {
	AddComment(comment Comment) error
	GetComments(postID string) ([]Comment, error)
	GetAllComment() ([]Comment, error)
	DeleteComment(commentID string) error
	UpdateComment(ctx context.Context, filter, updateFields bson.M) (Comment, error)
}

type UserRepositoryInterface interface {
	CreateUser(user User) error
	GetUserByUsername(username string) (User, error)
	GetUserByID(userID string) (User, error)
	GetUsers() ([]User, error)
}

type Repository struct {
	PostRepositoryInterface
	CommentRepositoryInterface
	UserRepositoryInterface
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		PostRepositoryInterface:    NewPostRepository(db.Collection("posts")),
		CommentRepositoryInterface: NewCommentRepository(db.Collection("comments")),
		UserRepositoryInterface:    NewUserRepository(db.Collection("users")),
	}
}

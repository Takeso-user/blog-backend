package pkg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Service for posts
type PostService struct {
	repository *PostRepository
}
type UserService struct {
	repository *UserRepository
}

func NewUserService(repository *UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(user User) error {
	return s.repository.CreateUser(user)
}

func (s *UserService) GetUserByUsername(username string) (User, error) {
	return s.repository.GetUserByUsername(username)
}

func (s *UserService) GetUserByID(userID string) (User, error) {
	return s.repository.GetUserByID(userID)
}
func NewPostService(repository *PostRepository) *PostService {
	return &PostService{repository: repository}
}

func (s *PostService) CreatePost(title, content, authorID string) error {
	post := Post{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
	return s.repository.CreatePost(post)
}

func (s *PostService) GetPosts() ([]Post, error) {
	return s.repository.GetPosts()
}

type CommentService struct {
	repository  *CommentRepository
	userService *UserService
}

func NewCommentService(repository *CommentRepository, userService *UserService) *CommentService {
	return &CommentService{repository: repository, userService: userService}
}

func (s *CommentService) AddComment(postID, userID, content string) error {
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	comment := Comment{
		PostID:    postID,
		UserID:    userID,
		Username:  user.Username,
		Content:   content,
		CreatedAt: time.Now(),
	}
	return s.repository.AddComment(comment)
}

func (s *CommentService) GetComments(postID string) ([]Comment, error) {
	return s.repository.GetComments(postID)
}

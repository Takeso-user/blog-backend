package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type PostService struct {
	repository *PostRepository
}
type UserService struct {
	repository *UserRepository
}
type CommentService struct {
	repository  *CommentRepository
	userService *UserService
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

func (s *UserService) GetUsers() ([]User, error) {
	return s.repository.GetUsers()
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

func (s *PostService) GetPostById(id string) (Post, error) {
	return s.repository.GetPostById(id)
}

func (s *PostService) DeletePost(id string) error {
	return s.repository.DeletePost(id)
}

func (s *PostService) UpdatePost(id primitive.ObjectID, input Post) (Post, error) {
	currentPost, err := s.repository.GetPostById(id.Hex())
	if err != nil {
		return Post{}, err
	}
	updateFields := bson.M{}

	if input.Title != "" {
		updateFields["title"] = input.Title
	}
	if input.Content != "" {
		updateFields["content"] = input.Content
	}

	updateFields["author_id"] = currentPost.AuthorID
	updateFields["created_at"] = currentPost.CreatedAt
	updatedPost, err := s.repository.UpdatePost(id, updateFields)
	log.Printf("updated post: %v", updatedPost)
	log.Printf("updated fields: %v", updateFields)
	if err != nil {
		return Post{}, err
	}
	return updatedPost, nil
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

func (s *CommentService) GetAllComment() ([]Comment, error) {
	return s.repository.GetAllComment()

}

func (s *CommentService) DeleteComment(id string) error {
	return s.repository.DeleteComment(id)
}

func (s *CommentService) UpdateComment(id primitive.ObjectID, input Comment) (Comment, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"content": input.Content,
		},
	}
	updatedComment, err := s.repository.UpdateComment(context.TODO(), filter, update)
	if err != nil {
		return Comment{}, err
	}
	return updatedComment, nil
}

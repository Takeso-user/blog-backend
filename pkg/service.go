package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type PostService struct {
	Repository PostRepositoryInterface
}
type UserService struct {
	Repository UserRepositoryInterface
}
type CommentService struct {
	Repository  CommentRepositoryInterface
	UserService *UserService
}

func NewUserService(repository UserRepositoryInterface) *UserService {
	return &UserService{Repository: repository}
}

func (s *UserService) CreateUser(user User) error {
	log.Println("Creating user:", user.Username)
	err := s.Repository.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	return err
}

func (s *UserService) GetUserByUsername(username string) (User, error) {
	log.Println("Getting user by username:", username)
	user, err := s.Repository.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user by username: %v", err)
	}
	return user, err
}

func (s *UserService) GetUserByID(userID string) (User, error) {
	log.Println("Getting user by ID:", userID)
	user, err := s.Repository.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
	}
	return user, err
}

func (s *UserService) GetUsers() ([]User, error) {
	log.Println("Getting all users")
	users, err := s.Repository.GetUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
	}
	return users, err
}

func NewPostService(repository PostRepositoryInterface) *PostService {
	return &PostService{Repository: repository}
}

func (s *PostService) CreatePost(title, content, authorID string) error {
	log.Println("Creating post:", title)
	post := Post{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
	err := s.Repository.CreatePost(post)
	if err != nil {
		log.Printf("Error creating post: %v", err)
	}
	return err
}

func (s *PostService) GetPosts() ([]Post, error) {
	log.Println("Getting all posts")
	posts, err := s.Repository.GetPosts()
	if err != nil {
		log.Printf("Error getting posts: %v", err)
	}
	return posts, err
}

func (s *PostService) GetPostById(id string) (Post, error) {
	log.Println("Getting post by ID:", id)
	post, err := s.Repository.GetPostByID(id)
	if err != nil {
		log.Printf("Error getting post by ID: %v", err)
	}
	return post, err
}

func (s *PostService) DeletePost(id string) error {
	log.Println("Deleting post by ID:", id)
	err := s.Repository.DeletePost(id)
	if err != nil {
		log.Printf("Error deleting post: %v", err)
	}
	return err
}

func (s *PostService) UpdatePost(id primitive.ObjectID, input Post) (Post, error) {
	log.Println("Updating post by ID:", id.Hex())
	currentPost, err := s.Repository.GetPostByID(id.Hex())
	if err != nil {
		log.Printf("Error getting post by ID: %v", err)
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
	updatedPost, err := s.Repository.UpdatePost(id, updateFields)
	if err != nil {
		log.Printf("Error updating post: %v", err)
		return Post{}, err
	}
	log.Printf("Updated post: %v", updatedPost)
	return updatedPost, nil
}

func NewCommentService(repository CommentRepositoryInterface, userService *UserService) *CommentService {
	return &CommentService{Repository: repository, UserService: userService}
}

func (s *CommentService) AddComment(postID, userID, content string) error {
	log.Println("Adding comment to post:", postID)
	user, err := s.UserService.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return err
	}
	comment := Comment{
		PostID:    postID,
		UserID:    userID,
		Username:  user.Username,
		Content:   content,
		CreatedAt: time.Now(),
	}
	err = s.Repository.AddComment(comment)
	if err != nil {
		log.Printf("Error adding comment: %v", err)
	}
	return err
}

func (s *CommentService) GetComments(postID string) ([]Comment, error) {
	log.Println("Getting comments for post:", postID)
	comments, err := s.Repository.GetComments(postID)
	if err != nil {
		log.Printf("Error getting comments: %v", err)
	}
	return comments, err
}

func (s *CommentService) GetAllComment() ([]Comment, error) {
	log.Println("Getting all comments")
	comments, err := s.Repository.GetAllComment()
	if err != nil {
		log.Printf("Error getting comments: %v", err)
	}
	return comments, err
}

func (s *CommentService) DeleteComment(id string) error {
	log.Println("Deleting comment by ID:", id)
	err := s.Repository.DeleteComment(id)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
	}
	return err
}

func (s *CommentService) UpdateComment(id primitive.ObjectID, input Comment) (Comment, error) {
	log.Println("Updating comment by ID:", id.Hex())
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"content": input.Content,
		},
	}
	updatedComment, err := s.Repository.UpdateComment(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		return Comment{}, err
	}
	log.Printf("Updated comment: %v", updatedComment)
	return updatedComment, nil
}

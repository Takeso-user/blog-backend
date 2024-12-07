package pkg

import (
	"context"
	"github.com/Takeso-user/in-mem-cache/cache"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostService struct {
	Repository PostRepositoryInterface
	Cache      *cache.Cache
}

type UserService struct {
	Repository UserRepositoryInterface
	Cache      *cache.Cache
}

type CommentService struct {
	Repository  CommentRepositoryInterface
	UserService *UserService
	Cache       *cache.Cache
}

func NewPostService(repository PostRepositoryInterface, cache *cache.Cache) *PostService {
	return &PostService{Repository: repository, Cache: cache}
}

func NewCommentService(repository CommentRepositoryInterface, userService *UserService, cache *cache.Cache) *CommentService {
	return &CommentService{Repository: repository, UserService: userService, Cache: cache}
}

func NewUserService(repository UserRepositoryInterface, cache *cache.Cache) *UserService {
	return &UserService{Repository: repository, Cache: cache}
}

func (s *UserService) CreateUser(user User) error {
	logrus.Println("Creating user:", user.Username)
	err := s.Repository.CreateUser(user)
	if err != nil {
		logrus.Printf("Error creating user: %v", err)
	}
	return err
}

func (s *UserService) GetUserByUsername(username string) (User, error) {
	logrus.Println("Getting user by username:", username)
	cachedUser, found := s.Cache.Get(username)
	if found {
		logrus.Printf("Found user in cache: %v", cachedUser)
		return cachedUser.(User), nil
	}

	user, err := s.Repository.GetUserByUsername(username)
	if err != nil {
		logrus.Printf("Error getting user by username: %v", err)
	}
	s.Cache.Set(username, user)
	return user, err
}

func (s *UserService) GetUserByID(userID string) (User, error) {
	logrus.Println("Getting user by ID:", userID)
	cachedUser, found := s.Cache.Get(userID)
	if found {
		logrus.Printf("Found user in cache: %v", cachedUser)
		return cachedUser.(User), nil
	}
	user, err := s.Repository.GetUserByID(userID)
	if err != nil {
		logrus.Printf("Error getting user by ID: %v", err)
	}
	s.Cache.Set(user.Username, user)
	return user, err
}

func (s *UserService) GetUsers() ([]User, error) {
	logrus.Println("Getting all users")
	users, err := s.Repository.GetUsers()
	if err != nil {
		logrus.Printf("Error getting users: %v", err)
	}
	return users, err
}

func (s *PostService) CreatePost(title, content, authorID string) error {
	logrus.Println("Creating post:", title)
	post := Post{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Content:   content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}
	err := s.Repository.CreatePost(post)
	if err != nil {
		logrus.Printf("Error creating post: %v", err)
	}
	return err
}

func (s *PostService) GetPosts() ([]Post, error) {
	logrus.Println("Getting all posts")
	posts, err := s.Repository.GetPosts()
	if err != nil {
		logrus.Printf("Error getting posts: %v", err)
	}
	return posts, err
}

func (s *PostService) GetPostById(id string) (Post, error) {
	logrus.Println("Getting post by ID:", id)
	cachedPost, found := s.Cache.Get(id)
	if found {
		logrus.Printf("Found post in cache: %v", cachedPost)
		return cachedPost.(Post), nil
	}
	post, err := s.Repository.GetPostByID(id)
	if err != nil {
		logrus.Printf("Error getting post by ID: %v", err)
	}
	s.Cache.Set(id, post)
	return post, err
}

func (s *PostService) DeletePost(id string) error {
	logrus.Println("Deleting post by ID:", id)
	err := s.Repository.DeletePost(id)
	if err != nil {
		logrus.Printf("Error deleting post: %v", err)
	}
	return err
}

func (s *PostService) UpdatePost(id primitive.ObjectID, input Post) (Post, error) {
	logrus.Println("Updating post by ID:", id.Hex())
	currentPost, err := s.Repository.GetPostByID(id.Hex())
	if err != nil {
		logrus.Printf("Error getting post by ID: %v", err)
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
		logrus.Printf("Error updating post: %v", err)
		return Post{}, err
	}
	logrus.Printf("Updated post: %v", updatedPost)
	return updatedPost, nil
}

func (s *CommentService) AddComment(postID, userID, content string) error {
	logrus.Println("Adding comment to post:", postID)
	user, err := s.UserService.GetUserByID(userID)
	if err != nil {
		logrus.Printf("Error getting user by ID: %v", err)
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
		logrus.Printf("Error adding comment: %v", err)
	}
	return err
}

func (s *CommentService) GetComments(postID string) ([]Comment, error) {
	logrus.Println("Getting comments for post:", postID)
	comments, err := s.Repository.GetComments(postID)
	if err != nil {
		logrus.Printf("Error getting comments: %v", err)
	}
	return comments, err
}

func (s *CommentService) GetAllComment() ([]Comment, error) {
	logrus.Println("Getting all comments")
	comments, err := s.Repository.GetAllComment()
	if err != nil {
		logrus.Printf("Error getting comments: %v", err)
	}
	return comments, err
}

func (s *CommentService) DeleteComment(id string) error {
	logrus.Println("Deleting comment by ID:", id)
	err := s.Repository.DeleteComment(id)
	if err != nil {
		logrus.Printf("Error deleting comment: %v", err)
	}
	return err
}

func (s *CommentService) UpdateComment(id primitive.ObjectID, input Comment) (Comment, error) {
	logrus.Println("Updating comment by ID:", id.Hex())
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"content": input.Content,
		},
	}
	updatedComment, err := s.Repository.UpdateComment(context.TODO(), filter, update)
	if err != nil {
		logrus.Printf("Error updating comment: %v", err)
		return Comment{}, err
	}
	logrus.Printf("Updated comment: %v", updatedComment)
	return updatedComment, nil
}

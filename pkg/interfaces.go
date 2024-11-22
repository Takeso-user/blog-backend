package pkg

type AuthorizationInterface interface {
	CreateUser(user User) error
	GetUserByUsername(username string) (User, error)
	GetUserByID(userID string) (User, error)
	GetUsers() ([]User, error)
}

type PostRepositoryInterface interface {
	CreatePost(post Post) error
	GetPosts() ([]Post, error)
}

type CommentRepositoryInterface interface {
	AddComment(comment Comment) error
	GetComments(postID string) ([]Comment, error)
}

type UserRepositoryInterface interface {
	CreateUser(user User) error
	GetUserByUsername(username string) (User, error)
	GetUserByID(userID string) (User, error)
	GetUsers() ([]User, error)
}

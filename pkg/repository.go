package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type PostRepository struct {
	collection *mongo.Collection
}

type UserRepository struct {
	collection *mongo.Collection
}

type CommentRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (r *UserRepository) CreateUser(user User) error {
	log.Println("Creating user:", user.Username)
	_, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (User, error) {
	log.Println("Getting user by username:", username)
	var user User
	err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		log.Printf("Error getting user by username: %v", err)
	} else {
		log.Printf("Found user: %v", user)
	}
	return user, err
}

func (r *UserRepository) GetUserByID(userID string) (User, error) {
	log.Println("Getting user by ID:", userID)
	var user User
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Error converting userID to ObjectID: %v", err)
		return user, err
	}
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
	}
	return user, err
}

func (r *UserRepository) GetUsers() ([]User, error) {
	log.Println("Getting all users")
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error getting users: %v", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Printf("Error decoding users: %v", err)
		return nil, err
	}
	return users, nil
}

func NewPostRepository(collection *mongo.Collection) *PostRepository {
	return &PostRepository{collection: collection}
}

func (r *PostRepository) CreatePost(post Post) error {
	log.Println("Creating post:", post.Title)
	_, err := r.collection.InsertOne(context.TODO(), post)
	if err != nil {
		log.Printf("Error creating post: %v", err)
	}
	return err
}

func (r *PostRepository) GetPosts() ([]Post, error) {
	log.Println("Getting all posts")
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error getting posts: %v", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var posts []Post
	if err = cursor.All(context.TODO(), &posts); err != nil {
		log.Printf("Error decoding posts: %v", err)
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetPostByID(id string) (Post, error) {
	log.Println("Getting post by ID:", id)
	var post Post
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting postID to ObjectID: %v", err)
		return post, err
	}
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&post)
	if err != nil {
		log.Printf("Error getting post by ID: %v", err)
	}
	return post, err
}

func (r *PostRepository) DeletePost(id string) error {
	log.Println("Deleting post by ID:", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting postID to ObjectID: %v", err)
		return err
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error deleting post: %v", err)
	}
	return err
}

func (r *PostRepository) UpdatePost(id primitive.ObjectID, updateFields bson.M) (Post, error) {
	log.Println("Updating post by ID:", id.Hex())
	var updatedPost Post
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedPost)
	if err != nil {
		log.Printf("Error updating post: %v", err)
	}
	return updatedPost, err
}

func NewCommentRepository(collection *mongo.Collection) *CommentRepository {
	return &CommentRepository{collection: collection}
}

func (r *CommentRepository) AddComment(comment Comment) error {
	log.Println("Adding comment to post:", comment.PostID)
	_, err := r.collection.InsertOne(context.TODO(), comment)
	if err != nil {
		log.Printf("Error adding comment: %v", err)
	}
	return err
}

func (r *CommentRepository) GetComments(postID string) ([]Comment, error) {
	log.Println("Getting comments for post:", postID)
	filter := bson.M{"post_id": postID}
	opts := options.Find().SetSort(bson.M{"created_at": 1})

	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Printf("Error getting comments: %v", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var comments []Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		log.Printf("Error decoding comments: %v", err)
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) GetAllComments() ([]Comment, error) {
	log.Println("Getting all comments")
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error getting comments: %v", err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var comments []Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		log.Printf("Error decoding comments: %v", err)
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) DeleteComment(id string) error {
	log.Println("Deleting comment by ID:", id)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting commentID to ObjectID: %v", err)
		return err
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
	}
	return err
}

func (r *CommentRepository) UpdateComment(ctx context.Context, filter, update bson.M) (Comment, error) {
	log.Println("Updating comment with filter:", filter)
	var updatedComment Comment
	err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedComment)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
	}
	return updatedComment, err
}

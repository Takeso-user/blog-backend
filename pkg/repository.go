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
	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (User, error) {
	var user User
	err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	log.Printf("found user: %v", user)
	return user, err
}

func (r *UserRepository) GetUserByID(userID string) (User, error) {
	var user User
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user, err
	}
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	return user, err
}

func (r *UserRepository) GetUsers() ([]User, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func NewPostRepository(collection *mongo.Collection) *PostRepository {
	return &PostRepository{collection: collection}
}

func (r *PostRepository) CreatePost(post Post) error {
	_, err := r.collection.InsertOne(context.TODO(), post)
	return err
}

func (r *PostRepository) GetPosts() ([]Post, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var posts []Post
	if err = cursor.All(context.TODO(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetPostById(id string) (Post, error) {
	var post Post
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return post, err
	}
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&post)
	return post, err
}

func (r *PostRepository) DeletePost(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	return err
}

func (r *PostRepository) UpdatePost(id primitive.ObjectID, updateFields bson.M) (Post, error) {
	var updatedPost Post
	err := r.collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": updateFields},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedPost)
	return updatedPost, err
}

func NewCommentRepository(collection *mongo.Collection) *CommentRepository {
	return &CommentRepository{collection: collection}
}

func (r *CommentRepository) AddComment(comment Comment) error {
	_, err := r.collection.InsertOne(context.TODO(), comment)
	return err
}

func (r *CommentRepository) GetComments(postID string) ([]Comment, error) {
	filter := bson.M{"post_id": postID}
	opts := options.Find().SetSort(bson.M{"created_at": 1})

	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var comments []Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) GetAllComment() ([]Comment, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("error closing cursor: %v", err)
		}
	}(cursor, context.TODO())

	var comments []Comment
	if err = cursor.All(context.TODO(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepository) DeleteComment(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	return err
}

func (r *CommentRepository) UpdateComment(ctx context.Context, filter, update bson.M) (Comment, error) {
	var updatedComment Comment
	err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedComment)
	return updatedComment, err
}

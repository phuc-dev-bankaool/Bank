package repositories

import (
	"app/be/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (ins *UserRepository) CreateUser(ctx context.Context, user *models.User) (id string, err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Tạo ID nếu chưa có
	if user.ID == bson.NilObjectID {
		user.ID = bson.NewObjectID()
	}

	_, err = ins.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	id = user.ID.Hex()
	return id, nil
}

func (ins *UserRepository) UpdateUser(ctx context.Context, user *models.User) (id string, err error) {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	result, err := ins.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	if result.UpsertedID != nil {
		id = user.ID.Hex()
	}
	return id, nil
}

func (ins *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ins.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}

func (ins *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ins.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &user, nil
}

func (ins *UserRepository) DeleteAllUsers(ctx context.Context) error {
	result, err := ins.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Error deleting users: %v", err)
		return err
	}
	log.Printf("Deleted %d users", result.DeletedCount)
	return nil
}

func (ins *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := ins.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

package repository

import (
	"context"
	"errors"

	"github.com/RibunLoc/microservices-learn/auth-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RedisMongo struct {
	Collection *mongo.Collection
	JwtSecret  string
}

func NewUserRepo(db *mongo.Database) *RedisMongo {
	return &RedisMongo{
		Collection: db.Collection("users"),
	}
}

func (r *RedisMongo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *RedisMongo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisMongo) FindByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(ctx, map[string]interface{}{"_id": oid}).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update Passowrd MongoDB
func (r *RedisMongo) UpdatePassword(ctx context.Context, userID, hash string) error {
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"password": hash}}
	res, err := r.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *RedisMongo) UpdateUserFields(ctx context.Context, userId string, updateFields bson.M) error {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return errors.New("invalid user ID")
	}
	filter := bson.M{"_id": oid}

	update := bson.M{"$set": updateFields}

	res, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

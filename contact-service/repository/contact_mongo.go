package repository

import (
	"context"

	"github.com/RibunLoc/WebPersonalBackend/contact-service/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactMongo struct {
	Collection *mongo.Collection
}

func NewContactRepo(db *mongo.Database) *ContactMongo {
	return &ContactMongo{Collection: db.Collection("contacts")}
}

func (r *ContactMongo) Create(ctx context.Context, c *model.Contact) error {
	_, err := r.Collection.InsertOne(ctx, c)
	return err
}

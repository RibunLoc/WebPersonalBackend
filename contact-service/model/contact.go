package model

import (
	"github.com/RibunLoc/WebPersonalBackend/contact-service/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Message   string             `bson:"message" json:"message"`
	CreatedAt *util.CustomTime   `bson:"created_at" json:"created_at"`
}

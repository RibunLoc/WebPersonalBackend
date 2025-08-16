package model

import (
	"github.com/RibunLoc/microservices-learn/auth-service/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Fullname  string             `bson:"full_name" json:"full_name"`
	Role      string             `bson:"role" json:"role"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt *util.CustomTime   `bson:"created_at" json:"created_at"`
}

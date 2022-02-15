package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty"`
	Email    string             `json:"email,omitempty" validate:"required,email"`
	Password string             `json:"password,omitempty" validate:"required"`
}

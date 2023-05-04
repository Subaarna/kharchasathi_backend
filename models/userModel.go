package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id, omitempty"`
	FirstName   string             `json:"firstName" validate:"required, min=3, max=30"`
	LastName    string             `json:"lastName" validate:"required, min=3, max=30"`
	Email       string             `json:"email" validate:"required, email"`
	Password    string             `json:"password" validate:"required, min=8, max=30"`
	PhoneNumber string             `json:"phoneNumber" validate:"omitempty, min=10, max=10"`
	Description string             `json:"description" validate:"omitempty, min=10, max=1000"`
	Balance     int                `json:"balance" validate:"omitempty"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=8, max=30"`
}

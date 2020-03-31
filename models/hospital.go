package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Hospital in-mem represents a hospitals collection document
type Hospital struct {
	ID      primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Address string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone   string             `json:"phone,omitempty" bson:"phone,omitempty"`
}
